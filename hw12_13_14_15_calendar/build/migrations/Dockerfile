#FROM gomicro/goose
#
#LABEL service="Calendar"
#LABEL microservice="Migrations"
#LABEL mainteiner="tabularasa31@gmail.com"
#
#ADD ../migrations/*.sql /migrations/
#
#ENV DBSTRING="user=postgres password=passpass host=localhost dbname=calendarDB sslmode=disable"
#
#CMD goose postgres "$DBSTRING" up
#

#
#FROM golang:1.19
#
#COPY . .
#RUN go install github.com/pressly/goose/v3/cmd/goose@latest
#
#ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
#RUN chmod +x /wait
#
#CMD /wait && /go/bin/goose -dir=migrations up

# builder
#FROM golang:1.15-alpine AS builder
#
#RUN apk add git
#
#WORKDIR /app
#COPY . .
#
## database migrations
#RUN CGO_ENABLED=0 go get github.com/pressly/goose/cmd/goose
#RUN mv /go/bin/goose /app/goose
#
#ENTRYPOINT [ "/bin/ash", "-c", "./goose -dir migrations postgres version"]

# builder
FROM golang:1.15-alpine AS builder
WORKDIR /app
COPY . .

# database migrations
RUN CGO_ENABLED=0 go get -u github.com/pressly/goose/cmd/goose
RUN mv /go/bin/goose /app/goose

ENTRYPOINT [ "/bin/ash", "-c", "./goose -dir migrations postgres ${GOOSE_DBSTRING} up"]