# --- Estágio de Compilação (Build) ---
FROM golang:1.26-alpine AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências para aproveitar o cache de camadas do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código fonte
# O Docker ignora arquivos listados no .dockerignore (como node_modules ou binários locais)
COPY . .

# Compila o binário
# GOOS=linux garante que o binário funcione no container Alpine
# CGO_ENABLED=0 cria um binário estático (sem dependências de C externas)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go
# --- Estágio Final (Runtime) ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Cria o usuário, mas vamos trabalhar no /app
RUN adduser -D appuser
WORKDIR /app

# Altera o dono da pasta para o appuser conseguir trabalhar nela se precisar
RUN chown appuser:appuser /app

# Copia o executável do builder para o /app atual
COPY --from=builder /app/server .

# SE sua API precisar ler a pasta de migrations por dentro do código Go, descomente a linha abaixo:
# COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8080

CMD ["./server"]