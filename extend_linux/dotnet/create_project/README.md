# Создание проекта

Для начала разработки собственных модулей расширения template создадим скелет проекта

В будущем его можно использовать для обогащения операций необходимой логикой работы.

## Создание шаблона проекта

### Создайте проект из шаблона dotnet

Используя visual studio создадим проект из шаблона dotnet grpc api.

Подробнее:

* <https://learn.microsoft.com/en-us/aspnet/core/grpc/?view=aspnetcore-8.0>

Добавляем в папку Protos необходимые прото файлы, для примера возьмем:

```proto
//Тип используемого синтаксиса
syntax = "proto3";
//Пространство имен, которое будет использоваться с этим сервисом
option csharp_namespace = "SnmpExample";
//Название пакета
package example;
//Имя сервиса
service MyProtoService
{
//Название функции
  rpc SendPing (PingRequest) returns (PingReply);
}
//Класс передаваемый сервису
message PingRequest
{
  string requestString = 1;
}
//Класс получаемый в ответ
message PingReply
{
  string respondString = 1;
}
```

Добавим ссылку на протофайлы в проект:

```xml
 <ItemGroup>
  <Protobuf ProtoRoot="../" Include="Protos/*.proto" AdditionalImportDirs="Protos/" OutputDir="$(IntermediateOutputPath)/%(RecursiveDir)" />
 </ItemGroup>
```

Реализуем простейший grpc сервис для подключенного протофайла:

```csharp
using Grpc.Core;

namespace SnmpExample.Services
{
    // Наша реализация grpc сервиса
    public class MyService : MyProtoService.MyProtoServiceBase
    {
        // Наша реализация rpc SendPing
        public override Task<PingReply> SendPing(PingRequest request, ServerCallContext context)
        {
            return Task.FromResult(new PingReply { RespondString = "Respond string"});
        }
    }
}
```

Скофигурируем приложение и grpc-сервер в `Program.cs`:

```csharp
var builder = WebApplication.CreateBuilder(args);
//Добавление grpc функционала в конструтор
builder.Services.AddGrpc();

var app = builder.Build();
//Добавление сервиса в приложение
app.MapGrpcService<MyService>();
app.Run();
```

На данном этапе пустой шаблон модуля расширения готов для локального запуска.

Для подключения модуля к системе необходимо настроить `Dockerfile` и `docker-compose.yaml`, подробную информацию об этом можно найти в директории `deploy`.

Пример готового шаблона модуля расширения находится в директории [project](./project).
