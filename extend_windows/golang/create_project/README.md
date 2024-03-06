# Создание проекта

## Создание шаблона проекта

Процесс создание проекта будет описан на примере редактора кода Visual Studio Code.

 - В Visual Studio Code откройте папку, в которой вы создадите корневой каталог приложения Go. Чтобы открыть папку, щелкните значок Обозреватель на панели действий и нажмите кнопку "Открыть папку".

 - Щелкните **"Создать папку"** на панели Обозреватель, а затем создайте корневой директор для примера приложения Go с именемsample-app

 - Нажмите кнопку **"Создать файл"** на панели Обозреватель, а затем назовите файл **`main.go`**

 - Нажмите кнопку **"Создать файл"** на панели Обозреватель, а затем назовите файл **`Makefile`** для генерации Proto-файлов исходного кода и бинарного файла приложения

 - Скопируйте следующие инструкции в **`Makefile`** файл:

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

 - Откройте терминал, **"Терминал"** -> **"Создать треминал"**, а затем выполните команду **`go mod init sample-app`**, чтобы инициализировать пример приложения Go.

 - Скопируйте следующий код в **`main.go`** файл:

```go
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

type microservice struct {
	pb.UnimplementedWindowsManagerServer
}

func (r *microservice) CollectMemory(
	ctx context.Context,
	req *pb.CollectWindowsMemoryRequest,
) (*pb.CollectWindowsMemoryResponse, error) {
	//реализация rpc
	//...

	return nil, errors.New("not implemented")
}

```

Ориентир gRPC-операций и сервиса идёт на конкретный модуль расширения, а именно сбор информации с Windows.

Подробнее описано здесь:
 - https://learn.microsoft.com/ru-ru/azure/developer/go/configure-visual-studio-code

## Подключение протофайлов

Протофайлы необходимы для определения gRPC-сервиса и контрактов (RPC), на основе которых можно взаимодействовать с сервисом.

Варианты подключения протофайлов к проекту:
 - Создать отдельную папку в корне проекта и поместить туда файлы с расширением **`.proto`**
 
 - Добавить через сабмодули проект в Gitlab, который содержит файлы с расширением **`.proto`**
   Пример команды для добавления сабмодуля: **`git submodule add https://gitlab.com/my_project/proto_example`**

После подключения прототофайлов разработчик сервиса может пользоваться контрактами и разрабатывать RPC.

Подробнее можно прочитать здесь:
 - https://grpc.io/docs/what-is-grpc/introduction/

Пример готового проекта расположен в папке [project](./project)
