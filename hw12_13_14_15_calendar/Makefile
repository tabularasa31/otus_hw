#include .env.example
#export

BIN_CALENDAR="./bin/calendar"
BIN_SCHEDULER="./bin/scheduler"
BIN_SENDER="./bin/sender"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# ==============================================================================
# Main
swag: ### swag init
	swag init -g internal/controller/http/v1/router.go

generate:
	 go generate ./...

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run_calendar:
	$(BIN_CALENDAR) -config ./config/config_local.yml

run_scheduler:
	$(BIN_SCHEDULER) -config ./config/scheduler_config.yml

run_sender:
	$(BIN_SENDER) -config ./config/sender_config.yml

run: build swag run_calendar run_scheduler run_sender

test:
	go test -cover -race ./internal/...

integration-test:
	go clean -testcache && go test -v -tags integration ./tests/...

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
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint: install-lint-deps ### check by golangci linter
	echo "Starting linters"
	golangci-lint run

# ==============================================================================
# Go migrate postgresql

migrate:
	goose -dir migrations up

# ==============================================================================
# Docker compose commands
up:
	echo "Starting docker environment"
	docker-compose -f docker-compose.yaml --env-file .env.example --project-name calendar_service up --build

down:
	docker-compose -f docker-compose.yaml --project-name calendar_service down


# ==============================================================================
# Docker support
DOCKER_IMG="calendar:develop"

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/calendar/Dockerfile .
.PHONY: build-img


run-img: build-img
	docker run $(DOCKER_IMG)


# ==============================================================================
.PHONY: build run build-img run-img test lint migrate integration-test
