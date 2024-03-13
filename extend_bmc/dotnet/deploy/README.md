# Развертывание модуля расширения

После реализации модуля расширения Bmc необходимо позаботиться о том, чтобы его можно было развернуть в Docker.

Образ необходим для того, чтобы можно было запустить приложение в docker-compose на стенде.

После того, как приложение будет развернуто можно будет проверить работу системной операции в UI EMS-a.

## Сборка в docker

Для сборки docker образа создадим в проекте `Dockerfile`

Подробнее об этом описано здесь:

- <https://docs.docker.com/reference/dockerfile/>
- <https://docs.docker.com/language/dotnet/>
- <https://learn.microsoft.com/ru-ru/visualstudio/containers/container-build?view=vs-2022>

`Dockerfile` проекта-примера:

```Dockerfile
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
LABEL ems.bmc.led-state=default
USER app
WORKDIR /app
EXPOSE 8080

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /src
COPY ["BmcHandler.csproj", "."]
RUN dotnet restore "./BmcHandler.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "./BmcHandler.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "./BmcHandler.csproj" -c Release -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "BmcHandler.dll"]
```

Выполним сборку docker образа:

```bash
docker build --tag bmc-handler-led ./
```

## Развертывание

В процессе работы система EMS ищет docker сервис с лейблом "ems.bmc.led-state.model" и значением,
которое совпадает с названием модели устройства или docker сервис с лейблом "ems.bmc.led-state.vendor" значение которого совпадает с названием производителем устройства.

Eсли не находит, то ищет сервис с лейблом "ems.bmc.led-state", если не находит, то ищет сервис с лейблом "ems.bmc".

Список лейблов с описанием:

Один на выбор:

- LABEL ems.bmc.led-state.model
  - Для привязки модуля расширения к операции для конкретной модели устройства.
  - Для этого объявите этот лейбл с названием модели устройства в качестве значения лейбла.
- LABEL ems.bmc.led-state.vendor
  - Для привязки модуля расширения к операции для производителя устройства.
  - Для этого объявите этот лейбл с названием производителя устройства в качестве значения лейбла.
- LABEL ems.bmc.led-state
  - Для привязки модуля расширения к операции для всех устройств.
  - Для этого объявите этот лейбл с любым значением.

Дополнительные:

- LABEL ems.service.secure=default - активация защищенного соединения, в случае использования tls на gRPC сервере.
- LABEL ems.service.port=8081 - порт, который прослушивает сервер.
- LABEL ems.grpc-service.healthcheck=81 - порт для проверки доступности сервера.

Лейблы необходимо указывать в Dockerfile проекта.

Для развертывания контейнера в EMS используется Docker-compose, который позволяет настроить межсервисное общение.

Для этого необходимо создать файл `docker-compose.yml` и определить в нем следующие атрибуты:

```yaml
networks:
 default:
  # Имя сети, необходимо чтобы контейнеры были в одной сети
  name: 'ems-network'

services:
  # Название вашего сервиса
  bmc-handler-led:
    # Процесс сборки через Dockerfile
    build:
      context: .
      dockerfile: Dockerfile
    ports:
    # Порты для обращения к сервису внешний:внутренний
      - 55555:8080
```

Для запуска проекта в директории с файлом docker-compose.yaml необходимо выполнить команду:

```bash
docker compose up -d
```

Убедитесь в том, что контейнер запущен:

```bash
docker ps -a | grep 'bmc-handler-led'
```

Пример вывода:

```bash
CONTAINER ID   IMAGE            COMMAND                  CREATED              STATUS              PORTS                    NAMES
122ffc95ec59   bmc-handler-led:latest   "dotnet BmcHandler.…"   About a minute ago   Up About a minute   0.0.0.0:42763->8080/tcp   project-bmc-handler-led-1
```

## Проверка в EMS

Развертывание проекта должно происходить на виртуальной машине с работоспособным EMS.

После завершения разработки поместите получившийся проект на стенд и запустите созданный ранее `docker-compose`.

Для чистоты проверки в UI EMS необходимо остановить Docker-контейнер стандартного модуля расширения по BMC, а именно контейнер `ems-bmc-manager-1`.

Алгоритм проверки в UI EMS:

- Авторизоваться в EMS
- Завести оборудование с сетевым интерфейсом по Redfish
- Зайти в раздел управлении
- Открыть задачу по Включению/выключению Led
- Отобразится список устройств со статусом Led

Подробнее можно прочитать в руководстве пользователя
