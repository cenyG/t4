package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"T4_test_case/config"
	"T4_test_case/internal/restserver/grpc"
	"T4_test_case/pkg/consul"
	"github.com/samber/lo"
)

type StorageServersProvider interface {
	GetStorageServersGrpcClients() []grpc.ChunkStorageClient
	GetStorageServersGrpcClientsMap() map[string]grpc.ChunkStorageClient
}

type storageServersProvider struct {
	consulWrapper consul.ConsulWrapper
	servers       concurrentMap
}

func NewStorageServersProvider(ctx context.Context) (StorageServersProvider, error) {
	consulWrapper, err := consul.NewConsulWrapper()
	if err != nil {
		return nil, err
	}

	ssr := storageServersProvider{
		consulWrapper: consulWrapper,
		servers: concurrentMap{
			m: make(map[string]grpc.ChunkStorageClient, 6),
		},
	}
	ssr.updateServersWorker(ctx)

	return &ssr, nil
}

func (sr *storageServersProvider) GetStorageServersGrpcClients() []grpc.ChunkStorageClient {
	return sr.servers.Values()
}

func (sr *storageServersProvider) GetStorageServersGrpcClientsMap() map[string]grpc.ChunkStorageClient {
	return sr.servers.Clone()
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
	services, err := sr.consulWrapper.GetServices("storage")
	if err != nil {
		slog.Error(fmt.Sprintf("[StorageServersProvider] sr.consulWrapper.Health().Service() error: %v", err))
		return
	}
	if len(services) == 0 {
		slog.Error("[StorageServersProvider] don't found storage servers")
	}

	for _, service := range services {
		id, port := service.Service.ID, service.Service.Port
		dockerName := strings.Replace(id, ".", "", -1)

		if _, ok := sr.servers.Get(id); !ok {
			slog.Info(fmt.Sprintf("[StorageServersProvider] fetch new storage server: %s - %s:%d", id, dockerName, port))

			client, cErr := grpc.NewChunkStorageClient(dockerName, port, id)
			if cErr != nil {
				slog.Error(fmt.Sprintf("[StorageServersProvider] pb.NewChunkStorageClient(%s, %d, %s) error: %v", dockerName, port, id, cErr))
				return
			}
			sr.servers.Set(id, client)
		}
	}

}

type concurrentMap struct {
	mu sync.RWMutex
	m  map[string]grpc.ChunkStorageClient
}

func (c *concurrentMap) Get(key string) (grpc.ChunkStorageClient, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}

func (c *concurrentMap) Values() []grpc.ChunkStorageClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return lo.Values(c.m)
}

func (c *concurrentMap) Clone() map[string]grpc.ChunkStorageClient {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// create new map with same values
	return lo.PickBy(c.m, func(key string, value grpc.ChunkStorageClient) bool {
		return true
	})
}

func (c *concurrentMap) Set(key string, val grpc.ChunkStorageClient) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = val
}
