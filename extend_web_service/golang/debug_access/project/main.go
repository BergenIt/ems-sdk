package main

import (
	"context"
	"fmt"
	pb "service/gen/cluster-contract"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	pb.RegisterServiceManagerServer(server, &m)

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
	pb.UnimplementedServiceManagerServer
}

// RPC по сбору списка виртуальных машин с гипервизра ESXI.
func (r *microservice) DebugAccess(ctx context.Context, req *pb.DebugServiceAccessRequest) (*pb.DebugServiceAccessResponse, error) {
	// Составление адреса в зависимости от входных данных.
	address := ""
	switch f := req.Address.(type) {
	case *pb.DebugServiceAccessRequest_AddressPort:
		address = fmt.Sprintf("%s:%d", f.AddressPort.Address, f.AddressPort.Port)
	case *pb.DebugServiceAccessRequest_Uri:
		address = f.Uri
	}

	log.Printf("got request with address %s", address)

	var (
		err          error
		availability bool
		out          *pb.DebugServiceAccessResponse = &pb.DebugServiceAccessResponse{
			Result: &pb.DebugAccessResult{
				Address: address,
			},
		}
	)

	// Пинг адреса зависимости от протокола.
	switch req.Protocol {
	case pb.ServiceProtocol_SERVICE_PROTOCOL_GRPC,
		pb.ServiceProtocol_SERVICE_PROTOCOL_WS,
		pb.ServiceProtocol_SERVICE_PROTOCOL_TCP,
		pb.ServiceProtocol_SERVICE_PROTOCOL_HTTP:
		availability, err = pingTcp(address)

		out.Result.State = determineAvailability(availability, err)

	case pb.ServiceProtocol_SERVICE_PROTOCOL_UDP:
		availability, err = pingUdp(address)

		out.Result.State = determineAvailability(availability, err)

	default:
		return nil, fmt.Errorf("unsuppoted protocol: %v", req.Protocol)
	}

	return out, nil
}

// Определение статуса сервиса.
func determineAvailability(availability bool, err error) pb.ServiceAvailableState {
	if err == nil && availability {
		return pb.ServiceAvailableState_SERVICE_AVAILABLE_STATE_AVAILABLE
	} else if err != nil {
		log.Printf("determine service availability: %s", err)
		return pb.ServiceAvailableState_SERVICE_AVAILABLE_STATE_UNAVAILABLE
	} else {
		return pb.ServiceAvailableState_SERVICE_AVAILABLE_STATE_UNAVAILABLE
	}
}

func pingTcp(address string) (bool, error) {
	// Подкотовка адреса.
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return false, fmt.Errorf("resolve tcp addr: %s", err)
	}

	var (
		c net.Conn
		d *net.Dialer = &net.Dialer{
			Timeout: 2000 * time.Millisecond,
		}
	)

	c, err = d.Dial("tcp", addr.String())
	if err != nil {
		return false, nil
	}
	defer c.Close()

	return true, nil
}

// Так как для пинга юдп протокола необходимо отправить пакет и получить ответ, мы оптимистично полагаем, что если
// отправив пакет, произошел таймаут, значит порт открыт.
func pingUdp(address string) (bool, error) {
	// Подкотовка адреса.
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return false, fmt.Errorf("resolve udp addr: %s", err)
	}

	// Объявление интерфейса соединения.
	var (
		c net.Conn
		d *net.Dialer = &net.Dialer{
			Timeout: 2000 * time.Millisecond,
		}
	)

	c, err = d.Dial("udp", addr.String())
	if err != nil {
		return false, fmt.Errorf("dial: %s", err)
	}

	_, err = c.Write([]byte("Ping"))
	if err != nil {
		return false, fmt.Errorf("write: %s", err)
	}

	c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	defer c.Close()

	rb := make([]byte, 1500)

	if _, err := c.Read(rb); err != nil {
		if e := err.(*net.OpError).Timeout(); e {
			return true, nil
		}
		return false, fmt.Errorf("read: %s", err)
	}

	return true, nil
}
