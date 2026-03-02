# -------------------------
# Stage 1: Build
# -------------------------
FROM golang:1.26-alpine AS builder

# Instalar dependencias del sistema necesarias para compilar
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copiar shared primero (referenciado por replace directive en go.mod)
# El contexto es la raíz del monorepo, así que el path es directo
COPY shared/ /shared/

# Copiar go.mod y go.sum del servicio para aprovechar cache de Docker
COPY services/auth/go.mod services/auth/go.sum ./
RUN go mod download

# Copiar el resto del código fuente del servicio
COPY services/auth/ .

# Compilar — binario estático para Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /bin/auth-service ./cmd/main.go

# -------------------------
# Stage 2: Runtime
# -------------------------
FROM alpine:3.21

# Certificados para HTTPS (OAuth2, SMTP) y zona horaria
RUN apk add --no-cache ca-certificates tzdata

# Usuario no-root por seguridad
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

# Solo copiar el binario compilado
COPY --from=builder /bin/auth-service .

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./auth-service"]