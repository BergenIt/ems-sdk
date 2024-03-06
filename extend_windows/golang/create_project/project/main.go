package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	pb "windows-handler/gen/cluster-contract"

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
	// Создаем инстанс сервиса.
	m := microservice{}

	// Создаем инстанс сервера.
	server := grpc.NewServer()

	// Регистрируем сервис.
	pb.RegisterWindowsManagerServer(server, &m)

	// Регистрируем рефлексию для сервиса, чтобы получать информацию об общедоступных RPC (опционально).
	reflection.Register(server)

	// Создаем листененра.
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		return fmt.Errorf("create listener: %s", err)
	}

	log.Printf("microservice start serving on port %q", listenPort)

	// Запускаем gRPC сервер.
	return server.Serve(lis)
}

// Инстанс сервиса с реализацией RPC.
type microservice struct {
	pb.UnimplementedWindowsManagerServer
}

// RPC по сбору инвентарных данных по ОЗУ с ОС Windows.
func (r *microservice) CollectMemory(
	ctx context.Context,
	req *pb.CollectWindowsMemoryRequest,
) (*pb.CollectWindowsMemoryResponse, error) {
	//реализация rpc
	//...

	return nil, errors.New("not implemented")
}
