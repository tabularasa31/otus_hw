FROM gomicro/goose

ADD ./migrations/*.sql .

ENV GOOSE_DRIVER="postgres"
ENV GOOSE_DBSTRING="postgres://postgres:passpass@postgres:5432/postgres?sslmode=disable"

CMD goose up
