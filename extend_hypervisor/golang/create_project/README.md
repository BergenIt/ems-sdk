# Создание проекта

## Создание шаблона проекта

Подключение модуля расширения к системе осуществляется с помощью следующих ключевых технологий:

- протокол [gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
- [docker-compose](https://docs.docker.com/compose/)

Для создания проекта на golang необходимо установить язык по следующей [инструкции](https://go.dev/doc/install).

Далее в директории проекта необходимо инициализировать файл `go.mod` для подгрузки внешних зависимостей с помощью команды `go mod init <название_корневого_модуля_по_усмотрению>`.

Далее необходимо [сгенерировать](#Подключение-протофайлов) код по необходимым протофайлам для работы gRPC сервера.

Для создания сервера `gRPC` в go нам необходимо:

1. Создать структуру, которая будет реализовывать интерфейс сервиса, описанного в протофайле, и сгенерированного в коде.

```golang
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
```

2. Создать инстанс структуры сервиса, создать сущность сервера gRPC, связать их, и запустить сервер.

```golang
    func run() error {
        // Создаем инстанс сервиса.
        m := microservice{}

        // Создаем инстанс сервера.
        server := grpc.NewServer()

        // Регистрируем сервис.
        pb.RegisterHypervisorManagerServer(server, &m)

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

## Подключение протофайлов


> Описываем зачем тут протофайлы

> Описываем как подключить к проекту протофайлы

> Описываем зачем мы их подключили и что можем с ними делать

> Указываем что пример итогового проекта находится в папке project
