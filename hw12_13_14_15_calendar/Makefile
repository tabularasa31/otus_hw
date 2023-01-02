include .env.example
export

BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

swag-v1: ### swag init
	swag init --parseDependency --parseInternal -g internal/controller/http/v1/router.go
.PHONY: swag-v1

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
.PHONY: build

run: swag-v1
	go mod tidy && go mod download && \
    DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 \
	$(BIN) -config ../../config/config.yml
.PHONY: run

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .
.PHONY: build-img

run-img: build-img
	docker run $(DOCKER_IMG)
.PHONY: run-img

version: build
	$(BIN) version
.PHONY: version

test:
	go test -race ./internal/...
.PHONY: test

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps ### check by golangci linter
	golangci-lint run
.PHONY: lint

migrate:
	goose -dir migrations up
.PHONY: migrate
