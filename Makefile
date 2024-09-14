g:
	protoc --go_out=. --go-grpc_out=. ./protos/chunk_storage.proto

build-docker-images:
	docker build -t t4/storage --target storage .
	docker build -t t4/rest --target rest .

dc-main:
	docker-compose -f docker-compose-main.yml up --build

dc-store:
	docker-compose -f docker-compose-storages.yml up --build

dc-store-more:
	docker-compose -f docker-compose-storages-more.yml up --build

dc-clean:
	docker-compose -f docker-compose-main.yml down
	docker-compose -f docker-compose-main.yml down -v
	docker-compose -f docker-compose-storages.yml down
	docker-compose -f docker-compose-storages-more.yml down

migrate-up:
	goose -dir migrations postgres "user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable" up

migrate-up-dev:
	goose -dir migrations postgres "user=user password=password dbname=t4 sslmode=disable" up

lint:
	golangci-lint run