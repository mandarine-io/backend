FROM golang:1.23.5-alpine3.21 AS deps

WORKDIR /app

COPY . .

ENV GO111MODULE=on

RUN go generate ./...
RUN go mod download

FROM deps AS build

WORKDIR /app

ENV CGO_ENABLED=0
ARG ARTIFACT_VERSION

RUN go build \
    -o ./build/server \
    -installsuffix "static" \
    -tags "" \
    -ldflags " \
        -X main.Version=${ARTIFACT_VERSION:-0.0.0} \
        -X main.GoVersion=$(go version | cut -d " " -f 3) \
        -X main.Compiler=$(go env CC) \
        -X main.Platform=$(go env GOOS)/$(go env GOARCH) \
    " \
    ./cmd/api/main.go

FROM alpine:3.21 AS runtime

WORKDIR /app

COPY --from=build /app/build/server /app/server

COPY config/config.default.yaml config/config.yaml
COPY locales locales
COPY migrations migrations
COPY templates templates

RUN apk update \
    && apk add --no-cache ca-certificates tzdata curl \
    && echo 'Etc/UTC' > /etc/timezone \
    && adduser --disabled-password --home /app --gecos '' gouser \
    && chown -R gouser /app

ENV TZ=Etc/UTC
ENV LANG=en_US.utf8
ENV LC_ALL=en_US.UTF-8

USER gouser

ENTRYPOINT [ "/app/server" ]