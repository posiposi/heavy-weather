FROM golang:1.26.6-bookworm

WORKDIR /app

ENV GOFLAGS=-buildvcs=false

RUN go install -tags='no_clickhouse no_libsql no_mssql no_postgres no_sqlite3 no_vertica no_ydb' \
    github.com/pressly/goose/v3/cmd/goose@v3.27.3 \
    && go clean -modcache -cache

CMD ["sleep", "infinity"]
