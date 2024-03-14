FROM golang:1.22-bullseye AS build

RUN apt-get update && \
    apt install -y protobuf-compiler && \
    apt clean
RUN go env -w GOSUMDB=off
RUN go env -w GO111MODULE=on
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /src

COPY . .
RUN make build

FROM debian AS release
LABEL ems.network.create-config=default
WORKDIR /app
COPY --from=build /src/bin ./bin
ENTRYPOINT ["./bin"]
