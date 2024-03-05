package main

import (
	"context"
	"errors"
	"fmt"
	pb "hypervisor/gen/cluster-contract"
	"log"
	"net"
	"os"

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
	pb.RegisterHypervisorManagerServer(server, &m)

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
	pb.UnimplementedHypervisorManagerServer
}

// RPC по сбору списка виртуальных машин с гипервизра ESXI.
func (r *microservice) CollectVirtialMachinesList(context.Context, *pb.CollectVirtialMachinesListRequest) (*pb.CollectVirtialMachinesListResponse, error) {
	//реализация rpc
	//...

	return nil, errors.New("not implemented")
}
