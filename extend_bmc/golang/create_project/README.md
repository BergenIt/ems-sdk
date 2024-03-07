# Создание проекта

Для начала разработки собственных модулей расширения Bmc создадим скелет проекта

В будущем его можно использовать для обогащения операций необходимой логикой работы.

## Создание шаблона проекта

Для создания проекта на golang необходимо установить sdk по следующей [инструкции](https://go.dev/doc/install).

Подключение модуля расширения к системе осуществляется с помощью следующих ключевых технологий:

- Протокол [gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
- [docker-compose](https://docs.docker.com/compose/)

В директории проекта необходимо инициализировать файл `go.mod` для подгрузки внешних зависимостей с помощью команды `go mod init bmc-handler`. Подробнее в документации [golang](https://go.dev/doc/tutorial/create-module).

## Подключение протофайлов

Для работы сервиса необходимы протофайлы. Полный набор протофайлов можно найти в корне проекта `sdk`, в директории `.proto`.

Копия набора прото-файлов для RPC `BmcFirmwareUpdate` расположены в директории `proto` [проекта](./project/).

Из данных протофайлов необходимо сгенерировать код для корректной работы сервера, для этого рекомендуется использовать следующий набор утилит:

- [protobuf-compiler](https://grpc.io/docs/protoc-installation/)
- [protoc-gen-go](https://pkg.go.dev/github.com/golang/protobuf/protoc-gen-go)
- [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)

Для установки в Ubuntu 22.04:

```sh
apt-get update && \
apt install -y protobuf-compiler && \
apt clean
go env -w GOSUMDB=off
go env -w GO111MODULE=on
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 && \
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

Рекомендуется использовать `Makefile` и утилиту `make` для работы с проектом, создадим его и распишем набор команд:

```makefile
MODULE_NAME=bmc-handler
gen:
	cd proto && \
	protoc --go_out=./.. \
	--go_opt=Mservice_bmc_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_available.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_cpu.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_cpu_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_template.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_disk.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_initial.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_memory.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_memory_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_operation_system.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_pci_slot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_uptime.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_boot_source.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_event.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_firmware_boot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_firmware.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_ipmi.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_led.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_power_state.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_power_usage_limit.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_power_usage.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_redfish.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_temperature.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=module=${MODULE_NAME} \
	--go-grpc_out=./.. \
	--go-grpc_opt=Mservice_bmc_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_available.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_cpu.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_cpu_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_template.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_disk.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_initial.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_memory.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_memory_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_operation_system.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_pci_slot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_uptime.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_boot_source.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_event.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_firmware_boot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_firmware.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_ipmi.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_led.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_power_state.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_power_usage_limit.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_power_usage.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_redfish.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_temperature.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=module=${MODULE_NAME} \
	service_bmc_manager.proto shared_common.proto \
	shared_device.proto shared_device_available.proto shared_device_cpu.proto \
	shared_device_cpu_utilization.proto shared_device_template.proto \
	shared_device_disk.proto shared_device_initial.proto shared_device_memory.proto \
	shared_device_memory_utilization.proto shared_device_operation_system.proto \
	shared_device_pci_slot.proto shared_device_uptime.proto \
	shared_device_boot_source.proto \
	shared_device_event.proto \
	shared_device_firmware_boot.proto \
	shared_device_firmware.proto \
	shared_device_ipmi.proto \
	shared_device_led.proto \
	shared_device_power_state.proto \
	shared_device_power_usage_limit.proto \
	shared_device_power_usage.proto \
	shared_device_redfish.proto \
	shared_device_temperature.proto

tidy:
	go mod tidy

build: gen tidy
	CGO_ENABLED=0 go build -o bin
```

Сгенерировать код из прото-файлов можно с помощью команды `make gen`.

С помощью команды `make build` можно сбилдить проект.

## Создание gRPC сервера

Для создания сервера `gRPC` в go нам необходимо:

1. Создать структуру, которая будет реализовывать интерфейс сервиса, описанного в протофайле, и сгенерированного в коде.

```golang
// Инстанс сервиса с реализацией RPC.
type microservice struct {
	pb.UnimplementedBmcManagerServer
}

// RPC по обновлению прошивки BMC.
func (r *microservice) BmcFirmwareUpdate(
	ctx context.Context,
	req *pb.BmcFirmwareUpdateRequest,
) (*pb.BmcFirmwareUpdateResponse, error) {
	//реализация rpc
	//...

	return nil, errors.New("not implemented")
}
```

2. Создать инстанс структуры сервиса, создать сущность сервера gRPC, связать их, и запустить сервер.

```golang
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

	log.Printf("microservice start serving on port %q", listenPort)

	// Запускаем gRPC сервер.
	return server.Serve(lis)
}
```

На данном этапе пустой шаблон модуля расширения готов для локального запуска.

Для подключения модуля к системе необходимо настроить `Dockerfile` и `docker-compose.yaml`, подробную информацию об этом можно найти в директории `deploy`.

Пример готового шаблона модуля расширения находится в директории `project`.
