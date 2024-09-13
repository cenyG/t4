package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"strings"

	"T4_test_case/config"
	"T4_test_case/internal/storageserver/pb"
	desc "T4_test_case/internal/storageserver/pb/chunkstorage"
	"T4_test_case/pkg/consul"
	"google.golang.org/grpc"
)

func main() {
	// validate config
	serviceName, port, healthCheckPort := config.Cfg.Storage.Name, config.Cfg.Storage.Port, config.Cfg.Common.HealthCheckPort

	parts := strings.Split(serviceName, ".")
	if len(parts) != 2 {
		log.Fatalf("wrong storage name: %s, must be like: 'storage.1'", config.Cfg.Storage.Name)
	}
	service, postfix := parts[0], parts[1]

	// register http health check for Consul
	go consul.RunHttpHealthCheck(healthCheckPort)

	// register service in Consul
	consulClient, err := consul.NewConsulWrapper()
	if err != nil {
		log.Fatalf("error connecting to Consul: %v", err)
	}

	err = consulClient.RegisterService(service, port, postfix, config.Cfg.Common.HealthCheckPort)
	if err != nil {
		log.Fatal(err)
	}

	// setup GRPC server
	address := fmt.Sprintf(":%s", config.Cfg.Storage.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to start: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterChunkStorageServer(s, pb.NewChunkStorageService())

	slog.Info(fmt.Sprintf("storage server %s listening on port: %s", serviceName, config.Cfg.Storage.Port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
