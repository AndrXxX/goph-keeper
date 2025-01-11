VERSION=1.0.0
DATE=11.01.2025

DATA_DIR=./data/
SERVER_PORT=8081

KEYS_PATH=${DATA_DIR}server/
PUBLIC_KEY_PATH=${KEYS_PATH}public_crypto.key
PRIVATE_KEY_PATH=${KEYS_PATH}private_crypto.key

CLIENT_PATH=./cmd/client/main.go
CLIENT_BINARY=clientapp

GENERATOR_PATH=./cmd/keygen/main.go

run: generate-keys run-server build-client run-client

stop: clean

generate-keys:
	if [ -d ${KEYS_PATH} ]; then echo ""; else mkdir -p ${KEYS_PATH}; fi
	go run -ldflags '-X main.privatePath=${PRIVATE_KEY_PATH} -X main.publicPath=${PUBLIC_KEY_PATH}' ${GENERATOR_PATH}

run-server:
	GOPH_KEEPER_SERVER_PORT=${SERVER_PORT} PRIVATE_CRYPTO_KEY_PATH=${PRIVATE_KEY_PATH} \
	PUBLIC_CRYPTO_KEY_PATH=${PUBLIC_KEY_PATH} BUILD_VERSION=${VERSION} BUILD_DATE=${DATE} \
	docker compose up -d

stop-server:
	docker compose down

build-client:
	go build -o ${DATA_DIR}${CLIENT_BINARY} -ldflags '-X main.buildVersion=${VERSION} -X main.buildDate=${DATE} -X main.serverHost=localhost:${SERVER_PORT}' ${CLIENT_PATH}

run-client: build-client
	${DATA_DIR}${CLIENT_BINARY}

test-coverage-total:
	go test ./... -coverprofile cover.out > /dev/null && go tool cover -func cover.out | grep total && rm cover.out

test-coverage-detail:
	go test ./... -coverprofile cover.out > /dev/null && go tool cover -func cover.out && rm cover.out

clean: stop-server
	go clean
	rm ${DATA_DIR}${CLIENT_BINARY}
