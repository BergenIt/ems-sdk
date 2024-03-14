package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "network/gen/cluster-contract"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listenPort = ":8080"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
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
	pb.RegisterNetworkManagerServer(server, &m)

	// Регистрируем рефлексию для сервиса, чтобы получать информацию об общедоступных RPC (опционально).
	reflection.Register(server)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}

	log.Printf("microservice start serving on port %q", port)

	// Запускаем gRPC сервер.
	return server.Serve(lis)
}

type microservice struct {
	pb.UnimplementedNetworkManagerServer
}

type sshConnInfo struct {
	addr  string
	login string
	pass  string
	port  int32
}

const createConfigCmd = "export"

func (m *microservice) CreateConfig(ctx context.Context, req *pb.CreateNetworkConfigRequest) (*pb.CreateNetworkConfigResponse, error) {
	fmt.Println("got request")

	info, err := extractSSHConnInfo(req.Device.Connectors)
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User:            info.login,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(info.pass),
		},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", info.addr, info.port), cfg)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(createConfigCmd); err != nil {
		return nil, err
	}

	res := pb.CreateNetworkConfigResponse{
		Result: &pb.OperationResult{
			DeviceId: req.Device.DeviceId,
			State:    pb.OperationState_OPERATION_STATE_SUCCESS,
			Output:   b.String(),
		},
	}

	return &res, nil
}

func extractSSHConnInfo(connectors []*pb.DeviceConnector) (sshConnInfo, error) {
	for _, conn := range connectors {
		for _, creds := range conn.Credentials {
			if creds.Protocol == pb.ConnectorProtocol_CONNECTOR_PROTOCOL_SSH {
				res := sshConnInfo{
					addr:  conn.Address,
					login: creds.Login,
					pass:  creds.Password,
					port:  creds.Port,
				}
				return res, nil
			}
		}
	}
	return sshConnInfo{}, fmt.Errorf("ssh credentials were not found")
}
