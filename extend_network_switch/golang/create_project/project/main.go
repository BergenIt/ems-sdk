package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	pb "network/gen/cluster-contract"

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

	lis, err := net.Listen("tcp", listenPort)
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

func (m *microservice) CreateConfig(context.Context, *pb.CreateNetworkConfigRequest) (*pb.CreateNetworkConfigResponse, error) {
	// реализация rpc
	// ...

	return nil, errors.New("not implemented")
}
