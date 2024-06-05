FROM golang:1.21 AS builder

WORKDIR /build

ADD go.mod .

RUN go mod download
COPY . .
RUN go install github.com/mattn/go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o server -v ./cmd/


FROM golang:1.21
EXPOSE 1090

WORKDIR /build

# ADD .env /build

COPY --from=builder /build/server /build/server


RUN chmod +x /build/server

CMD ["./server"]