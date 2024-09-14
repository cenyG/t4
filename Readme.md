
# T4_test_case

**T4_test_case** is a distributed file storage system that splits large files into chunks and stores them across multiple storage servers. 

## Features
- Split large files into chunks and store them on different storage servers.
- Support for file Upload and Download.
- Files are not stored on the REST server; They are streamed during both upload and download for efficiency
- GRPC communication between REST server and Storage servers.
- Dockerized environment for both the REST server and Storage servers.
- Database migrations managed using Goose.
- Service discovery for Storage servers using Consul

## Project Structure
- **REST Server**: Handles file uploads and downloads via HTTP.
- **Storage Server**: Stores file chunks and communicates with the REST server via GRPC.
- **Protobuf**: Used for defining GRPC services.
- **PostgreSQL**: Stores metadata about the files.
- **Goose**: Used for database migrations.
- **Consul**: Used for service discovery.
- **Static site**: Separated static website with front-end part for testing.



## Installation and Setup

### Step 1: Generate GRPC Code

Run the following command to generate the Go code from your `.proto` files:

```bash
make g
```

### Step 2: Start the Docker Containers

1. You firstly need to start docker-compose with Consul, Static-Site, Rest and Postgres:

    ```bash
    make dc-main
    ```

2. Then start storage servers:

    ```bash
    make dc-store
    ```

3. Optional: You can add more Storage servers by:

    ```bash
    make dc-store-more
    ```

### Step 3: Test with Web

Once everything is set up and running, you can go to the following URLs to test the application:
- [http://127.0.0.1:8081/](http://127.0.0.1:8081/)
- [http://localhost:8081/](http://localhost:8081/)