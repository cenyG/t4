package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"strings"

	"T4_test_case/internal/restserver/model"
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/services"
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

// Download - download file chunks from store servers and proxy them to writer
func (d *downloadFileUseCase) Download(ctx context.Context, file *model.File, writer io.Writer) error {
	slog.Info("[uploadFileUseCase] start download file")

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

	slog.Info(fmt.Sprintf("[uploadFileUseCase] download file %s success", file.Name))

	return nil
}
