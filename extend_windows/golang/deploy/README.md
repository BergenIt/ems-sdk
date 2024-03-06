# Развертывание модуля расширения

После написания кода приложения модуля расширения Windows необходимо позаботиться о том, чтобы его можно было развернуть в Docker.

Образ необходим для того, чтобы можно было запустить приложение в docker-compose на стенде.

После того, как приложение будет развернуто можно будет проверить работу модуля расширения Windows в UI EMS-a.

## Сборка в docker

В данной разделе будет описано каким образом будет формироваться Dockerfile, а также процесс сборки образа.

Для сборки приложения в контейнере докер необходимо сформировать Docker-образ приложения.

Для этого, как правильно, необходимо создать Dockerfile, который содержит информацию о базовом образе окружения, в котором будет работать контейнер и его приложение.

Подробнее об этом описано здесь:
 - https://docs.docker.com/reference/dockerfile/
 - https://docs.docker.com/language/golang/build-images/

Для возможности подключения модуля расширения в Dockerfile необходимо указать один из следующих лейблов:
 - **`ems.windows.memory.operation-system`** для операции по определенной операционной системе устройства или названием и версией операционной системы устройства (разделенных пробелом)
 - **`ems.windows.memory`** для операции без привязке к ОС

В случае нашего модуля расширения, Dockerfile для [проекта](../memory_inventory/project/) будет выглядеть следующим образом:

```docker
# Указание базового образа
FROM golang:1.22-bullseye AS build

# Указание лейбла операции
LABEL ems.windows.memory=default

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

Для развертывания используется Docker-compose, который позволяет конфигурировать запуск приложений.

Для запуска приложения в контейнере необходимо создать файл **`docker-compose.yaml`** и определить в нем следующие атрибуты:
 - Прописать имя нетворка, чтобы контейнеры были в одной сети `ems-network`
 - Имя контейнера **`windows-handler`**
 - Название образа, на основне которого будет запущен контейнер (или путь к докер-файлу). Укажем название образа, которое использовали на предыдущем этапе **`docker-windows-handler:latest`**, 
 - (необязательно) Порт работы gRPC-сервис **`8080`** на внешний порт системы **`55728`**, чтобы можно было отправлять запросы к приложению, запущенному в контейнере

Подробнее об этом описано здесь:
 - https://docs.docker.com/compose/
 - https://docs.docker.com/compose/compose-application-model/

В случае нашего модуля расширения файл docker-compose.yaml для [проекта](../memory_inventory/project/) будет выглядеть следующим образом:

```yaml
version: "3.9"

networks:
  default:
    name: 'ems-network'

services:
  windows-handler:
    image: docker-windows-handler:latest
    ports:
      - 55001:8080
```

Для того, чтобы поднять композицию достаточно перейти в директорию, где был создан файл **`docker-compose.yaml`** и в терминале ввести команду:
 - `docker compose up -d`

Если приложения корректно запущено, то при выполнении команды **`docker-compose.yaml`** мы увидим вот такой результат:

```table
CONTAINER ID   IMAGE                           COMMAND   CREATED          STATUS         PORTS                                       NAMES
69d9a7d31c20   docker-windows-handler:latest   "./bin"   37 seconds ago   Up 3 seconds   0.0.0.0:55001->8080/tcp, :::55001->8080/tcp   project-windows-handler-1
```

Проверить работу сервиса можно через Postman, подробнее об этом описано здесь:
 - https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/

Для проверки в UI EMS необходимо:
 - Авторизоваться в EMS
 - Завести оборудование с сетевым интерфейсом по WMI
 - Зайти в проводник
 - Найти заведенное устройство
 - Открыть карточку устройства
 - Перейти на вкладку **`Оперативная память`** (примечание: данные собираются раз в 15 минут, поэтому возможно придётся подождать)

Подробнее можно прочитать в руководстве пользователя
