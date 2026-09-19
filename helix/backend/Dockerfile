## ---- build stage ----------------------------------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download 2>/dev/null || true

COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /helix-server ./cmd/server

## ---- runtime stage ----------------------------------------------------
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S helix && adduser -S helix -G helix

WORKDIR /app
COPY --from=builder /helix-server ./helix-server

USER helix
EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=5 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./helix-server"]
