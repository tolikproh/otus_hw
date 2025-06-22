include .env
export

BIN_API := "./bin/calendar"
BIN_SCHEDURLER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"
PROTO_DIR="./api"
GENERATE_DIR="./internal/pb"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build run build-img run-img version test lint generate

build:
	go build -v -o $(BIN_API) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDURLER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) --config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

 test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose -dir ./migrations up

migdown:
	goose -dir ./migrations down

generate:
	rm -rf $(GENERATE_DIR)
	mkdir -p $(GENERATE_DIR)
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GENERATE_DIR) \
		--go-grpc_out=$(GENERATE_DIR) \
		--grpc-gateway_out=$(GENERATE_DIR) \
		$(PROTO_DIR)/*.proto
	go generate -v ./...