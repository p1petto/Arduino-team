FROM golang:1.21 AS builder

WORKDIR /build

COPY . .

RUN go mod download && \
    go install github.com/mattn/go-sqlite3 && \
    CGO_ENABLED=1 GOOS=linux go build -o server -v ./cmd/


FROM ubuntu:latest
EXPOSE 1090

WORKDIR /build
COPY --from=builder /build/server /build/server
ENTRYPOINT ["/build/server"]