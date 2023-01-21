include .env.example
export
BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)


# ==============================================================================
# Main

swag: ### swag init
	swag init --parseDependency --parseInternal -g internal/controller/http/v1/router.go

generate:
	 go generate ./...

run:
	go mod tidy && go mod download && \
    DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 \
	$(BIN) -config ../../config/config.yml
.PHONY: run

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
.PHONY: build

test:
	go test -cover -race ./internal/...

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Tools commands

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps ### check by golangci linter
	echo "Starting linters"
	golangci-lint run

# ==============================================================================
# Go migrate postgresql

migrate:
	goose -dir migrations up
.PHONY: migrate


# ==============================================================================
# Docker compose commands


# ==============================================================================
# Docker support
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .
.PHONY: build-img

run-img: build-img
	docker run $(DOCKER_IMG)
.PHONY: run-img
