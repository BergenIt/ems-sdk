# Создание проекта

Данный документ описывает создание минимально доступного рабочего скелета модуля расширения без реализации функционала RPC.

В будущем его можно использовать для обогащения операций необходимой логикой работы.

## Создание шаблона проекта

Подключение модуля расширения к системе осуществляется с помощью следующих ключевых технологий:

- [протокол gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
- [docker-compose](https://docs.docker.com/compose/)

Для создания проекта на golang необходимо установить sdk по следующей [инструкции](https://go.dev/doc/install).

Далее в директории проекта необходимо инициализировать файл `go.mod` для подгрузки внешних зависимостей с помощью команды `go mod init <название_корневого_модуля_по_усмотрению>`.

## Подключение протофайлов

Для работы сервиса необходимы протофайлы. Полный набор протофайлов можно найти в корне проекта `sdk`, в директории `.proto`.

Копия набора прото-файлов для RPC `PutSettings` расположены в директории `proto` [проекта](./project/).

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
    MODULE_NAME=<название_корневого_модуля_по_усмотрению>
    gen:
        cd proto && \
        protoc --go_out=./.. \
        --go_opt=Mservice_sso_center.proto=${MODULE_NAME}/gen/cluster-contract \
        --go_opt=module=${MODULE_NAME} \
        --go-grpc_out=./.. \
        --go-grpc_opt=Mservice_sso_center.proto=${MODULE_NAME}/gen/cluster-contract \
        --go-grpc_opt=module=${MODULE_NAME} \
        service_sso_center.proto

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
        pb.UnimplementedSsoCenterServer
    }

    // RPC для установления настроек LDAP-авторизации на BMC.
    func (r *microservice) PutSettings(context.Context, *pb.PutSsoSettingsRequest) (*pb.PutSsoSettingsResponse, error) {
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
    pb.RegisterSsoCenterServer(server, &m)

    // Создаем листененра.
    lis, err := net.Listen("tcp", listenPort)
    if err != nil {
        return fmt.Errorf("create listener: %s", err)
    }

    // Запускаем gRPC сервер.
    return server.Serve(lis)
}
```

На данном этапе пустой шаблон модуля расширения готов для локального запуска.

Для подключения модуля к системе необходимо настроить `Dockerfile` и `docker-compose.yaml`, подробную информацию об этом можно найти в директории `deploy`.

Пример готового шаблона модуля расширения находится в директории `project`.
