# Развертывание модуля расширения

После реализации модуля расширения template необходимо позаботиться о том, чтобы его можно было развернуть в Docker.

Образ необходим для того, чтобы можно было запустить приложение в docker-compose на стенде.

После того, как приложение будет развернуто можно будет проверить работу сбора инвентарных данных ОЗУ в UI EMS-a.

## Сборка в docker

Для сборки docker образа создадим в проекте `Dockerfile`

Подробнее об этом описано здесь:

- <https://docs.docker.com/reference/dockerfile/>
- <https://docs.docker.com/language/dotnet/>
- <https://learn.microsoft.com/ru-ru/visualstudio/containers/container-build?view=vs-2022>

`Dockerfile` проекта-примера:

```Dockerfile
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
USER app
WORKDIR /app
EXPOSE 8080

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["BmcManager.csproj", "."]
RUN dotnet restore "./BmcManager.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "./BmcManager.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "./BmcManager.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "BmcManager.dll"]
```

Выполним сборку docker образа:

```bash
docker build --tag example ./
```

## Развертывание

Для развертывания контейнера в EMS используется Docker-compose, который позволяет настроить межсервисное общение.

Для этого необходимо создать файл `docker-compose.yml` и определить в нем следующие атрибуты:

```yaml
networks:
 default:
  # Имя сети, необходимо чтобы контейнеры были в одной сети
  name: 'ems-network'

services:
  # Название вашего сервиса
  snmp-example:
   # Название созданного вами образа
    image: example:latest
    ports:
   # Порты для обращения к сервису внешний:внутренний
      - 42763:8080
```

Для запуска проекта в директории с файлом docker-compose.yaml необходимо выполнить команду:

```bash
docker compose up -d
```

Убедитесь в том, что контейнер запущен:

```bash
docker ps -a | grep 'snmp-example'
```

Пример вывода:

```bash
CONTAINER ID   IMAGE            COMMAND                  CREATED              STATUS              PORTS                    NAMES
122ffc95ec59   example:latest   "dotnet BmcManager.…"   About a minute ago   Up About a minute   0.0.0.0:42763->8080/tcp   1-snmp-example-1
```

## Проверка в системе

Развертывание проекта должно происходить на виртуальной машине с работоспособным EMS.

После завершения разработки поместите получившийся проект на стенд и запустите созданный ранее `docker-compose`.

Для проверки получения температуры по протоколу SNMP необходимо:

1) Авторизоваться в EMS
2) Добавить или настроить существующий шаблон для сбора температуры
    - Шаблон мониторинга -> Системные метрики -> Датчик температур -> SNMP
3) Завести оборудование или настроить заведенное введя данные для подключения по SNMP
4) Подождать немного (пока произведется фоновый опрос оборудования)
5) Полученная температура должна отобразиться на вкладке Общая информация

Подробнее можно прочитать в руководстве пользователя
