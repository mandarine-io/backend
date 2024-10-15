FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./cmd ./cmd
COPY ./docs ./docs
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN go build -o ./build/server ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/build/server /app/server

COPY ./config ./config
COPY ./locales ./locales
COPY ./migrations ./migrations
COPY ./templates ./templates

RUN cp config/config.example.yaml config/config.yaml

ENV MANDARINE_CONFIG__FILE=config/config.yaml

RUN adduser --disabled-password \
  --home /app \
  --gecos '' gouser && chown -R gouser /app
USER gouser

ENTRYPOINT /app/server