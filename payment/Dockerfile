FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copiar proto (copiado pelo docker-compose context)
COPY microservices-proto/ /build/microservices-proto/

# Copiar código do payment
COPY microservices/payment/ /build/microservices/payment/

WORKDIR /build/microservices/payment

# Download dependências
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/payment ./cmd/main.go

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/payment .

EXPOSE 3001

CMD ["./payment"]
