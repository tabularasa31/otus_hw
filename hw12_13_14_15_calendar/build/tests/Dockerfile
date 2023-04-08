FROM golang:1.17-alpine

RUN mkdir -p /opt/integration_tests
WORKDIR /opt/integration_tests

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

COPY . .

CMD ["go", "test", "-v", "-tags", "integration", "./tests/."]

