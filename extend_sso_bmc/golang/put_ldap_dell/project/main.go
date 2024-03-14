package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	pb "sso_center/gen/cluster-contract"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	LISTEN_PORT = ":8080"

	HTTP_PORT            = 80
	HTTPS_PORT           = 443
	HTTP_PROTOCOL        = "http://"
	HTTPS_PROTOCOL       = "https://"
	ACCOUNT_SERVICE_PAGE = "/redfish/v1/AccountService/"
	MANAGERS_PAGE        = "/redfish/v1/Managers/"
	SET_CA_PAGE          = `/redfish/v1/Managers/iDRAC.Embedded.1/Oem/Dell/DelliDRACCardService/Actions/DelliDRACCardService.ImportSSLCertificate`
	BASE_DN              = "dc=bergen,dc=ems"
	SSO_HOST             = "some_sso_host"
	SSO_PORT             = "1234"
	CA_LOCAL_PATH        = "./roots.pem"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run() error {
	port := os.Getenv("ServicePort")
	if port == "" {
		port = LISTEN_PORT
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
	lis, err := net.Listen("tcp", port)
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
	log.Printf("got request with state %v", req.TargetState.String())

	// Поиск редфиш кредов
	creds, address, err := findCreds(req.Device.Connectors, pb.ConnectorProtocol_CONNECTOR_PROTOCOL_REDFISH)
	if err != nil {
		return nil, fmt.Errorf("find redfish creds: %s", err)
	}

	// Создание редфиш клиента
	redfishClient := newRedfishClient(creds.Login, creds.Password, address, creds.Port)

	// Установка настроек в зависимости от статуса
	if err := putSettings(redfishClient, req.TargetState); err != nil {
		return nil, fmt.Errorf("put ldap settings: %s", err)
	}

	// Если статус - активировать, то загружаем дополнительные параметры и сертификат
	if req.TargetState == pb.SsoState_SSO_STATE_ACTIVE {
		if err := setLDAPAttrsDell(redfishClient, req.SsoDn, req.SsoPassword); err != nil {
			return nil, fmt.Errorf("set ldap attrs: %s", err)
		}

		if err := loadLDAPCA(redfishClient); err != nil {
			return nil, fmt.Errorf("load CA: %s", err)
		}
	}

	return &pb.PutSsoSettingsResponse{
		Result: &pb.OperationResult{
			DeviceId: req.Device.DeviceId,
			State:    pb.OperationState_OPERATION_STATE_SUCCESS,
		},
	}, nil
}

func findCreds(in []*pb.DeviceConnector, protocol pb.ConnectorProtocol) (*pb.Credential, string, error) {
	for _, connector := range in {
		for _, creds := range connector.Credentials {
			if creds.Protocol == protocol {
				fmt.Println(connector.GetAddress())
				if creds.Login == "" || creds.Password == "" {
					return nil, "", fmt.Errorf("login or password can not be empty")
				}

				return creds, connector.Address, nil
			}
		}
	}

	return nil, "", fmt.Errorf("creds not found")
}

// Включение/выключение лдапа на бмс Dell
func putSettings(client *RedfishClient, state pb.SsoState) error {
	body := createLDAPManageBody(state)

	if err := client.PatchData(ACCOUNT_SERVICE_PAGE, body); err != nil {
		return err
	}

	return nil
}

func loadLDAPCA(client *RedfishClient) error {
	ca, err := loadCA()
	if err != nil {
		return fmt.Errorf("load CA: %s", err)
	}

	if err := client.PostData(SET_CA_PAGE, createLoadCABody(ca)); err != nil {
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
	b, err := client.GetPage(MANAGERS_PAGE)
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

func loadCA() (string, error) {
	ca, err := os.ReadFile(CA_LOCAL_PATH)
	if err != nil {
		return "", fmt.Errorf("read body error: %s", err)
	}

	return string(ca), nil
}
