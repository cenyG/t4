package grpc

import (
	"context"
	"io"
	"strconv"

	"T4_test_case/config"
	desc "T4_test_case/pb/chunkstorage"
	"T4_test_case/pkg/helper"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChunkStorageClient interface {
	UploadChunkStream(ctx context.Context, fileName string, chunkIndex int32, chunkData io.Reader, chunkSize int64) error
	DownloadChunkStream(ctx context.Context, fileName string, chunkIndex int32, writer io.Writer) error
	ServerStats(ctx context.Context, req *desc.ServerStatsRequest) (*desc.ServerStatsResponse, error)
	DeleteChunk(ctx context.Context, req *desc.DeleteChunkRequest) (*desc.DeleteChunkResponse, error)
	ServiceID() string
}

// chunkStorageClient - grpc client
type chunkStorageClient struct {
	client    desc.ChunkStorageClient
	serviceId string
}

type bufReadResult struct {
	err error
	n   int
}

func (s *chunkStorageClient) ServiceID() string {
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
	defer stream.CloseAndRecv()

	buffer := make([]byte, 1024*1024)
	totalSent := int64(0)

	for totalSent < chunkSize {
		// close connection if user don't send anything for a long time
		n, readErr := helper.WithTimeout[int](ctx, config.Cfg.Rest.UploadBytesWaitTime, func() (int, error) {
			return chunkData.Read(buffer)
		})
		if readErr != nil {
			return errors.Wrap(readErr, "failed while read from client")
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

// ServerStats - server stats
func (s *chunkStorageClient) ServerStats(ctx context.Context, req *desc.ServerStatsRequest) (*desc.ServerStatsResponse, error) {
	return s.client.ServerStats(ctx, req)
}

// DeleteChunk - delete chunk from server
func (s *chunkStorageClient) DeleteChunk(ctx context.Context, req *desc.DeleteChunkRequest) (*desc.DeleteChunkResponse, error) {
	return s.client.DeleteChunk(ctx, req)
}
