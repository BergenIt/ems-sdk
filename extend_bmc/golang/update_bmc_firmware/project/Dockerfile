# Указание базового образа
FROM golang:1.22-bullseye AS build

### Установка компилятора protobuf
RUN apt-get update && \
    apt install -y protobuf-compiler && \
    apt clean
###

### Установка Golang Env-s
RUN go env -w GOSUMDB=off
RUN go env -w GO111MODULE=on
###

### Установка пакетов для генерация protobuf-файлов
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
###

WORKDIR /src

# Копирование исходных файлов в образ решения
COPY . .
# Запуск Makefile для компиляции файлов protobuf и бинарного файла приложения
RUN make build

# Использование многоэтапной сборки для уменьшения размера образа
FROM scratch AS release
# Указание лейбла операции
LABEL ems.bmc.bmc-firmware-update=default
WORKDIR /app
# Копирование бинарного файла приложения, полученного на предыдущем этапе сборки
COPY --from=build /src/bin ./bin
ENTRYPOINT ["./bin"]