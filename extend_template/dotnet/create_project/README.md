# Создание проекта

## Создание шаблона проекта

1) Создайте проект из шаблона "Служба ASP.NET Core gRPC"
    * Введите название для проекта (в моем примере проект будет называться SnmpExample)
    * Выберите платформу .NET 8.0
    *Если пакет Grpc.AspNetCore не был установлен автоматически - установите его вручную, тк данный пакет необходим для работы приложения!*
2) Автоматически будут созданы следующие папки и файлы:
    * Папка Properties содержит файл launchSettings.json, который определяет параметры запуска сервиса.
    * Папка Protos содержит файлы с определением сервисов и сообщений, используемых для взаимодействия по сети (По умолчанию в этой папке определен один файл greet.proto.)
    * Папка Services содержит файлы с реализацией сервисов (По умолчанию в этой папке определен один файл GreeterService.cs.)
    * Файл appsettings.json - стандартный файл конфигурации приложения ASP.NET Core.
    * Файл appsettings.Development.json - файл конфигурации приложения для стадии разработки.
    * Файл Program.cs содержит стандартный класс Program, с которого начинается выполнение приложения ASP.NET Core.
    * Файл SnmpExample.csproj - стандартный файл конфигурации проекта C#.

3) gRPC использует подход "contract-first", то есть вначале определяется контракт - общее определение сервиса, которое определяет механизм взаимодействия.
Добавляем в папку Protos файл my_proto.proto:
```s
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
4) Необходимо добавить ссылку на протофайлы в проекте. Для этого добавляем в файл SnmpExample.csproj следующие строки между тегами <Project...> и </Project>
```s
	<ItemGroup>
		<Protobuf ProtoRoot="../" Include="Protos/*.proto" AdditionalImportDirs="Protos/" OutputDir="$(IntermediateOutputPath)/%(RecursiveDir)" />
	</ItemGroup>
```
5) Добавляем в папку Services файл MyService.cs:
```s
//Необходимая для работы библиотека из пакета Grpc.AspNetCore
using Grpc.Core;
//Пространство имен
namespace SnmpExample.Services
{
  //Объявление сервиса grpc
    public class MyService : MyProtoService.MyProtoServiceBase
    {
      //Описание действий и результата в ответ на полученный запрос SendPing
        public override Task<PingReply> SendPing(PingRequest request, ServerCallContext context)
        {
          //Создается новый объект PingReply и передается в ответ
            return Task.FromResult(new PingReply { RespondString = "Respond string"});
        }
    }
}
```
6) В файле Program.cs пишем:
```s
//Пространство имен
using SnmpExample.Services;
//Передача аргументов конструктору
var builder = WebApplication.CreateBuilder(args);
//Добавление grpc функционала в конструтор
builder.Services.AddGrpc();
//Создание приложения с помощью конструктора
var app = builder.Build();
//Добавление сервиса в приложение
app.MapGrpcService<MyService>();
//Запуск приложения
app.Run();
```

## Локальный запуск
После сборки и запуска приложения, откроется консоль, в которой будет написаны адреса доступные для отправки запросов приложению. Grpc использует защищенное соединение, поэтому адрес запущенного сервиса будет начинаться с https://

## Пример итогового проекта находится в папке project
