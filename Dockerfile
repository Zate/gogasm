# ---------------
# Build Container
# ---------------
FROM golang:alpine AS build-env

# install build tools
RUN apk add --no-cache tzdata bash wget curl git build-base openssl-dev \
    && mkdir -p $$GOPATH/bin

COPY . /go/src/github.com/xhochn/gogasm
WORKDIR /go/src/github.com/xhochn/gogasm

RUN go install ./cmd/

# --------------------
# Production Container
# --------------------
FROM alpine:3.8

# hadolint ignore=DL3018
RUN apk add --update --no-cache ca-certificates openssl-dev

WORKDIR /app
COPY --from=build-env /go/bin/cmd /app/

ENTRYPOINT ["./cmd"]
