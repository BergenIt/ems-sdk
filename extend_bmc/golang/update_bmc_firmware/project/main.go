package main

import (
	rfish "bmc-handler/adapter"
	pb "bmc-handler/gen/cluster-contract"
	redfish "bmc-handler/service"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ErrNotValidURL     = errors.New("firmware url can't be proceed, its not valid")
	ErrEmptySftpPort   = errors.New("sftp port can't be empty, please enter port to env variable")
	ErrUnknownProtocol = errors.New("can't proccess update via unknown protocol")
)

const (
	listenPort = ":8080"

	sftpPort = "2222"
)

const (
	insecureHTTP = "http://"
	secureHTTPS  = "https://"
	SFTP         = "sftp://"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run() error {
	// Создаем инстанс сервиса.
	m := microservice{}

	// Создаем инстанс сервера.
	server := grpc.NewServer()

	// Регистрируем сервис.
	pb.RegisterBmcManagerServer(server, &m)

	// Регистрируем рефлексию для сервиса, чтобы получать информацию об общедоступных RPC (опционально).
	reflection.Register(server)

	// Создаем листененра.
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		return fmt.Errorf("create listener: %s", err)
	}

	log.Printf("microservice start serving on port %q\n", listenPort)

	// Запускаем gRPC сервер.
	return server.Serve(lis)
}

// Инстанс сервиса с реализацией RPC.
type microservice struct {
	pb.UnimplementedBmcManagerServer
}

// RPC по обновлению прошивки BMC.
func (r *microservice) BmcFirmwareUpdate(
	ctx context.Context,
	req *pb.BmcFirmwareUpdateRequest,
) (*pb.BmcFirmwareUpdateResponse, error) {
	resp := &pb.BmcFirmwareUpdateResponse{
		Result: &pb.OperationResult{
			DeviceId: req.Device.DeviceId,
		},
	}

	// Перебор коннекторов для получения коннектора с протоколом Redfish
	for _, connector := range req.Device.Connectors {
		redfishCreds, _ := getRedfishCreds(connector.Credentials)
		if redfishCreds != nil {
			err := proccessUpdate(req.FirmwareUrl, redfishCreds, connector.Address)
			if err != nil {
				resp.Result.State = pb.OperationState_OPERATION_STATE_FAILED
				resp.Result.Output = err.Error()
			} else {
				resp.Result.State = pb.OperationState_OPERATION_STATE_SUCCESS
			}
			return resp, nil
		}
	}

	return nil, errors.New("not implemented")
}

// Получение кредов по протоколу Redfish
func getRedfishCreds(in []*pb.Credential) (*pb.Credential, error) {
	if len(in) == 0 {
		return nil, fmt.Errorf("redfish creds list is empty")
	}

	for _, creds := range in {
		if creds.Protocol == pb.ConnectorProtocol_CONNECTOR_PROTOCOL_REDFISH {
			if creds.Login == "" || creds.Password == "" {
				return nil, fmt.Errorf("redfish creds can not be empty")
			}

			return creds, nil
		}
	}

	return nil, fmt.Errorf("redfish creds not found")
}

func proccessUpdate(firmwareUrl string, redfishCreds *pb.Credential, ip string) error {
	log.Printf("Start update bmc process with Firmware URL: %s, Address: %s\n", firmwareUrl, ip)

	sftpPath, err := sftpBuilder(firmwareUrl, sftpPort)
	if err != nil {
		return err
	}

	host := redfish.Host{
		IP:       ip,
		User:     redfishCreds.Login,
		Password: redfishCreds.Password,
		IsHTTPS:  redfishCreds.Port == 443,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	cfg := redfish.RedfishCFG{
		Vendor: "Huawei",
		Model:  "2288H V5",
		SimpleUpdate: redfish.SimpleUpdate{
			UpdateURL:     "/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
			UpdateRequest: `{"ImageURI": "{filePath}", "TransferProtocol": "{protocol}"}`,
		},
		Optional: redfish.Optional{},
	}

	gofish := rfish.New()
	redfishService := redfish.NewRedfish(gofish)

	return redfishService.SimpleUpdate(sftpPath, cfg, host)
}

func sftpBuilder(firmwareUrl, sftpPortValue string) (string, error) {
	if firmwareUrl == "" || !strings.Contains(firmwareUrl, "http") {
		return "", ErrNotValidURL
	}

	if sftpPortValue == "" {
		return "", ErrEmptySftpPort
	}

	u, err := url.Parse(firmwareUrl)
	if err != nil {
		return firmwareUrl, fmt.Errorf("sftpBuilder failed: %w", err)
	}

	host := strings.Split(u.Host, ":")[0]
	path := strings.TrimPrefix(u.Path, "/")

	var ipHost string

	ips, _ := net.LookupIP(host)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			ipHost = ipv4.String()
		}
	}

	return fmt.Sprintf("sftp://%s:%s@%s:%s/%s",
		"minio", "minio_key", ipHost, sftpPort,
		path), nil

}
