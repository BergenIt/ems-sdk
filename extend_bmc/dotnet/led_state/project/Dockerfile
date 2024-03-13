FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
LABEL ems.bmc.led-state=default
USER app
WORKDIR /app
EXPOSE 8080

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /src
COPY ["BmcHandler.csproj", "."]
RUN dotnet restore "./BmcHandler.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "./BmcHandler.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "./BmcHandler.csproj" -c Release -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "BmcHandler.dll"]
