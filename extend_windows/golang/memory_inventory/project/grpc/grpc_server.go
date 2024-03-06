package cmd

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "windows-handler/gen/cluster-contract"

	"windows-handler/memory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listenPort = ":8080"
)

type microservice struct {
	pb.UnimplementedWindowsManagerServer
}

func Serve() error {
	m := microservice{}
	server := grpc.NewServer()

	pb.RegisterWindowsManagerServer(server, &m)
	reflection.Register(server)

	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		return fmt.Errorf("create listener: %s", err)
	}

	log.Printf("microservice start serving on port %q", listenPort)
	return server.Serve(lis)
}

func (r *microservice) CollectMemory(
	ctx context.Context,
	req *pb.CollectWindowsMemoryRequest,
) (*pb.CollectWindowsMemoryResponse, error) {
	resp := &pb.CollectWindowsMemoryResponse{
		Memory: &pb.DeviceMemory{
			DeviceIdentity: &pb.DeviceDataIdentity{
				DeviceId: req.Device.DeviceId,
				Source:   pb.ServiceSource_SERVICE_SOURCE_WINDOWS_MANAGER,
			},
		},
	}

	for _, connector := range req.Device.Connectors {
		wmiCreds, _ := getWMICreds(connector.Credentials)
		if wmiCreds != nil {
			memoryInventory, err := memory.GetMemoryInv(ctx, connector.Address, wmiCreds)
			if err != nil {
				continue
			}

			if memoryInventory != nil {
				resp.Memory.Memories = append(resp.Memory.Memories, memoryInventory...)

				break
			}
		}
	}

	return resp, nil
}

func getWMICreds(in []*pb.Credential) (*pb.Credential, error) {
	if len(in) == 0 {
		return nil, fmt.Errorf("wmi creds list is empty")
	}

	for _, creds := range in {
		if creds.Protocol == pb.ConnectorProtocol_CONNECTOR_PROTOCOL_WMI {
			if creds.Login == "" || creds.Password == "" {
				return nil, fmt.Errorf("wmi creds can not be empty")
			}

			return creds, nil
		}
	}

	return nil, fmt.Errorf("wmi creds not found")
}
