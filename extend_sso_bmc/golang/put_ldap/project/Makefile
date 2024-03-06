MODULE_NAME=sso_center
gen:
	cd proto && \
	protoc --go_out=./.. \
	--go_opt=Mservice_sso_center.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=module=${MODULE_NAME} \
	--go-grpc_out=./.. \
	--go-grpc_opt=Mservice_sso_center.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=module=${MODULE_NAME} \
	service_sso_center.proto

tidy:
	go mod tidy

build: gen tidy
	CGO_ENABLED=0 go build -o bin