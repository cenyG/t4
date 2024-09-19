package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"T4_test_case/internal/restserver/grpc"
	"T4_test_case/internal/restserver/model"
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/services"
	desc "T4_test_case/pb/chunkstorage"
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

// Upload - read file from use Reader and proxy chunks to storage servers
func (u *uploadFileUseCase) Upload(ctx context.Context, reader io.Reader, filename string, size int64) (int64, error) {
	slog.Info("[UploadFileUseCase] start upload file")

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
			// if error occurs delete already uploaded files from servers to optimize storage
			go u.deleteChunks(storageServerClients[:i+1], filename)
			return 0, errors.Wrap(err, "storageClient.UploadChunkStream")
		}
		storageServers = append(storageServers, storageClient.ServiceID())
	}

	// save info about loaded file
	id, err := u.fileRepo.SaveFile(ctx, model.File{
		Name:    filename,
		Servers: strings.Join(storageServers, ","),
	})
	if err != nil {
		// may be delete files from storages
		return 0, errors.Wrap(err, "u.fileRepo.SaveFile")
	}

	slog.Info(fmt.Sprintf("[UploadFileUseCase] file %s uploaded successfully", filename))

	return id, nil
}

// deleteChunks - remove file from all servers if fail to proceed
func (u *uploadFileUseCase) deleteChunks(clients []grpc.ChunkStorageClient, filename string) {
	for index, client := range clients {
		_, err := client.DeleteChunk(context.Background(), &desc.DeleteChunkRequest{
			FileName:   filename,
			ChunkIndex: int32(index),
		})
		if err != nil {
			slog.Error(fmt.Sprintf("[UploadFileUseCase] file %s delete error: %v", filename, err))
		}
	}
}
