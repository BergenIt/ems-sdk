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
