package main

import (
	"T4_test_case/pkg/interceptors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"strings"

	"T4_test_case/config"
	"T4_test_case/internal/storageserver/service"
	desc "T4_test_case/pb/chunkstorage"
	"T4_test_case/pkg/consul"
	"google.golang.org/grpc"
)

func main() {
	config.Init()
	// validate config
	serviceId, port, healthCheckPort := config.Cfg.Storage.Name, config.Cfg.Storage.Port, config.Cfg.Common.HealthCheckPort

	parts := strings.Split(serviceId, ".")
	if len(parts) != 2 {
		log.Fatalf("wrong storage name: %s, must be like: 'storage.1'", config.Cfg.Storage.Name)
	}
	serviceName, postfix := parts[0], parts[1]

	// register http health check for Consul
	go consul.RunHttpHealthCheck(healthCheckPort)

	// register service in Consul
	consulClient, err := consul.NewConsulWrapper()
	if err != nil {
		log.Fatalf("error connecting to Consul: %v", err)
	}

	err = consulClient.RegisterService(serviceName, port, postfix, config.Cfg.Common.HealthCheckPort)
	if err != nil {
		log.Fatal(err)
	}

	// setup GRPC server
	address := fmt.Sprintf(":%s", config.Cfg.Storage.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to start: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryPanicRecoveryInterceptor()),
		grpc.StreamInterceptor(interceptors.StreamPanicRecoveryInterceptor()),
	)
	desc.RegisterChunkStorageServer(s, service.NewChunkStorageService())

	slog.Info(fmt.Sprintf("storage server %s listening on port: %s", serviceId, config.Cfg.Storage.Port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
