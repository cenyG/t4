package pb

import (
	"T4_test_case/config"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	desc "T4_test_case/internal/storageserver/pb/chunkstorage"
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

// UploadChunk - upload file chunk
func (i *Implementation) UploadChunk(stream desc.ChunkStorage_UploadChunkServer) error {
	slog.Info(fmt.Sprintf("start upload chunk on %s", config.Cfg.Storage.Name))

	var file *os.File
	for {
		// Чтение данных из потока
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileName := chunk.GetFileName()
		chunkIndex := chunk.GetChunkIndex()

		// Create file
		if file == nil {
			err := os.MkdirAll(_filesPath, os.ModePerm)
			if err != nil {
				return err
			}
			filePath := filepath.Join(_filesPath, fmt.Sprintf("%s.%d", fileName, chunkIndex))
			file, err = os.Create(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
		}

		// Write chunk to file
		_, err = file.Write(chunk.GetChunkData())
		if err != nil {
			return err
		}
	}

	return stream.SendAndClose(&desc.UploadChunkResponse{})
}

// DownloadChunk - upload file chunk
func (i *Implementation) DownloadChunk(req *desc.DownloadChunkRequest, stream desc.ChunkStorage_DownloadChunkServer) error {
	fileName := fmt.Sprintf("%s/%s.%d", _filesPath, req.GetFileName(), req.GetChunkIndex())

	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open chunk file: %s", fileName)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read file part: %v", err)
		}

		// send data to client
		if err := stream.Send(&desc.DownloadChunkResponse{Data: buf[:n]}); err != nil {
			return fmt.Errorf("failed to send file part: %v", err)
		}
	}

	return nil
}
