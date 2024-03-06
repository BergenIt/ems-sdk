# Создание проекта

## Создание шаблона проекта

Подключение модуля расширения к системе осуществляется с помощью следующих ключевых технологий:

- протокол [gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
- [docker-compose](https://docs.docker.com/compose/)

Для создания проекта на golang необходимо установить язык по следующей [инструкции](https://go.dev/doc/install).

Далее в директории проекта необходимо инициализировать файл `go.mod` для подгрузки внешних зависимостей с помощью команды `go mod init windows-handler`.

Далее необходимо [сгенерировать](#Подключение-протофайлов) код по необходимым протофайлам для работы gRPC сервера.

Для создания сервера `gRPC` в go нам необходимо:

1. Создать структуру, которая будет реализовывать интерфейс сервиса, описанного в протофайле, и сгенерированного в коде.

```golang
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
```

2. Создать инстанс структуры сервиса, создать сущность сервера gRPC, связать их, и запустить сервер.

```golang
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
```

На данном этапе пустой шаблон модуля расширения готов для локального запуска.

Для запуска рекомендуется использовать `Makefile` и утилиту `make`, создадим его и распишем набор команд:


```makefile
MODULE_NAME=windows-handler
gen:
	cd proto && \
	protoc --go_out=./.. \
	--go_opt=Mservice_windows_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_memory.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=module=${MODULE_NAME} \
	--go-grpc_out=./.. \
	--go-grpc_opt=Mservice_windows_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_memory.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=module=${MODULE_NAME} \
	service_windows_manager.proto shared_common.proto shared_device.proto shared_device_memory.proto

tidy:
	go mod tidy

build: gen tidy
	CGO_ENABLED=0 go build -o bin
```

С помощью команды `make build` можно сбилдить проект под работу в скретче.

Для подключения модуля к системе необходимо настроить `Dockerfile` и `docker-compose.yaml`, подробную информацию об этом можно найти в директории `deploy`.

## Подключение протофайлов

Для работы сервиса необходимо протофайлы. Полный набор протофайлов можно найти в корне проекта `sdk`, в директории `.proto`.

Ограниченный набор прото-файлов для RPC `CollectMemory` расположен в директории `proto` проекта `create_project`.

Из данных протофайлов необходимо сгенерировать код для корректной работы сервера, для этого рекомендуется использовать следующий набор утилит:

- protobuf-compiler
- protoc-gen-go
- protoc-gen-go-grpc

Рекомендуемые команды для устанвоки в Ubuntu 22.04:

```sh
apt-get update && \
apt install -y protobuf-compiler && \
apt clean
go env -w GOSUMDB=off
go env -w GO111MODULE=on
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 && \
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

Сгенерировать код из прото-файлов можно с помощью команды `make gen`.

Пример готового шаблона модуля расширения находится в директории `project`.
