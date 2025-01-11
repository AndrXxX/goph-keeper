CLIENT_PATH=./cmd/client/main.go
CLIENT_BINARY=client
VERSION=1.0.0
DATE=11.01.2025

run: run-server build-client run-client

run-server:
	docker compose up -d

stop-server:
	docker compose down

build-client:
	go build -o ${CLIENT_BINARY} -ldflags '-X main.buildVersion=${VERSION} -X main.buildDate=${DATE}' ${CLIENT_PATH}

run-client: build-client
	./${CLIENT_BINARY}

clean: stop-server
	go clean
	rm ${BINARY_NAME}
