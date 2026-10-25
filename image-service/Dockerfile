# ------------------------------------------------------------------------------
# STAGE 1: Compilación de binario estático en Go
# ------------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ ./pkg/
COPY image-service/ ./image-service/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/image-service ./image-service

# ------------------------------------------------------------------------------
# STAGE 2: Imagen mínima de ejecución
# ------------------------------------------------------------------------------
FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -g '' appuser && \
    mkdir -p /app/uploads && \
    chown -R appuser:appuser /app/uploads

USER appuser

COPY --from=builder /bin/image-service /app/image-service

EXPOSE 8002

ENTRYPOINT ["/app/image-service"]
