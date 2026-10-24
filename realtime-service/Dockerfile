# ------------------------------------------------------------------------------
# STAGE 1: Compilación de binario estático en Go
# ------------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ ./pkg/
COPY realtime-service/ ./realtime-service/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/realtime-service ./realtime-service

# ------------------------------------------------------------------------------
# STAGE 2: Imagen mínima de ejecución
# ------------------------------------------------------------------------------
FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /bin/realtime-service /app/realtime-service

EXPOSE 8001

ENTRYPOINT ["/app/realtime-service"]
