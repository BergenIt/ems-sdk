package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	pb "sso_center/gen/cluster-contract"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listenPort                 = ":8080"
	HTTPPort                   = 80
	HTTPsPort                  = 443
	HTTPProtocol               = "http://"
	HTTPsProtocol              = "https://"
	accountServicePageEndpoint = "/redfish/v1/AccountService/"
	managersPageEndpoint       = "/redfish/v1/Managers/"
	CAEndpoint                 = `/redfish/v1/Managers/iDRAC.Embedded.1/Oem/Dell/DelliDRACCardService/Actions/DelliDRACCardService.ImportSSLCertificate`
	BASE_DN                    = "dc=bergen,dc=ems"
	SSO_HOST                   = "some_sso_host"
	SSO_PORT                   = "1234"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run() error {
	port := os.Getenv("ServicePort")
	if port == "" {
		port = listenPort
	}

	// Создаем инстанс сервиса.
	m := microservice{}

	// Создаем инстанс сервера.
	server := grpc.NewServer()

	// Регистрируем сервис.
	pb.RegisterSsoCenterServer(server, &m)

	// Регистрируем рефлексию для сервиса, чтобы получать информацию об общедоступных RPC (опционально).
	reflection.Register(server)

	// Создаем листененра.
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		return fmt.Errorf("create listener: %s", err)
	}

	log.Printf("microservice start serving on port %q", port)

	// Запускаем gRPC сервер.
	return server.Serve(lis)
}

// Инстанс сервиса с реализацией RPC.
type microservice struct {
	pb.UnimplementedSsoCenterServer
}

// RPC для установления настроек LDAP-авторизации на BMC.
func (r *microservice) PutSettings(ctx context.Context, req *pb.PutSsoSettingsRequest) (*pb.PutSsoSettingsResponse, error) {
	creds, address, err := findCreds(req.Device.Connectors, pb.ConnectorProtocol_CONNECTOR_PROTOCOL_REDFISH)
	if err != nil {
		return nil, fmt.Errorf("find redfish creds: %s", err)
	}

	redfishClient := newRedfishClient(creds.Login, creds.Password, address, creds.Port)

	if err := putSettings(redfishClient, req.TargetState); err != nil {
		return nil, fmt.Errorf("put ldap settings: %s", err)
	}

	if req.TargetState == pb.SsoState_SSO_STATE_ACTIVE {
		if err := setLDAPAttrsDell(redfishClient, req.SsoDn, req.SsoPassword); err != nil {
			return nil, fmt.Errorf("set ldap attrs: %s", err)
		}

		if err := loadLDAPCA(); err != nil {
			return nil, fmt.Errorf("load CA: %s", err)
		}
	}

	return nil, errors.New("not implemented")
}

func findCreds(in []*pb.DeviceConnector, protocol pb.ConnectorProtocol) (*pb.Credential, string, error) {
	for _, connector := range in {
		for _, creds := range connector.Credentials {
			if creds.Protocol == protocol {
				if creds.Login == "" || creds.Password == "" {
					return nil, connector.Address, fmt.Errorf("login or password can not be empty")
				}

				return creds, "", nil
			}
		}
	}

	return nil, "", fmt.Errorf("creds not found")
}

// Включение/выключение лдапа на бмс Dell
func putSettings(client *RedfishClient, state pb.SsoState) error {
	body := createLDAPManageBody(state)

	if err := client.PatchData(accountServicePageEndpoint, body); err != nil {
		return err
	}

	return nil
}

func loadLDAPCA(client *RedfishClient) error {
	ca, err := DownloadAcmeCAFromStorage(d.configEms.S3Host, d.ca)
	if err != nil {
		return fmt.Errorf("download CA error: %s", err)
	}

	if err := client.PostData(CAEndpoint, createLoadCABody(ca)); err != nil {
		return fmt.Errorf("post CA error: %s", err)
	}

	return nil
}

type ManagersDell struct {
	Members []struct {
		OdataID string `json:"@odata.id"`
	} `json:"Members"`
}

// Установка настроек лдапа в бмс Dell (работает только при включении лдапа)
func setLDAPAttrsDell(client *RedfishClient, ssoDn, ssoPassword string) error {
	b, err := client.GetPage(managersPageEndpoint)
	if err != nil {
		return fmt.Errorf("get managers page error: %s", err)
	}

	managers := &ManagersDell{}
	if err := json.Unmarshal(b, managers); err != nil {
		return fmt.Errorf("unmarshal managers page error: %s", err)
	}

	if len(managers.Members) <= 0 {
		return fmt.Errorf("managers not found")
	}

	body := createLDAPAttrsBodyDell(ssoDn, ssoPassword)
	attributesEndpoint := managers.Members[0].OdataID + "/Attributes"

	if err := client.PatchData(attributesEndpoint, body); err != nil {
		return err
	}

	return nil
}

func DownloadAcmeCAFromStorage(addr string, ca *x509.CertPool) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            ca,
				InsecureSkipVerify: true,
			},
		},
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/roots.pem", addr))
	if err != nil {
		return "", fmt.Errorf("download ca error: %s", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error downloading ca: %s", string(b))
	}

	return string(b), nil
}
