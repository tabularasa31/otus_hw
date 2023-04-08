# Step 1: Modules caching
FROM golang:1.17.1-alpine3.14 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.17.1-alpine3.14 as builder
COPY --from=modules /go/pkg /go/pkg
COPY .. /sender
WORKDIR /sender

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /bin/sender ./cmd/calendar_sender/*

# Step 2: Final
FROM scratch

LABEL service="Calendar"
LABEL microservice="Calendar Sender"
LABEL mainteiner="tabularasa31@gmail.com"

COPY --from=builder /sender/config /config
COPY --from=builder /bin/sender /sender

CMD ["/sender"]