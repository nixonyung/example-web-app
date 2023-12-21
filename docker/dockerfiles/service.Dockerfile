# (ref.) [Faster Multi-Platform Builds: Dockerfile Cross-Compilation Guide](https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/)

#
# builder
#
FROM golang:1.21 AS builder

# package: path relative to $GOPATH/src
ARG package
ARG cmd=./cmd/main

# (ref.) [Optimise docker build for go](https://medium.com/@marcin.niemira/optimise-docker-build-for-go-c03d6eb8b4b)
RUN go env -w GOCACHE=/gocache && go env -w GOMODCACHE=/gomodcache

WORKDIR $GOPATH/src/$package
RUN --mount=source=src,target=$GOPATH/src \
    --mount=type=cache,target=/gocache \
    --mount=type=cache,target=/gomodcache \
    go build -o /build/main $cmd

#
# runner
#
FROM golang:1.21 AS runner
RUN apt-get update && apt-get install -y \
    curl && rm -rf /var/lib/apt/lists/*
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.8.9
COPY --from=builder /build/main /main

RUN groupadd --system appgroup && useradd --system appuser --groups appgroup
USER appuser
# (ref.) [Docker Shell vs. Exec Form](https://emmer.dev/blog/docker-shell-vs.-exec-form/)
# (ref.) [Is it possible to run prefork mode inside docker?](https://github.com/gofiber/fiber/issues/1021#issuecomment-730537971)
CMD ["sh", "-c", "/main"]
