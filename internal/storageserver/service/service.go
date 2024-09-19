package service

import (
	"T4_test_case/config"
	desc "T4_test_case/pb/chunkstorage"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ricochet2200/go-disk-usage/du"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	_filesPath = "/app/files"
)

type Implementation struct {
	desc.UnimplementedChunkStorageServer
}

func NewChunkStorageService() *Implementation {
	return &Implementation{}
}

// UploadChunk - endpoint read file chunk from stream and save to file
func (i *Implementation) UploadChunk(stream desc.ChunkStorage_UploadChunkServer) error {
	slog.Info(fmt.Sprintf("[Upload] start upload chunk to %s", config.Cfg.Storage.Name))

	var file *os.File
	for {
		// read data from steam
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileName, chunkIndex := chunk.GetFileName(), chunk.GetChunkIndex()

		// create file
		if file == nil {
			err = os.MkdirAll(_filesPath, os.ModePerm)
			if err != nil {
				slog.Info(fmt.Sprintf("[Upload] os.MkdirAll err: %v", err))
				return err
			}
			filePath := getFilePath(fileName, chunkIndex)
			file, err = os.Create(filePath)
			if err != nil {
				slog.Info(fmt.Sprintf("[Upload] os.Create err: %v", err))
				return err
			}
			defer file.Close()
		}

		// Write chunk to file
		_, err = file.Write(chunk.GetChunkData())
		if err != nil {
			slog.Info(fmt.Sprintf("[Upload] file.Write err: %v", err))
			return err
		}
	}

	slog.Info(fmt.Sprintf("[Upload] success chunk uploaded to %s", config.Cfg.Storage.Name))

	return stream.SendAndClose(&desc.UploadChunkResponse{})
}

// DownloadChunk - endpoint read file chunk from local storage and send to stream
func (i *Implementation) DownloadChunk(req *desc.DownloadChunkRequest, stream desc.ChunkStorage_DownloadChunkServer) error {
	slog.Info(fmt.Sprintf("[Download] start download chunk from %s", config.Cfg.Storage.Name))

	fileName := getFilePath(req.GetFileName(), req.GetChunkIndex())
	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open chunk file: %s", fileName)
	}
	defer file.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to read file part: %v", err)
		}

		// send data to client
		if err := stream.Send(&desc.DownloadChunkResponse{Data: buf[:n]}); err != nil {
			return status.Errorf(codes.Internal, "failed to send file part: %v", err)
		}
	}

	slog.Info(fmt.Sprintf("[Download] success download chunk file %s from %s", fileName, config.Cfg.Storage.Name))

	return nil
}

func (i *Implementation) DeleteChunk(_ context.Context, req *desc.DeleteChunkRequest) (*desc.DeleteChunkResponse, error) {
	filePath := getFilePath(req.GetFileName(), req.ChunkIndex)
	slog.Info(fmt.Sprintf("[Delete] start delete chunk file %s from %s", filePath, config.Cfg.Storage.Name))

	err := os.Remove(filePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "os.Remove(%s) err: %v", filePath, err)
	}

	slog.Info(fmt.Sprintf("[Delete] success delete chunk file %s from %s", filePath, config.Cfg.Storage.Name))

	return &desc.DeleteChunkResponse{}, nil
}

// ServerStats - server stats like disc space
func (i *Implementation) ServerStats(context.Context, *desc.ServerStatsRequest) (*desc.ServerStatsResponse, error) {
	info := du.NewDiskUsage(_filesPath)

	return &desc.ServerStatsResponse{
		DiscTotal: info.Size(),
		DiscUsed:  info.Used(),
		DiscAvail: info.Available(),
	}, nil
}

func getFilePath(fileName string, chunkIndex int32) string {
	return filepath.Join(_filesPath, fmt.Sprintf("%s.%d", fileName, chunkIndex))
}
