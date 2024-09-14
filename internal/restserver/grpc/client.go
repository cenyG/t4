package grpc

import (
	"context"
	"io"
	"strconv"

	desc "T4_test_case/pb/chunkstorage"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChunkStorageClient interface {
	UploadChunkStream(ctx context.Context, fileName string, chunkIndex int32, chunkData io.Reader, chunkSize int64) error
	DownloadChunkStream(ctx context.Context, fileName string, chunkIndex int32, writer io.Writer) error
	GetServiceID() string
}

// chunkStorageClient - grpc client
type chunkStorageClient struct {
	client    desc.ChunkStorageClient
	serviceId string
}

func (s *chunkStorageClient) GetServiceID() string {
	return s.serviceId
}

// NewChunkStorageClient .
func NewChunkStorageClient(host string, port int, serviceId string) (ChunkStorageClient, error) {
	addr := host + ":" + strconv.Itoa(port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to storage server: %s", addr)
	}

	return &chunkStorageClient{
		client:    desc.NewChunkStorageClient(conn),
		serviceId: serviceId,
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

		sendErr := stream.Send(&desc.UploadChunkRequest{
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

// DownloadChunkStream - load chunk from storage server
func (s *chunkStorageClient) DownloadChunkStream(ctx context.Context, fileName string, chunkIndex int32, writer io.Writer) error {
	req := &desc.DownloadChunkRequest{
		FileName:   fileName,
		ChunkIndex: chunkIndex,
	}

	stream, err := s.client.DownloadChunk(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to download file part")
	}

	for {
		resp, streamErr := stream.Recv()
		if streamErr == io.EOF {
			break
		}
		if streamErr != nil {
			return errors.Wrap(streamErr, "error receiving file part")
		}

		_, err = writer.Write(resp.Data)
		if err != nil {
			return errors.Wrap(err, "failed to write chunk")
		}
	}

	return nil
}
