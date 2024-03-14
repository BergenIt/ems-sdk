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

FROM scratch AS release
LABEL ems.sso.put-setting=default
WORKDIR /app
COPY --from=build /src/bin ./bin
ENTRYPOINT ["./bin"]
```

Подробнее здесь:

- <https://docs.docker.com/reference/dockerfile/>
- <https://docs.docker.com/language/golang/build-images/>

Сбилдить проект можно с помощью команды `docker build .`

Запустить решение в контейнере можно с помощью команды `docker run -dp 127.0.0.1:8080:8080 <id_образа>`.

## Развертывание

В процессе работы система EMS ищет docker сервис с лейблом "ems.sso.put-setting.model" и значением, которое совпадает с названием модели устройства или docker сервис с лейблом "ems.sso.put-setting.vendor" значение которого совпадает с названием производителем устройства. Eсли не находит (или устройство не было найдено по Ip), то ищет сервис с лейблом "ems.sso.put-setting", если не находит, то ищет сервис с лейблом "ems.sso".

Список лейблов с описанием:

Один на выбор:

- LABEL ems.sso.put-setting.vendor
  - Для привязки модуля расширения к операции для конкретных вендоров устройств.
  - Для этого объявите этот лейбл с моделью устройства.
- LABEL ems.sso.put-setting.model
  - Для привязки модуля расширения к операции для конкретных моделей устройств.
  - Для этого объявите этот лейбл с моделью устройства.
- LABEL ems.sso.put-setting
  - Для привязки модуля расширения к операции для всех устройств.
  - Для этого объявите этот лейбл с любым значением.

Дополнительные:

- LABEL ems.service.secure=default - активация защищенного соединения.
- LABEL ems.service.port=8081 - порт, который прослушивает сервер.
- LABEL ems.grpc-service.healthcheck=81 - порт для проверки доступности сервера.

Лейблы необходимо указывать в Dockerfile проекта.

Далее необходимо создать `docker-compose.yaml` файл конфигурации для развертывания сервиса через `docker-compose`.

Важно прописать extra_hosts следующим образом:

```yaml
extra_hosts:
  - traefik:<IP master-cluster>
```

Важно определить в конфигурации работу в сети `ems-network`, пример рабочего файла конфигурации выглядит следующим образом:

```yaml
version: "3.9"

networks:
  default:
    name: 'ems-network'

services:
  sso-center-put-ldap-handler-huawei:
    extra_hosts:
      - traefik:10.1.18.17
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
    hostname: sso-center-put-ldap-handler-huawei
    environment:
      ServicePort: :8080
    deploy:
      resources:
        limits:
          cpus: "3"
          memory: 2000M
        reservations:
          cpus: "0.5"
          memory: 400M
```

Подробнее об этом здесь:

- <https://docs.docker.com/compose/>
- <https://docs.docker.com/compose/compose-application-model/>

Для запуска модуля расширения в эксплуатацию необходимо поместить папку `project` на ВМ, где запущен EMS и выполнить команду `docker-compose up --build -d` из корня проекта.

Если приложения корректно запущено, то при выполнении команды `docker ps | grep project-sso-center-put-ldap-handler-1` мы увидим вот такой результат:

```bash
CONTAINER ID   IMAGE                     COMMAND   CREATED         STATUS         PORTS     NAMES
059ed982d404   project-sso-center-put-ldap-handler   "./bin"   3 seconds ago   Up 2 seconds             project-sso-center-put-ldap-handler-1
```

Проверить работу сервиса можно через Postman, подробнее об этом описано [здесь](https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/).

Для проверки в UI EMS необходимо:

- Авторизоваться в EMS
- Зайти во вкладку управление
- Выбрать операцию "Определить EMS, как центр идентификации пользователей BMC"
- Выполнить операцию на необходимом оборудовании

Подробнее можно прочитать в руководстве пользователя.
