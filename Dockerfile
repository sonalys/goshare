FROM golang:1.23 AS builder

RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV VERSION=production

RUN task build

FROM scratch

COPY --from=builder /app/bin/server /usr/local/bin/server
COPY --from=builder /app/bin/migration /usr/local/bin/migration