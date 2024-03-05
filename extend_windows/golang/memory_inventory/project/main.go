package main

import (
	"log"

	grpc_server "windows-handler/grpc"
)

func main() {
	if err := grpc_server.Serve(); err != nil {
		log.Fatalf("run: %s", err)
	}
}
