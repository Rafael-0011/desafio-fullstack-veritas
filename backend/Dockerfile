# Stage 1: Compilação do binário Go
FROM golang:alpine AS builder

WORKDIR /app

# Instala certificados para dependências HTTPS
RUN apk --no-cache add ca-certificates

# Cache e download de dependências
COPY go.mod go.sum* ./
RUN go mod download

# Copia todo o código fonte do backend
COPY . .

# Compilação estática com otimização de tamanho
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Stage 2: Imagem de execução leve
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/server /app/server

# Exposição da porta 8080 conforme o edital
EXPOSE 8080

ENV PORT=8080

CMD ["/app/server"]
