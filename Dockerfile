FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV VERSION=production

RUN go build -a -ldflags "-extldflags '-static' -X main.version=$VERSION" -o /app/bin/server /app/cmd/server

FROM scratch

COPY --from=builder /app/bin/server /usr/local/bin/server