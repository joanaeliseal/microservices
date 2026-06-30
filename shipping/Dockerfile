FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copiar proto (copiado pelo docker-compose context)
COPY microservices-proto/ /build/microservices-proto/

# Copiar código do shipping
COPY microservices/shipping/ /build/microservices/shipping/

WORKDIR /build/microservices/shipping

# Download dependências
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/shipping ./cmd/main.go

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/shipping .

EXPOSE 3002

CMD ["./shipping"]
