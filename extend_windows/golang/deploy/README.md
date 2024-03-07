# Развертывание модуля расширения

После реализации модуля расширения Windows необходимо позаботиться о том, чтобы его можно было развернуть в Docker.

Образ необходим для того, чтобы можно было запустить приложение в docker-compose на стенде.

После того, как приложение будет развернуто можно будет проверить работу сбора инвентарных данных ОЗУ в UI EMS-a.

## Сборка в docker

Для сборки приложения в контейнере докер необходимо сформировать Docker-образ приложения.

Для этого, как правильно, необходимо создать Dockerfile, который содержит информацию о базовом образе окружения, в котором будет работать контейнер и его приложение.

Подробнее об этом описано здесь:

- <https://docs.docker.com/reference/dockerfile/>
- <https://docs.docker.com/language/golang/build-images/>

Важно для корректной работы модуля расширения прописать лейблы в Dockerfile.

В случае нашего модуля расширения, Dockerfile для [проекта](../memory_inventory/project/) будет выглядеть следующим образом:

```dockerfile
# Указание базового образа
FROM golang:1.22-bullseye AS build

### Установка компилятора protobuf
RUN apt-get update && \
    apt install -y protobuf-compiler && \
    apt clean
###

### Установка Golang Env-s
RUN go env -w GOSUMDB=off
RUN go env -w GO111MODULE=on
###

### Установка пакетов для генерация protobuf-файлов
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
###

WORKDIR /src

# Копирование исходных файлов в образ решения
COPY . .
# Запуск Makefile для компиляции файлов protobuf и бинарного файла приложения
RUN make build

# Использование многоэтапной сборки для уменьшения размера образа
FROM scratch AS release
# Указание лейбла операции
LABEL ems.windows.memory=default
WORKDIR /app
# Копирование бинарного файла приложения, полученного на предыдущем этапе сборки
COPY --from=build /src/bin ./bin
ENTRYPOINT ["./bin"]
```

После формирования Dockerfile можно сформировать сам образ, для этого в терминале необходимо выполнить команду и дождаться завершения процесса сборки:

- `docker build --tag docker-windows-handler .`

После формирования образа в терминале можно ввести команду `docker images` и увидеть, что образ успешно сформирован.

```table
REPOSITORY                TAG       IMAGE ID       CREATED          SIZE
docker-windows-handler    latest    33f51e228eb4   17 seconds ago   17.4MB
```

## Развертывание

В процессе работы система EMS ищет docker сервис с лейблом "ems.windows.memory.operation-system" и значением, которое совпадает с названием операционной системы устройства или которое совпадает с названием и версией операционной системы устройства (разделенных пробелом).

Eсли не находит, то ищет сервис с лейблом "ems.windows.memory", если не находит, то ищет сервис с лейблом "ems.windows".

Список лейблов с описанием:

Один на выбор:

- LABEL ems.windows.memory.operation-system
  - Для привязки модуля расширения к операции для конкретных операционных систем устройств.
  - Для этого объявите этот лейбл с названием и версией операционной системы устройства, разделенных пробелом, в качестве значения лейбла.
- LABEL ems.windows.memory
  - Для привязки модуля расширения к операции для всех устройств.
  - Для этого объявите этот лейбл с любым значением.

Дополнительные:

- LABEL ems.service.secure=default - активация защищенного соединения, в случае использования tls на gRPC сервере.
- LABEL ems.service.port=8081 - порт, который прослушивает сервер.
- LABEL ems.grpc-service.healthcheck=81 - порт для проверки доступности сервера.

Лейблы необходимо указывать в Dockerfile проекта.

Далее необходимо создать `docker-compose.yaml` файл конфигурации для развертывания сервиса через `docker-compose`.

Важно определить в конфигурации работу в сети `ems-network`, пример рабочего файла конфигурации выглядит следующим образом:

```yaml
version: "3.9"

networks:
  default:
    name: 'ems-network'

services:
  windows-handler:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 55001:8080
```

Подробнее об этом описано здесь:

- <https://docs.docker.com/compose/>
- <https://docs.docker.com/compose/compose-application-model/>

Для запуска модуля расширения в эксплуатацию необходимо поместить папку `project` на ВМ, где запущен EMS и выполнить команду `docker-compose up --build -d` из корня проекта.

Если приложения корректно запущено, то при выполнении команды `docker ps | grep project-windows-handler-1` мы увидим вот такой результат:

```bash
CONTAINER ID   IMAGE                           COMMAND   CREATED          STATUS         PORTS                                         NAMES
69d9a7d31c20   docker-windows-handler:latest   "./bin"   37 seconds ago   Up 3 seconds   0.0.0.0:55001->8080/tcp, :::55001->8080/tcp   project-windows-handler-1
```

Проверить работу сервиса можно через Postman, подробнее об этом описано здесь:

- <https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/>

Для проверки в UI EMS необходимо:

- Авторизоваться в EMS
- Завести оборудование с сетевым интерфейсом по WMI
- Зайти в проводник
- Найти заведенное устройство
- Открыть карточку устройства
- Перейти на вкладку **`Оперативная память`** (примечание: данные собираются раз в 15 минут, поэтому возможно придётся подождать)

Подробнее можно прочитать в руководстве пользователя
