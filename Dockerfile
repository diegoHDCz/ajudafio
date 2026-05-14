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

# Adiciona certificados CA para permitir chamadas HTTPS externas
RUN apk --no-cache add ca-certificates

# Cria um usuário não-root por segurança (boa prática de produção)
RUN adduser -D appuser
USER appuser

WORKDIR /home/appuser/

# Copia apenas o executável final do estágio de build
COPY --from=builder /app/server .

# Porta que a aplicação escuta (ajuste conforme seu código)
EXPOSE 8080

# Executa o binário
CMD ["./server"]