# Создание проекта

Для начала разработки собственных модулей расширения Bmc создадим скелет проекта

В будущем его можно использовать для обогащения операций необходимой логикой работы.

## Создание шаблона проекта

### Создайте проект из шаблона dotnet

Используя Visual Studio создадим проект из шаблона Dotnet gRPC API.

Подробнее:

* <https://learn.microsoft.com/en-us/aspnet/core/grpc/?view=aspnetcore-8.0>

Подключение модуля расширения к системе осуществляется с помощью следующих ключевых технологий:

- Протокол [gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
- [docker-compose](https://docs.docker.com/compose/)

Для работы сервиса необходимы протофайлы. Полный набор протофайлов можно найти в корне проекта `sdk`, в директории `.proto`.

Копия набора прото-файлов для RPC `CollectLedState` расположены в директории `Protos` [проекта](./project/).

Добавим ссылку на протофайлы в проект:

```xml
  <ItemGroup>
		<Protobuf ProtoRoot="./" Include="Protos/*.proto" AdditionalImportDirs="Protos/" OutputDir="$(IntermediateOutputPath)/%(RecursiveDir)" />
	</ItemGroup>
```

Реализуем простейший grpc сервис для подключенного протофайла:

```csharp
using Grpc.Core;

using ToolCluster.V4;

namespace BmcHandler.Services
{
    public class BmcHandlerService : BmcManager.BmcManagerBase
    {
        // RPC по сбору статуса LED.
        public override async Task<CollectBmcLedStateResponse> CollectLedState(CollectBmcLedStateRequest request, ServerCallContext context)
        {
            throw new RpcException(new Status(StatusCode.Unimplemented, ""));
        }
    }
}
```

Сконфигурируем приложение и grpc-сервер в `Program.cs`:

```csharp
using BmcHandler.Services;

var builder = WebApplication.CreateBuilder(args);

// Добавление gRPC функционала в контейнер.
builder.Services.AddGrpc();
// Добавление gRPC рефлексии в контейнер.
builder.Services.AddGrpcReflection();

var app = builder.Build();

// Настройка конвейера gRPC сервиса.
app.MapGrpcService<BmcHandlerService>();
// Сопоставление входящих запросов со службой рефлексии gRPC.
app.MapGrpcReflectionService();

await app.RunAsync();
```

На данном этапе пустой шаблон модуля расширения готов для локального запуска.

Для подключения модуля к системе необходимо настроить `Dockerfile` и `docker-compose.yaml`, подробную информацию об этом можно найти в директории `deploy`.

Пример готового шаблона модуля расширения находится в директории [project](./project).
