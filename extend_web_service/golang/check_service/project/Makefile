MODULE_NAME=service_manager
gen:
	cd proto && \
	protoc --go_out=./.. \
	--go_opt=Mservice_service_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_service_instance.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=module=${MODULE_NAME} \
	--go-grpc_out=./.. \
	--go-grpc_opt=Mservice_service_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_service_instance.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=module=${MODULE_NAME} \
	service_service_manager.proto shared_service_instance.proto

tidy:
	go mod tidy

build: gen tidy
	CGO_ENABLED=0 go build -o bin