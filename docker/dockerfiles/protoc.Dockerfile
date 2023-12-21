# (ref.) [gRPC - Go - Quick start](https://grpc.io/docs/languages/go/quickstart/)

FROM golang:1.21 AS builder
# (ref.) [Best practices for Dockerfile instructions - RUN - apt-get](https://docs.docker.com/develop/develop-images/instructions/#apt-get)
RUN apt-get update && apt-get install --no-install-recommends -y \
    curl \
    unzip && rm -rf /var/lib/apt/lists/*
RUN go env -w GOCACHE=/gocache && go env -w GOMODCACHE=/gomodcache
ADD https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-linux-x86_64.zip protoc.zip
RUN unzip protoc.zip -d protoc && mv protoc/bin/protoc /usr/local/bin && mv protoc/include /usr/local
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

ENTRYPOINT ["protoc"]
