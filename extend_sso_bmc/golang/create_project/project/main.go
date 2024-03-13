package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	pb "sso_center/gen/cluster-contract"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listenPort = ":8080"
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
func (r *microservice) PutSettings(context.Context, *pb.PutSsoSettingsRequest) (*pb.PutSsoSettingsResponse, error) {
	//реализация rpc
	//...

	return nil, errors.New("not implemented")
}
