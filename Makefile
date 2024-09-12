g:
	protoc --go_out=./internal/restserver --go-grpc_out=./internal/restserver ./protos/chunk_storage.proto
	protoc --go_out=./internal/storageserver --go-grpc_out=./internal/storageserver ./protos/chunk_storage.proto