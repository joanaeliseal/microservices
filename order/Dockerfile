FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copiar proto (copiado pelo docker-compose context)
COPY microservices-proto/ /build/microservices-proto/

# Copiar código do order
COPY microservices/order/ /build/microservices/order/

WORKDIR /build/microservices/order

# Download dependências
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/order ./cmd/main.go

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/order .

EXPOSE 3000

CMD ["./order"]
