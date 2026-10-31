# syntax=docker/dockerfile:1
ARG GO_VERSION=1.27.3
FROM golang:${GO_VERSION}-alpine AS builder
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=${GOPROXY}
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY app ./app
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/admin ./app/admin/api

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g 10001 app && adduser -D -u 10001 -G app app \
    && mkdir -p /app/var/uploads /app/var/avatars /app/var/downloads /app/var/generated \
    && chown -R app:app /app
WORKDIR /app
COPY --from=builder /out/admin /app/admin
COPY deploy/admin-api.yaml /app/etc/admin-api.yaml
COPY LICENSE app/admin/api/internal/generator/LICENSE.ruoyi /app/licenses/
USER app
EXPOSE 8888
ENTRYPOINT ["/app/admin"]
CMD ["-f", "/app/etc/admin-api.yaml"]
