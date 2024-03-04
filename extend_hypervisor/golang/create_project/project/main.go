package main

import (
	"context"
	"errors"
	"fmt"
	pb "hypervisor/gen/cluster-contract"
	"log"
	"net"

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
	m := microservice{}
	server := grpc.NewServer()

	pb.RegisterHypervisorManagerServer(server, &m)
	reflection.Register(server)

	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		return fmt.Errorf("create listener: %s", err)
	}

	log.Printf("microservice start serving on port %q", listenPort)
	return server.Serve(lis)
}

type microservice struct {
	pb.UnimplementedHypervisorManagerServer
}

func (r *microservice) CollectVirtialMachinesList(context.Context, *pb.CollectVirtialMachinesListRequest) (*pb.CollectVirtialMachinesListResponse, error) {
	//реализация rpc
	//...

	return nil, errors.New("not implemented")
}
