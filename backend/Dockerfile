FROM golang:1.22.12-alpine3.21 AS builder

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build \
    -mod=readonly \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w -buildid=" \
    -o /out/server \
    ./cmd/server

FROM scratch

COPY --from=builder --chown=65532:65532 /out/server /server

USER 65532:65532

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/server"]
