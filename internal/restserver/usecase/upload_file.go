package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"T4_test_case/internal/restserver/model"
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/services"
	"github.com/pkg/errors"
)

type UploadFileUseCase interface {
	Upload(ctx context.Context, reader io.Reader, filename string, size int64) (int64, error)
}

type uploadFileUseCase struct {
	storageServersProvider services.StorageServersProvider
	fileRepo               repo.FileRepository
}

func NewUploadFileUseCase(storageServersResolver services.StorageServersProvider, fileRepo repo.FileRepository) UploadFileUseCase {
	return &uploadFileUseCase{
		storageServersProvider: storageServersResolver,
		fileRepo:               fileRepo,
	}
}

func (u *uploadFileUseCase) Upload(ctx context.Context, reader io.Reader, filename string, size int64) (int64, error) {
	slog.Info("[uploadFileUseCase] start upload file")

	storageServerClients := u.storageServersProvider.GetStorageServersGrpcClients()
	storageClientsCount := int64(len(storageServerClients))
	if storageClientsCount == 0 {
		return 0, errors.New("no storageServerClients found")
	}

	chunkSize := size / storageClientsCount

	var storageServers []string
	// send chunks to storage servers
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
			return 0, errors.Wrap(err, "storageClient.UploadChunkStream")
		}
		storageServers = append(storageServers, storageClient.GetServiceID())
	}

	// save info about loaded file
	id, err := u.fileRepo.SaveFile(ctx, model.File{
		Name:    filename,
		Servers: strings.Join(storageServers, ","),
	})
	if err != nil {
		return 0, errors.Wrap(err, "u.fileRepo.SaveFile")
	}

	slog.Info(fmt.Sprintf("[uploadFileUseCase] file %s uploaded successfully", filename))

	return id, nil

}
