# Развертывание модуля расширения

## Сборка в docker
Для работы микросервисов в EMS необходимо собирать их в Docker-образ.

Для этого необходимо:
1) Добавить в проект файл `Dockerfile`
Для проекта примера файл выглядит так:
```s
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
USER app
WORKDIR /app
EXPOSE 8080

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["SnmpExample.csproj", "."]
RUN dotnet restore "./SnmpExample.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "./SnmpExample.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "./SnmpExample.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "SnmpExample.dll"]
```
2) Команда для сборки образа
```s
docker build --tag {Название ораза} {Путь папки проекта}
```
Пример выполнения команды из папки проекта
```s
docker build --tag example ./
```
3) Посмотреть созданные образы можно вызвав команду
```s
docker images
```
Пример выполнения команды
```s
REPOSITORY                                  TAG                          IMAGE ID       CREATED          SIZE
example                                     latest                       b6cac163644f   26 minutes ago   219MB
```

Подробные инструкции по сборке и настройке контейнеров docker:
https://learn.microsoft.com/ru-ru/visualstudio/containers/container-build?view=vs-2022
https://docs.docker.com/language/dotnet/containerize/

## Развертывание

Для развертывания контейнера в EMS используется Docker-compose, который позволяет настроить межсервисное общение.

Для этого необходимо создать файл `docker-compose.yml` и определить в нем следующие атрибуты:
```s
networks:
	default:
		name: 'ems-network'  //Имя сети, необходимо чтобы контейнеры были в одной сети

services:
  snmp-example:   //Название вашего сервиса
    image: example:latest  //Название созданного вами образа
    ports:
      - 7777:8080  //Порты для обращения к сервису внешний:внутренний
```
2) Для того, чтобы создать контейнер достаточно перейти в директорию, где был создан файл docker-compose.yaml и в терминале ввести команду:
```s
docker compose up -d
```
3) Посмотреть созданные контейнеры можно вызвав команду
```s
docker container ls
```
Пример выполнения команды
```s
CONTAINER ID   IMAGE            COMMAND                  CREATED              STATUS              PORTS                    NAMES
122ffc95ec59   example:latest   "dotnet SnmpExample.…"   About a minute ago   Up About a minute   0.0.0.0:7777->8080/tcp   1-snmp-example-1
```
## Локальный запуск
После сборки контейнер автоматически запустится. Вы можете отправлять к нему запросу на localhost:{Порты указанные в docker-compose.yaml}.
Для остановки контейнера необходимо выполнить команду:
```s
docker stop {Имя контейнера}
```
Для повторного запуска не нужно снова собирать контейнер, можно выполнить команду
```s
docker start {Имя контейнера}
```
## Проверка в системе
После завершения разработки поместите получившийся проект на стенд.
Как только код пройдет ci/cd и стенд перезапустится, можно проверять в системе.

Для проверки получения температуры по протоколу SNMP необходимо:
1) Авторизоваться в EMS
2) Добавить или настроить существующий шаблон для сбора температуры
    * Шаблон мониторинга -> Системные метрики -> Датчик температур -> SNMP
3) Завести оборудование или настроить заведенное введя данные для подключения по SNMP
4) Подождать немного (пока произведется фоновый опрос оборудования)
5) Полученная температура должна отобразиться на вкладке Общая информация

Подробнее можно прочитать в руководстве пользователя

## Пример итогового проекта находится в папке project