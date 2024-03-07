# Развертывание модуля расширения

Данный документ описывает подход к развертыванию модуля расширения с помощью `docker` и `docker-compose`.

## Сборка в docker

Для сборки проекта необходим [docker](https://docs.docker.com/build/building/packaging/) и [docker-compose](https://docs.docker.com/compose/).

Dockerfile рекомендуется расположить в корне проекта.

Важно для корректной работы модуля расширения прописать лейблы в Dockerfile.

Пример рабочего докерфайла:

```dockerfile
FROM golang:1.22-bullseye AS build

RUN apt-get update && \
    apt install -y protobuf-compiler && \
    apt clean
RUN go env -w GOSUMDB=off
RUN go env -w GO111MODULE=on
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /src

COPY . .
RUN make build

FROM debian AS release
WORKDIR /app
COPY --from=build /src/bin ./bin
ENTRYPOINT ["./bin"]
```

Сбилдить проект можно с помощью команды `docker build .`

Запустить решение в контейнере можно с помощью команды `docker run -dp 127.0.0.1:8080:8080 <id_образа>`.

## Развертывание

В процессе работы система EMS ищет docker сервис с лейблом "ems.network.create-config.vendor" и значением, которое совпадает
с названием операционной системы устройства или которое совпадает с названием и версией операционной системы устройства (разделенных пробелом).

Eсли не находит, то ищет сервис с лейблом "ems.network.create-config.vendor", если не находит, то ищет сервис с лейблом "ems.network".

Список лейблов с описанием:

- LABEL ems.network.create-config.vendor=default - включение модуля расширения в обработку системой EMS.
- LABEL ems.service.secure=default - активация защищенного соединения.
- LABEL ems.service.port=8081 - порт, который прослушивает сервер.
- LABEL ems.grpc-service.healthcheck=81 - порт для проверки доступности сервера.

Лейблы необходимо указывать в Dockerfile проекта.

Далее необходимо создать `docker-compose.yaml` файл конфигурации для развертывания сервиса через `docker-compose`.

Важно определить в конфигурации работу в сети `ems-network`, пример рабочего файла конфигурации выглядит следующим образом:

```
version: "3.9"

networks:
  default:
    name: 'ems-network'

services:
  network-manager:
    build:
      context: .
      dockerfile: Dockerfile
    restart: unless-stopped
    logging:
      options:
        max-size: '50M'
        max-file: '5'
    ulimits:
      core:
        hard: 0
        soft: 0
    hostname: network-manager
    environment:
      ServicePort: 8081
    deploy:
      resources:
        limits:
          cpus: "3"
          memory: 2000M
        reservations:
          cpus: "0.5"
          memory: 400M
```


Для запуска модуля расширения в эксплуатацию необходимо поместить папку `project` на ВМ, где запущен EMS и выполнить команду `docker-compose up
--build -d` из корня проекта.

Если приложения корректно запущено, то при выполнении команды `docker ps | grep project-network-manager-1` мы увидим вот такой результат:

```table
CONTAINER ID   IMAGE                     COMMAND   CREATED         STATUS         PORTS     NAMES
059ed982d404   project-network-manager   "./bin"   3 seconds ago   Up 2 seconds             project-network-manager-1
```

Проверить работу сервиса можно через Postman, подробнее об этом описано [здесь](https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/).
