package services

import (
	"T4_test_case/config"
	"T4_test_case/internal/restserver/model"
	"T4_test_case/internal/restserver/pb"
	"context"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/samber/lo"
	"log/slog"
	"time"
)

type StorageServersProvider interface {
	GetStorageServersGrpcClients() []pb.ChunkStorageClient
}

type storageServersProvider struct {
	consulClient *consulapi.Client
	servers      map[model.StorageServerAddress]pb.ChunkStorageClient
}

func NewStorageServersProvider(ctx context.Context) (StorageServersProvider, error) {
	consulCfg := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	ssr := storageServersProvider{
		consulClient: client,
		servers:      make(map[model.StorageServerAddress]pb.ChunkStorageClient, 6),
	}
	ssr.updateServersWorker(ctx)

	return &ssr, nil
}

func (sr *storageServersProvider) GetStorageServersGrpcClients() []pb.ChunkStorageClient {
	return lo.Values(sr.servers)
}

func (sr *storageServersProvider) updateServersWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				sr.updateServer()
			}

			time.Sleep(config.Cfg.UpdateStorageServersTimeout)
		}
	}()
}

func (sr *storageServersProvider) updateServer() {
	services, _, err := sr.consulClient.Health().Service("storage", "", true, nil)
	if err != nil {
		slog.Error("sr.consulClient.Health().Service() error: %v", err)
		return
	}

	for _, service := range services {
		address := model.StorageServerAddress{
			Port: service.Service.Port,
			Host: service.Service.Address,
		}

		if _, ok := sr.servers[address]; !ok {
			slog.Info("fetch new storage server: %s", address)

			client, cErr := pb.NewChunkStorageClient(address.Host, address.Port)
			if cErr != nil {
				slog.Error("pb.NewChunkStorageClient(%s, %d) error: %v", address.Host, address.Port, cErr)
				return
			}
			sr.servers[address] = client
		}
	}

}
