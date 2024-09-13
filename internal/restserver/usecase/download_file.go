package usecase

import (
	"T4_test_case/internal/restserver/model"
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/services"
	"context"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type DownloadFileUseCase interface {
	Download(ctx context.Context, file *model.File, writer io.Writer) error
}

type downloadFileUseCase struct {
	storageServersProvider services.StorageServersProvider
	fileRepo               repo.FileRepository
}

func NewDownloadFileUseCase(storageServersProvider services.StorageServersProvider, fileRepo repo.FileRepository) DownloadFileUseCase {
	return &downloadFileUseCase{
		storageServersProvider: storageServersProvider,
		fileRepo:               fileRepo,
	}
}

func (d *downloadFileUseCase) Download(ctx context.Context, file *model.File, writer io.Writer) error {
	storageServers := strings.Split(file.Servers, ",")
	storageServersMap := d.storageServersProvider.GetStorageServersGrpcClientsMap()

	for _, storageServer := range storageServers {
		if _, ok := storageServersMap[storageServer]; !ok {
			return errors.Errorf("storage server %s not exist founded", storageServer)
		}
	}

	var chunkIndex int32
	for _, storageServer := range storageServers {
		chunkStorageClient := storageServersMap[storageServer]
		err := chunkStorageClient.DownloadChunkStream(ctx, file.Name, chunkIndex, writer)
		if err != nil {
			return errors.Wrap(err, "chunkStorageClient.DownloadChunkStream")
		}

		chunkIndex += 1
	}

	return nil
}
