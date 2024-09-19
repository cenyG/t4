package services

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"T4_test_case/config"
	"T4_test_case/internal/restserver/grpc"
	"T4_test_case/internal/restserver/model"
	desc "T4_test_case/pb/chunkstorage"
	"T4_test_case/pkg/collections"
	"T4_test_case/pkg/consul"
	"github.com/samber/lo"
)

type StorageServersProvider interface {
	GetStorageServersGrpcClients() []grpc.ChunkStorageClient
	GetStorageServersGrpcClientsMap() map[string]grpc.ChunkStorageClient
}

type storageServersProvider struct {
	consulWrapper consul.ConsulWrapper
	serversMap    collections.ConcurrentMap[string, grpc.ChunkStorageClient]
	serversStats  collections.ConcurrentArray[model.ServerInfo]
}

func NewStorageServersProvider(ctx context.Context) (StorageServersProvider, error) {
	consulWrapper, err := consul.NewConsulWrapper()
	if err != nil {
		return nil, err
	}

	ssr := storageServersProvider{
		consulWrapper: consulWrapper,
		serversMap:    collections.NewConcurrentMap[string, grpc.ChunkStorageClient](6),
		serversStats:  collections.NewConcurrentArray[model.ServerInfo](6),
	}
	ssr.updateServersWorker(ctx)

	return &ssr, nil
}

// GetStorageServersGrpcClients - get first N servers with the highest available storage
func (s *storageServersProvider) GetStorageServersGrpcClients() []grpc.ChunkStorageClient {
	serversMap := s.serversMap.Clone()

	infos := s.serversStats.Values()[0:config.Cfg.StorageServersCount]
	infos = lo.Shuffle(infos) // shuffle to randomize serversMap order for more load balancing

	res := make([]grpc.ChunkStorageClient, len(infos))
	for _, info := range infos {
		res = append(res, serversMap[info.ServerID])
	}

	return s.serversMap.Values()
}

func (s *storageServersProvider) GetStorageServersGrpcClientsMap() map[string]grpc.ChunkStorageClient {
	return s.serversMap.Clone()
}

func (s *storageServersProvider) updateServersWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				s.updateServersConsul()
				s.updateServersStats(ctx)
			}

			time.Sleep(config.Cfg.UpdateStorageServersInterval)
		}
	}()
}

// updateServersConsul - fetch health serversMap from Consul and update only if found new server
func (s *storageServersProvider) updateServersConsul() {
	services, err := s.consulWrapper.GetServices("storage")
	if err != nil {
		slog.Error(fmt.Sprintf("[StorageServersProvider] s.consulWrapper.Health().Service() error: %v", err))
		return
	}
	if len(services) == 0 {
		slog.Error("[StorageServersProvider] don't found storage serversMap")
	}

	for _, service := range services {
		id, port := service.Service.ID, service.Service.Port
		dockerName := strings.Replace(id, ".", "", -1)

		if _, ok := s.serversMap.Get(id); !ok {
			slog.Info(fmt.Sprintf("[StorageServersProvider] fetch new storage server: %s - %s:%d", id, dockerName, port))

			client, cErr := grpc.NewChunkStorageClient(dockerName, port, id)
			if cErr != nil {
				slog.Error(fmt.Sprintf("[StorageServersProvider] pb.NewChunkStorageClient(%s, %d, %s) error: %v", dockerName, port, id, cErr))
				return
			}
			s.serversMap.Set(id, client)
		}
	}
}

// updateServersStorage - update serversMap info from grpc clients and write sorted by ServerInfo.Avail array
func (s *storageServersProvider) updateServersStats(ctx context.Context) {
	serversMap := s.serversMap.Clone()
	serversStats := make([]model.ServerInfo, len(serversMap))

	for _, server := range serversMap {
		resp, err := server.ServerStats(ctx, &desc.ServerStatsRequest{})
		if err != nil {
			slog.Error(fmt.Sprintf("[StorageServersProvider] error while get stats for server %s: %v", server.ServiceID(), err))
			continue
		}

		serversStats = append(serversStats, model.ServerInfo{
			ServerID: server.ServiceID(),
			Total:    resp.DiscTotal,
			Avail:    resp.DiscAvail,
			Used:     resp.DiscUsed,
		})
	}

	sort.Slice(serversStats, func(i, j int) bool {
		return serversStats[i].Avail > serversStats[j].Avail
	})

	s.serversStats.SetAll(serversStats)
}
