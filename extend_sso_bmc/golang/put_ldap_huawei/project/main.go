package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	pb "sso_center/gen/cluster-contract"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	LISTEN_PORT = ":8080"

	HTTP_PORT                = 80
	HTTPS_PORT               = 443
	HTTP_PROTOCOL            = "http://"
	HTTPS_PROTOCOL           = "https://"
	ACCOUNT_SERVICE_PAGE     = "/redfish/v1/AccountService/"
	LDAPPageEndpoint         = "/redfish/v1/AccountService/LdapService/"
	SessionsPageEndpoint     = "/redfish/v1/SessionService/Sessions/"
	LDAPSetCAPageEndpoint    = "/redfish/v1/AccountService/LdapService/LdapControllers/1/Actions/HwLdapController.ImportCert"
	LDAPSettingsPageEndpoint = "/redfish/v1/AccountService/LdapService/LdapControllers/1"
	BASE_DN                  = "dc=bergen,dc=ems"
	SSO_HOST                 = "some_sso_host"
	SSO_PORT                 = "1234"
	CA_LOCAL_PATH            = "/app/roots.pem"
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

	headers, err := redfishClient.GetAllHeadersFromPage(
		SessionsPageEndpoint,
		http.MethodPost,
		strings.NewReader(fmt.Sprintf(`{"UserName": "%s","Password": "%s"}`, redfishClient.Username, redfishClient.Password)),
	)
	if err != nil {
		return nil, fmt.Errorf("get headers error: %s", err)
	}

	location := headers.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("empty location header")
	}

	defer func() {
		// закрытие сессии
		if err := redfishClient.DeleteByPage(location); err != nil {
			log.Printf("failed close huawei session: %s", err)
		}
	}()

	xauth := headers.Get("X-Auth-Token")
	if xauth == "" {
		return nil, fmt.Errorf("empty xauth header")
	}

	if err := manageLDAP(redfishClient, req.TargetState, xauth); err != nil {
		return nil, fmt.Errorf("manage ldap error: %s", err)
	}

	if req.TargetState == pb.SsoState_SSO_STATE_ACTIVE {
		if err := loadLDAPCA(redfishClient, xauth); err != nil {
			return nil, fmt.Errorf("upload CA error: %s", err)
		}

		if err := setLDAPSettings(redfishClient, xauth); err != nil {
			return nil, fmt.Errorf("set ldap settings error: %s", err)
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

func loadLDAPCA(client *RedfishClient, xauth string) error {
	var headers map[string]string = map[string]string{
		"X-Auth-Token": xauth,
	}

	ca, err := loadCa()
	if err != nil {
		return fmt.Errorf("load CA: %s", err)
	}

	if err := client.PostDataWithHeaders(LDAPSetCAPageEndpoint, createLDAPSetCABody(string(ca)), headers); err != nil {
		return fmt.Errorf("run redfish request error: %s", err)
	}

	return nil
}

func manageLDAP(client *RedfishClient, state pb.SsoState, xauth string) error {
	etag, err := client.GetHeaderFromPage(LDAPPageEndpoint, "ETag", http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("get etag error: %s", err)
	}

	var headers map[string]string = map[string]string{
		"If-Match":     etag,
		"X-Auth-Token": xauth,
	}

	if err := client.PatchDataWithHeaders(LDAPPageEndpoint, createLDAPManageBody(state), headers); err != nil {
		return fmt.Errorf("run redfish request error: %s", err)
	}

	return nil
}

// Установка настроек лдапа в бмс Huawei (работает только при включении лдапа)
func setLDAPSettings(client *RedfishClient, xauth string) error {
	etag, err := client.GetHeaderFromPage(LDAPSettingsPageEndpoint, "ETag", http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("get etag error: %s", err)
	}

	var headers map[string]string = map[string]string{
		"If-Match":     etag,
		"X-Auth-Token": xauth,
	}

	if err := client.PatchDataWithHeaders(LDAPSettingsPageEndpoint, createLDAPSettingsBody(), headers); err != nil {
		return fmt.Errorf("run redfish request error: %s", err)
	}

	return nil
}

func loadCa() (string, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://traefik:7071/roots.pem", nil)
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	response, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("http GET: %s", err)
	}
	defer response.Body.Close()

	ca, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %s", err)
	}

	return string(ca), nil
}
