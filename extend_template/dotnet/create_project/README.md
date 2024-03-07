# Создание проекта

Для начала разработки собственных модулей расширения Template создадим скелет проекта.

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

### Реализация отправки Snmp запросов
Большинство операций данного сервиса используют Snmp протокол для общения.

Для отправки Snmp запросов будем использовать библиотеку `Lextm.SharpSnmpLib` - [github](https://github.com/lextudio/sharpsnmplib?tab=readme-ov-file)

Создадим класс c данными необходимыми для подключения:

```csharp
public record SnmpCredential(
    string Ip
    string Login,
    int Port,
    int Version,
    string SecurityName,
    string SecurityLevel,
    string? Community,
    string? Context,
    string? AuthProtocol,
    string? AuthKey,
    string? PrivateProtocol,
    string? PrivateKey);
```

Реализуем отправку запросов:

```csharp
string SendRequest(SnmpCredential credential, string oidTemplate, int port, int timeout)
{
    string result;

    // Переменные используемые библиотекой Lextm.SharpSnmpLib
    IPEndPoint endpoint = new(IPAddress.Parse(credential.Ip), port);
    OctetString community = new(credential.Community);
    ObjectIdentifier oid = new(oidTemplate);
    VersionCode versionCode = (credential.Version == 1 || credential.Version == 2) ? VersionCode.V2 : VersionCode.V3;

    // Попытка получить данные по конкретному пути (реализация функции Get)
    // Например - При запросе 1.2.3.4 вернет 1.2.3.4
    string resultGet = "Null";

    try
    {
        GetRequestMessage message = new(0, versionCode, community, new List<Variable> { new(oid) });

        ISnmpMessage response = message.GetResponse(timeout, endpoint);

        if (response.Pdu().ErrorStatus.ToInt32() == 0)
        {
            resultGet = response.Pdu().Variables.FirstOrDefault().Data.ToString();
        }
    }
    catch (Exception ex)
    {
        result = "Error: " + ex.Message;
    }

    //Возвращаем если получен валидный ответ
    if (resultGet != "NoSuchObject" && resultGet != "Null")
    {
        return resultGet;
    }

    // Если по конкретному пути получить данные не вышло пробуем получить следующий по списку (реализация функции Walk)
    // Например - При запросе 1.2.3.4 вернет 1.2.3.4 или 1.2.3.4.0
    
    List<Variable> resultGetBulk = new();
    
    try
    {
        GetBulkRequestMessage message = new(0, versionCode, community, 0, 1, new List<Variable> { new(oid) });

        ISnmpMessage response = message.GetResponse(timeout, endpoint);
        if (response.Pdu().ErrorStatus.ToInt32() == 0)
        {
            resultGetBulk = response.Pdu().Variables.ToList();
        }
    }
    catch (Exception ex)
    {
        return "Error: " + ex.Message;
    }

    //При возврате пустого массива возвращаем ошибку
    return resultGetBulk.Count == 0 ? "Error: No valid oid" : resultGetBulk.First().Data.ToString();
}
```

На данном этапе пустой шаблон модуля расширения готов для локального запуска.

Пример готового шаблона модуля расширения находится в директории [project](./project).
