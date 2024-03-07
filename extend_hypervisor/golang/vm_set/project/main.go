package main

import (
	"context"
	"errors"
	"fmt"
	pb "hypervisor/gen/cluster-contract"
	"log"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	listenPort  = ":8080"
	ipv4Pattern = `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
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
	pb.UnimplementedHypervisorManagerServer
}

// RPC по сбору списка виртуальных машин с гипервизра ESXI.
func (r *microservice) CollectVirtualMachinesList(ctx context.Context, req *pb.CollectVirtualMachinesListRequest) (*pb.CollectVirtualMachinesListResponse, error) {
	log.Printf("got reqiest with vendor")
	// Проверка на то, что тип гипервизора ESXI.
	if req.Hypervisor.GetType() != pb.HypervisorType_HYPERVISOR_TYPE_ESXI {
		return nil, errors.New("only esxi hypervisor supported")
	}
	fmt.Printf("%+v\n", req.Hypervisor)
	address := req.Hypervisor.GetAddress()

	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("hypervisor address can not be empty")
	}

	// Парсинг адреса.
	u, err := soap.ParseURL(address)
	if err != nil {
		return nil, fmt.Errorf("parse URL [%s]: %s", address, err)
	}

	// Установка кредов.
	u.User = url.UserPassword(req.Hypervisor.GetLogin(), req.Hypervisor.GetPassword())

	// Создание клиента.
	client, err := govmomi.NewClient(context.Background(), u, true)
	if err != nil {
		return nil, fmt.Errorf("create client: %s", err)
	}

	// Сбор данных с гипервизора, используя библиотеку govmomi.
	virtualMachines := []mo.VirtualMachine{}
	m := view.NewManager(client.Client)
	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, fmt.Errorf("create container view: %s", err)
	}
	defer v.Destroy(ctx)

	if err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"guest", "config", "summary"}, virtualMachines); err != nil {
		return nil, fmt.Errorf("retrieve: %s", err)
	}

	// Создание выходной структуры.
	resp := pb.CollectVirtualMachinesListResponse{
		VirtualMachines: &pb.VirtualMachines{
			Identity: &pb.HypervisorIdentity{
				HypervisorId: req.GetHypervisor().GetIdentity().GetHypervisorId(),
				HypervisorOwnerId: req.GetHypervisor().GetIdentity().GetHypervisorOwnerId(),
				AccessObjectId: req.GetHypervisor().GetIdentity().GetAccessObjectId(),
			},
			VirtualMachines: make([]*pb.VirtualMachine, 0, len(virtualMachines)),
		},
	}

	// Парсинг ответа от гипервизора.
	regexp, err := regexp.Compile(ipv4Pattern)
	if err != nil {
		return nil, fmt.Errorf("compile IPv4 pattern: %s", err)
	}

	for _, vm := range virtualMachines {
		machine := pb.VirtualMachine{
			Name: vm.Summary.Config.Name,
		}

		for _, ni := range vm.Guest.Net {
			if ni.Network != "" {

				machine.Networks = append(machine.Networks, &pb.VirtualMachineNetwork{
					Ips: getIPv4(ni.IpAddress, regexp),
					Mac: wrapperspb.String(ni.MacAddress),
				})
			}
		}

		resp.VirtualMachines.VirtualMachines = append(resp.VirtualMachines.VirtualMachines, &machine)
	}

	return &resp, nil
}

// Получение доступных IP-адресов формата IPv4.
func getIPv4(ips []string, re *regexp.Regexp) []string {
	out := make([]string, 0, len(ips))
	for _, ip := range ips {
		if re.MatchString(ip) {
			out = append(out, ip)
		}
	}

	return out
}
