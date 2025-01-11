SERVER_PORT=8081
CLIENT_PATH=./cmd/client/main.go
CLIENT_BINARY=client
VERSION=1.0.0
DATE=11.01.2025

run: run-server build-client run-client

run-server:
	GOPH_KEEPER_SERVER_PORT=${SERVER_PORT} docker compose up -d

stop-server:
	docker compose down

build-client:
	go build -o ./data/${CLIENT_BINARY} -ldflags '-X main.buildVersion=${VERSION} -X main.buildDate=${DATE} -X main.serverHost=localhost:${SERVER_PORT}' ${CLIENT_PATH}

run-client: build-client
	./data/${CLIENT_BINARY}

clean: stop-server
	go clean
	rm ./data/${CLIENT_BINARY}
