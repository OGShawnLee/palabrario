# ------------------------------------------------------------------------------
# STAGE 1: Compilación de binario estático en Go
# ------------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar ca-certificates y git por si se necesitan
RUN apk add --no-cache ca-certificates git

# Copiar manifiestos de dependencias desde la raíz del proyecto
COPY go.mod go.sum ./

# Descargar módulos y cachear capas
RUN go mod download

# Copiar el código de backend y paquetes compartidos
COPY pkg/ ./pkg/
COPY core-service/ ./core-service/

# Compilar binario estático optimizado (-s -w elimina símbolos de depuración)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/core-service ./core-service

# ------------------------------------------------------------------------------
# STAGE 2: Imagen mínima de ejecución
# ------------------------------------------------------------------------------
FROM alpine:3.21

WORKDIR /app

# Certificados SSL para conexiones seguras
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario sin privilegios por seguridad
RUN adduser -D -g '' appuser
USER appuser

# Copiar binario compilado desde el builder
COPY --from=builder /bin/core-service /app/core-service

EXPOSE 8000

ENTRYPOINT ["/app/core-service"]
