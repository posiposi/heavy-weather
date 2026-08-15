FROM golang:1.26.6-bookworm

WORKDIR /app

ENV GOFLAGS=-buildvcs=false

CMD ["sleep", "infinity"]
