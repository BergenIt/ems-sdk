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