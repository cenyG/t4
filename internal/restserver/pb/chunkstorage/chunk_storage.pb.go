// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: protos/chunk_storage.proto

package chunkstorage

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Сообщение для отправки части файла
type UploadChunkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// File name
	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// Chunk index
	ChunkIndex int32 `protobuf:"varint,2,opt,name=chunk_index,json=chunkIndex,proto3" json:"chunk_index,omitempty"`
	// Chunk bytes
	ChunkData []byte `protobuf:"bytes,3,opt,name=chunk_data,json=chunkData,proto3" json:"chunk_data,omitempty"`
}

func (x *UploadChunkRequest) Reset() {
	*x = UploadChunkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chunk_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadChunkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadChunkRequest) ProtoMessage() {}

func (x *UploadChunkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chunk_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadChunkRequest.ProtoReflect.Descriptor instead.
func (*UploadChunkRequest) Descriptor() ([]byte, []int) {
	return file_protos_chunk_storage_proto_rawDescGZIP(), []int{0}
}

func (x *UploadChunkRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *UploadChunkRequest) GetChunkIndex() int32 {
	if x != nil {
		return x.ChunkIndex
	}
	return 0
}

func (x *UploadChunkRequest) GetChunkData() []byte {
	if x != nil {
		return x.ChunkData
	}
	return nil
}

// Ответ от сервера хранения
type UploadChunkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadChunkResponse) Reset() {
	*x = UploadChunkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chunk_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadChunkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadChunkResponse) ProtoMessage() {}

func (x *UploadChunkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chunk_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadChunkResponse.ProtoReflect.Descriptor instead.
func (*UploadChunkResponse) Descriptor() ([]byte, []int) {
	return file_protos_chunk_storage_proto_rawDescGZIP(), []int{1}
}

var File_protos_chunk_storage_proto protoreflect.FileDescriptor

var file_protos_chunk_storage_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x22, 0x71, 0x0a, 0x12, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x22, 0x15, 0x0a,
	0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0x64, 0x0a, 0x0c, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x12, 0x54, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x12, 0x20, 0x2e, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x42, 0x11, 0x5a, 0x0f, 0x70, 0x62,
	0x2f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_chunk_storage_proto_rawDescOnce sync.Once
	file_protos_chunk_storage_proto_rawDescData = file_protos_chunk_storage_proto_rawDesc
)

func file_protos_chunk_storage_proto_rawDescGZIP() []byte {
	file_protos_chunk_storage_proto_rawDescOnce.Do(func() {
		file_protos_chunk_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_chunk_storage_proto_rawDescData)
	})
	return file_protos_chunk_storage_proto_rawDescData
}

var file_protos_chunk_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protos_chunk_storage_proto_goTypes = []any{
	(*UploadChunkRequest)(nil),  // 0: chunkstorage.UploadChunkRequest
	(*UploadChunkResponse)(nil), // 1: chunkstorage.UploadChunkResponse
}
var file_protos_chunk_storage_proto_depIdxs = []int32{
	0, // 0: chunkstorage.ChunkStorage.UploadChunk:input_type -> chunkstorage.UploadChunkRequest
	1, // 1: chunkstorage.ChunkStorage.UploadChunk:output_type -> chunkstorage.UploadChunkResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_chunk_storage_proto_init() }
func file_protos_chunk_storage_proto_init() {
	if File_protos_chunk_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_chunk_storage_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*UploadChunkRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_chunk_storage_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*UploadChunkResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_chunk_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_chunk_storage_proto_goTypes,
		DependencyIndexes: file_protos_chunk_storage_proto_depIdxs,
		MessageInfos:      file_protos_chunk_storage_proto_msgTypes,
	}.Build()
	File_protos_chunk_storage_proto = out.File
	file_protos_chunk_storage_proto_rawDesc = nil
	file_protos_chunk_storage_proto_goTypes = nil
	file_protos_chunk_storage_proto_depIdxs = nil
}
