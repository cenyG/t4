package pb

import (
	"T4_test_case/internal/restserver/pb/chunkstorage"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"strconv"
)

type ChunkStorageClient interface {
	UploadChunkStream(ctx context.Context, fileName string, chunkIndex int32, chunkData io.Reader, chunkSize int64) error
}

// ChunkStorageClient .
type chunkStorageClient struct {
	client chunkstorage.ChunkStorageClient
}

// NewChunkStorageClient .
func NewChunkStorageClient(host string, port int) (ChunkStorageClient, error) {
	addr := host + ":" + strconv.Itoa(port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to storage server: %s", addr)
	}

	return &chunkStorageClient{
		client: chunkstorage.NewChunkStorageClient(conn),
	}, nil
}

// UploadChunkStream - sent chunk to storage server
func (s *chunkStorageClient) UploadChunkStream(ctx context.Context, fileName string, chunkIndex int32, chunkData io.Reader, chunkSize int64) error {
	stream, err := s.client.UploadChunk(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create upload stream")
	}

	buffer := make([]byte, 1024*1024)
	totalSent := int64(0)

	for totalSent < chunkSize {
		n, readErr := chunkData.Read(buffer)
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return errors.Wrap(readErr, "failed to read chunk")
		}

		sendErr := stream.Send(&chunkstorage.UploadChunkRequest{
			FileName:   fileName,
			ChunkIndex: chunkIndex,
			ChunkData:  buffer[:n],
		})
		if sendErr != nil {
			return errors.Wrap(sendErr, "failed to send chunk")
		}
		totalSent += int64(n)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return errors.Wrap(err, "stream.CloseAndRecv()")
	}

	return nil
}
