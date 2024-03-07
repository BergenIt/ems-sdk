# Развертывание модуля расширения

После реализации модуля расширения из примера необходимо позаботиться о том, чтобы его можно было развернуть в Docker.

Образ необходим для того, чтобы можно было запустить приложение в docker-compose на стенде.

После того, как приложение будет развернуто можно будет проверить работу сбора данных о нагрузке на CPU .

## Сборка в docker

Для сборки docker образа создадим в проекте `Dockerfile`

Подробнее об этом описано здесь:

- <https://docs.docker.com/reference/dockerfile/>
- <https://docs.docker.com/language/dotnet/>
- <https://learn.microsoft.com/ru-ru/visualstudio/containers/container-build?view=vs-2022>

Для возможности подключения модуля расширения к системе в `Dockerfile` необходимо указать лейбл.

Для модуля расширения Template подходят следующие лейблы:
* Для привязки модуля расширения ко всем операциям и всем устройствам - `ems.template`
* Для привязки модуля к конкретной системной операции к общему лейблу добавляется:
    * Получение данных для заведения устройства `.device-initial`
    * Получение данных об операционной системе `.operation-system`
    * Получение данных доступности `.available`
    * Получение инвентарных данных `.inventory`
    * Получение данных об энергопотреблении `.power-usage`
	* Получение данных о времени работы устройства `.uptime`
	* Получение данных о процессорах `.cpu`
	* Получение данных о использовании оперативной памяти `.memory-utilization`
	* Получение данных о нагрузки на процессор `.cpu-utilization`
	* Получение данных о дисках `.disk`
	* Получение данных о устройствах в PCI слотах `.pci-slot`
	* Получение метаданных сенсоров дисков `.smart-sensor-meta`
	* Получение инвентарных данных с дисков `.smart-inventory-meta`
	* Установка агентов `.agent-set`
	* Перезагрузка устройства `.reboot`
	* Выключение устройства `.power-off`
	* Вызов скрипта `.script-invoke`
	* Установка приложения при помощи скрипта `.script-soft-install`
    
* Для привязки модуля расширения к конкретной модели устройства к лейблу системной операции добавляется `.model`
* Для привязки модуля расширения к конкретному производителю устройства к лейблу системной операции добавляется `.vendor`
* Для привязки модуля расширения к конкретной операционной системой устройства к лейблу системной операции добавляется `.operation-system`

Например для операции получение температуры подходят следущие лейблы:

- **`ems.template.temperature.model`**
    - Для привязки модуля расширения к операции для конкретных моделей устройств.
- **`ems.template.temperature.vendor`**
    - Для привязки модуля расширения к операции для конкретных производителям устройств.
- **`ems.template.temperature.operation-system`**
    - Для привязки модуля расширения к операции для конкретным операционным системам.
- **`ems.template.temperature`**
    - Для привязки модуля расширения к операции для всех устройств.
- **`ems.template`**
    - Для привязки модуля расширения ко всем устройствам

В различных случаях докер файлы моглядеть по разному. 
Для того, чтобы развернуть приложение из примера нужно следующее:

*	Базовый контейнер с .net SDK версии 8.0
*	В конитейнере должна быть установлена zlib1g-dev для обработки прото файлов
*	Исходные коды должны быть скопированы в контейнер, а затем скомпилированы для среды linux-x64
*	Для уменьшения размеров финального образа, можно использовать двухэтапную сборку, тогда понадобится еще один базовый образ с Runtime версией .net 8.0, в который нужно скопировать
бинарные файлы, а затем запустить их.

Можно использовать следующий докер файл:
	
	FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

	RUN apt-get update \
		&& apt-get install -y --no-install-recommends \
		clang zlib1g-dev

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


Для развертывания микросервисов в ЕМС используется Docker-compose, который позволяет конфигурировать запуск приложений.  
Для запуска приложения в контейнере необходимо создать файл docker-compose.yaml и прописать в нем атрибуты запуска приложения.

Подробнее об этом описано здесь:

https://docs.docker.com/compose/
https://docs.docker.com/compose/compose-application-model/


Для приложения из примера docker-compose.yml можно использовать самый простейший docker-compose. В передаче переменных среды в контейнер надобности нет.
Нужно просто указать путь к докер образу и к докер файлу.

Вполне достаточно будет вот такого файла:
	
	version: '3.4'

	services:
	  linux-manager:
		image: linuxmanager:latest
		build:
		  context: .
		  dockerfile: Dockerfile

Для того, чтобы поднять композицию достаточно перейти в директорию, где был создан файл docker-compose.yaml и в терминале ввести команду:

	docker compose up -d

## Проверка в системе

Развертывание проекта должно происходить на виртуальной машине с работоспособным EMS.

После завершения разработки поместите получившийся проект на стенд и запустите созданный ранее `docker-compose`.

Для проверки получения нагрузки на процессоры по протоколу SSH необходимо:

1) Авторизоваться в EMS
2) Завести оборудование или настроить заведенное введя данные для подключения по SSH
3) Подождать немного (пока произведется фоновый опрос оборудования)
4) Полученная температура должна отобразиться на вкладке Процессоры

Подробнее можно прочитать в руководстве пользователя
