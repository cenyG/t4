// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0
// source: protos/chunk_storage.proto

package chunkstorage

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ChunkStorage_UploadChunk_FullMethodName   = "/chunkstorage.ChunkStorage/UploadChunk"
	ChunkStorage_DownloadChunk_FullMethodName = "/chunkstorage.ChunkStorage/DownloadChunk"
)

// ChunkStorageClient is the client API for ChunkStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChunkStorageClient interface {
	// Send file part to store server
	UploadChunk(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadChunkRequest, UploadChunkResponse], error)
	// Load file part from store server
	DownloadChunk(ctx context.Context, in *DownloadChunkRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadChunkResponse], error)
}

type chunkStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewChunkStorageClient(cc grpc.ClientConnInterface) ChunkStorageClient {
	return &chunkStorageClient{cc}
}

func (c *chunkStorageClient) UploadChunk(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadChunkRequest, UploadChunkResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ChunkStorage_ServiceDesc.Streams[0], ChunkStorage_UploadChunk_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadChunkRequest, UploadChunkResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChunkStorage_UploadChunkClient = grpc.ClientStreamingClient[UploadChunkRequest, UploadChunkResponse]

func (c *chunkStorageClient) DownloadChunk(ctx context.Context, in *DownloadChunkRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadChunkResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ChunkStorage_ServiceDesc.Streams[1], ChunkStorage_DownloadChunk_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadChunkRequest, DownloadChunkResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChunkStorage_DownloadChunkClient = grpc.ServerStreamingClient[DownloadChunkResponse]

// ChunkStorageServer is the server API for ChunkStorage service.
// All implementations must embed UnimplementedChunkStorageServer
// for forward compatibility.
type ChunkStorageServer interface {
	// Send file part to store server
	UploadChunk(grpc.ClientStreamingServer[UploadChunkRequest, UploadChunkResponse]) error
	// Load file part from store server
	DownloadChunk(*DownloadChunkRequest, grpc.ServerStreamingServer[DownloadChunkResponse]) error
	mustEmbedUnimplementedChunkStorageServer()
}

// UnimplementedChunkStorageServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChunkStorageServer struct{}

func (UnimplementedChunkStorageServer) UploadChunk(grpc.ClientStreamingServer[UploadChunkRequest, UploadChunkResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadChunk not implemented")
}
func (UnimplementedChunkStorageServer) DownloadChunk(*DownloadChunkRequest, grpc.ServerStreamingServer[DownloadChunkResponse]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadChunk not implemented")
}
func (UnimplementedChunkStorageServer) mustEmbedUnimplementedChunkStorageServer() {}
func (UnimplementedChunkStorageServer) testEmbeddedByValue()                      {}

// UnsafeChunkStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChunkStorageServer will
// result in compilation errors.
type UnsafeChunkStorageServer interface {
	mustEmbedUnimplementedChunkStorageServer()
}

func RegisterChunkStorageServer(s grpc.ServiceRegistrar, srv ChunkStorageServer) {
	// If the following call pancis, it indicates UnimplementedChunkStorageServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChunkStorage_ServiceDesc, srv)
}

func _ChunkStorage_UploadChunk_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChunkStorageServer).UploadChunk(&grpc.GenericServerStream[UploadChunkRequest, UploadChunkResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChunkStorage_UploadChunkServer = grpc.ClientStreamingServer[UploadChunkRequest, UploadChunkResponse]

func _ChunkStorage_DownloadChunk_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadChunkRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChunkStorageServer).DownloadChunk(m, &grpc.GenericServerStream[DownloadChunkRequest, DownloadChunkResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChunkStorage_DownloadChunkServer = grpc.ServerStreamingServer[DownloadChunkResponse]

// ChunkStorage_ServiceDesc is the grpc.ServiceDesc for ChunkStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChunkStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chunkstorage.ChunkStorage",
	HandlerType: (*ChunkStorageServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadChunk",
			Handler:       _ChunkStorage_UploadChunk_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadChunk",
			Handler:       _ChunkStorage_DownloadChunk_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protos/chunk_storage.proto",
}
