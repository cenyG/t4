g:
	protoc --go_out=./internal/restserver --go-grpc_out=./internal/restserver ./protos/chunk_storage.proto
	protoc --go_out=./internal/storageserver --go-grpc_out=./internal/storageserver ./protos/chunk_storage.proto

build-docker-images:
	docker build -t t4/storage --target storage .
	docker build -t t4/rest --target rest .

dc-main:
	docker-compose -f docker-compose-main.yml up --build

dc-store:
	docker-compose -f docker-compose-storages.yml up --build

dcd:
	docker-compose -f docker-compose-main.yml down
	docker-compose -f docker-compose-storages.yml down

migrate-up:
	goose -dir migrations postgres "user=user password=password dbname=t4 sslmode=disable" up