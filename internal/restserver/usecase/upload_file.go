package usecase

import (
	"T4_test_case/internal/restserver/services"
	"context"
	"github.com/pkg/errors"
	"io"
	"log/slog"
)

type UploadFileUseCase interface {
	Upload(ctx context.Context, reader io.Reader, filename string, size int64) error
}

type uploadFileUseCase struct {
	storageServersProvider services.StorageServersProvider
}

func NewUploadFileUseCase(storageServersResolver services.StorageServersProvider) UploadFileUseCase {
	return &uploadFileUseCase{
		storageServersProvider: storageServersResolver,
	}
}

func (u *uploadFileUseCase) Upload(ctx context.Context, reader io.Reader, filename string, size int64) error {
	storageServerClients := u.storageServersProvider.GetStorageServersGrpcClients()
	storageClientsCount := int64(len(storageServerClients))
	if storageClientsCount == 0 {
		return errors.New("no storage storageServerClients found")
	}

	chunkSize := size / storageClientsCount

	for i, storageClient := range storageServerClients {
		var idx = int64(i)
		var currentChunkSize int64
		if idx == storageClientsCount-1 {
			// send all remaining data to last client
			currentChunkSize = size - (chunkSize * idx)
		} else {
			currentChunkSize = chunkSize
		}

		limitedReader := io.LimitReader(reader, currentChunkSize)

		err := storageClient.UploadChunkStream(ctx, filename, int32(i), limitedReader, currentChunkSize)
		if err != nil {
			return errors.Wrap(err, "storageClient.UploadChunkStream")
		}
	}

	slog.Info("file %s uploaded successfully", filename)

	return nil

}
