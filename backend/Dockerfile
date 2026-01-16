FROM golang:1.25.5 AS builder

WORKDIR /workspace

RUN --mount=type=bind,source=.,target=/workspace \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -a -tags netgo \
    -ldflags "-extldflags '-static' -s -w \
    -X main.version=$(git describe --tag --abbrev=0 2>/dev/null || echo 'dev') \
    -X main.revision=$(git rev-list -1 HEAD 2>/dev/null || echo 'unknown') \
    -X main.build=$(git describe --tags 2>/dev/null || echo 'dev')" \
    -o /bin/mawinter ./cmd/mawinter

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder --chmod=555 /bin/mawinter /mawinter

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/mawinter"]

CMD ["serve", "--port", "8080"]
