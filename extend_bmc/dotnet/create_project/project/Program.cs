using BmcHandler.Services;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddGrpc();
// Add gRPC reflection to the container.
builder.Services.AddGrpcReflection();

var app = builder.Build();

// Configure the HTTP request pipeline.
app.MapGrpcService<BmcHandlerService>();
// Map incoming requests to the gRPC reflection service.
app.MapGrpcReflectionService();

await app.RunAsync();
