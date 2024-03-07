MODULE_NAME=bmc-handler
gen:
	cd proto && \
	protoc --go_out=./.. \
	--go_opt=Mservice_bmc_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_available.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_cpu.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_cpu_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_template.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_disk.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_initial.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_memory.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_memory_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_operation_system.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_pci_slot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_uptime.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_boot_source.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_event.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_firmware_boot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_firmware.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_ipmi.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_led.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_power_state.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_power_usage_limit.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_power_usage.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_redfish.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=Mshared_device_temperature.proto=${MODULE_NAME}/gen/cluster-contract \
	--go_opt=module=${MODULE_NAME} \
	--go-grpc_out=./.. \
	--go-grpc_opt=Mservice_bmc_manager.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_common.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_available.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_cpu.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_cpu_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_template.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_disk.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_initial.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_memory.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_memory_utilization.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_operation_system.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_pci_slot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_uptime.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_boot_source.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_event.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_firmware_boot.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_firmware.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_ipmi.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_led.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_power_state.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_power_usage_limit.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_power_usage.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_redfish.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=Mshared_device_temperature.proto=${MODULE_NAME}/gen/cluster-contract \
	--go-grpc_opt=module=${MODULE_NAME} \
	service_bmc_manager.proto shared_common.proto \
	shared_device.proto shared_device_available.proto shared_device_cpu.proto \
	shared_device_cpu_utilization.proto shared_device_template.proto \
	shared_device_disk.proto shared_device_initial.proto shared_device_memory.proto \
	shared_device_memory_utilization.proto shared_device_operation_system.proto \
	shared_device_pci_slot.proto shared_device_uptime.proto \
	shared_device_boot_source.proto \
	shared_device_event.proto \
	shared_device_firmware_boot.proto \
	shared_device_firmware.proto \
	shared_device_ipmi.proto \
	shared_device_led.proto \
	shared_device_power_state.proto \
	shared_device_power_usage_limit.proto \
	shared_device_power_usage.proto \
	shared_device_redfish.proto \
	shared_device_temperature.proto

tidy:
	go mod tidy

build: gen tidy
	CGO_ENABLED=0 go build -o bin