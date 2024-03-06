MODULE_NAME=network

all: build

build: gen tidy
	CGO_ENABLED=0 go build -o bin

gen:
	cd proto && \
	protoc --go_out=./.. \
	--go_opt=Mservice_network_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=module=${MODULE_NAME} \
	--go-grpc_out=./.. \
	--go-grpc_opt=Mservice_network_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=module=${MODULE_NAME} \
	service_network_manager.proto shared_common.proto shared_device.proto

tidy:
	go mod tidy

clean:
	rm -rf bin gen
