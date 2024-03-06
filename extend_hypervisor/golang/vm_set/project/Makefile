MODULE_NAME=hypervisor
gen:
	cd proto && \
	protoc --go_out=./.. \
	--go_opt=Mservice_hypervisor_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_hypervisor.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_virtual_machine.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=module=${MODULE_NAME} \
	--go-grpc_out=./.. \
	--go-grpc_opt=Mservice_hypervisor_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_hypervisor.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_virtual_machine.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=module=${MODULE_NAME} \
	service_hypervisor_manager.proto shared_hypervisor.proto shared_virtual_machine.proto

tidy:
	go mod tidy

build: gen tidy
	CGO_ENABLED=0 go build -o bin