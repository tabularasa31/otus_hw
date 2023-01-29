# Step 1: Modules caching
FROM golang:1.17.1-alpine3.14 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.17.1-alpine3.14 as builder
COPY --from=modules /go/pkg /go/pkg
COPY .. /app
WORKDIR /app

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /bin/app ./cmd/calendar/*

# Step 2: Final
FROM scratch

LABEL service="Calendar"
LABEL microservice="Calendar main"
LABEL mainteiner="tabularasa31@gmail.com"

COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app

CMD ["/app"]