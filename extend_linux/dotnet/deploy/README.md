# Развертывание модуля расширения

## Сборка в docker

Микросервисы EMS запускаются в среде Docker.

*	Инструкция по сборке образов docker: https://docs.docker.com/language/dotnet/containerize/

В различных случаях докер файлы моглядеть по разному. Для того, чтобы развернуть приложение из примера
можно использовать следующий:
	
	FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

	RUN apt-get update \
		&& apt-get install -y --no-install-recommends \
		clang zlib1g-dev

	ARG CLIENT__KEY_BUILD
	ARG CLIENT__PUBLIC_KEY_BUILD

	WORKDIR "/app"

	COPY "./SshExample/SshExample.csproj" "./SshExample/SshExample.csproj"
	RUN dotnet restore "./SshExample/SshExample.csproj" -r linux-x64

	COPY . .
	RUN dotnet build "./SshExample/SshExample.csproj" -c Release -o /app/build
	RUN dotnet publish "./SshExample/SshExample.csproj" \
		--maxcpucount \
		--no-restore \
		-c Release \
		-o /app/publish \
		-r linux-x64 \
		-p:TrimmerRemoveSymbols=true \
		-p:InvariantGlobalization=true \
		-p:UseSystemResourceKeys=true \
		-p:HttpActivityPropagationSupport=false \
		-p:EnableCompressionInSingleFile=true \
		-p:EnableUnsafeBinaryFormatterSerialization=false \
		-p:EnableUnsafeUTF7Encoding=false \
		-p:IncludeNativeLibrariesForSelfExtract=true \
		-p:IncludeAllContentForSelfExtract=true \
		-p:DebugType=none


	FROM  mcr.microsoft.com/dotnet/runtime:8.0 AS final
	LABEL ems.linux=default
	WORKDIR /app
	EXPOSE 8080
	COPY --from=build /app/publish .
	ENTRYPOINT ["./SshExample"]


Докер файл нужно создать в директории с проектом, рядом с файлом SshExample.sln.  
После этого выполнить команду

*	docker build --tag docker-linux-manager .
## Развертывание


Для развертывания используется Docker-compose, который позволяет конфигурировать запуск приложений.  
Для запуска приложения в контейнере необходимо создать файл docker-compose.yaml и прописать в нем атрибуты запуска приложения.

Подробнее об этом описано здесь:

https://docs.docker.com/compose/
https://docs.docker.com/compose/compose-application-model/

Для приложения из примера docker-compose.yml будет выглядить следующим образом:
	
	version: '3.4'

	services:
	  linux-manager:
		image: linuxmanager:latest
		environment:
		  
		  Kestrel__EndPoints__Http__Url: http://*:4545
		  Kestrel__EndPoints__Http__Protocols: Http1

		build:
		  context: .
		  dockerfile: Dockerfile

Для того, чтобы поднять композицию достаточно перейти в директорию, где был создан файл docker-compose.yaml и в терминале ввести команду:

	docker compose up -d

Если приложения корректно запущено, то при выполнении команды docker-compose.yaml мы увидим вот такой результат: