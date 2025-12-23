# ==============================================================================
# Dockerfile otimizado para produção - BookHub API
# Utiliza multi-stage build para criar uma imagem final mínima e segura
# ==============================================================================

# ==============================================================================
# STAGE 1: Builder
# Esta etapa compila o código Go e gera o binário da aplicação
# ==============================================================================
FROM golang:1.25-alpine AS builder

# Instala certificados CA necessários para HTTPS e git para dependências
# O tzdata é necessário para suporte a timezones
RUN apk add --no-cache ca-certificates git tzdata

# Define o diretório de trabalho dentro do container
WORKDIR /build

# Copia os arquivos de dependências primeiro para aproveitar o cache do Docker
# Se go.mod e go.sum não mudarem, as dependências não serão baixadas novamente
COPY go.mod go.sum ./

# Baixa as dependências do Go
# O cache do Docker manterá essas dependências entre builds
RUN go mod download

# Copia o restante do código fonte
# Isso é feito após o download das dependências para otimizar o cache
COPY . .

# Compila o binário com otimizações para produção:
# - CGO_DISABLED=0: Desabilita o CGO para criar um binário estático
# - GOOS=linux: Compila para Linux (compatível com o container Alpine)
# - -ldflags="-s -w": Remove símbolos de debug e tabela de símbolos (reduz tamanho)
# - -trimpath: Remove caminhos absolutos do binário (melhora reprodutibilidade)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /build/bookhub \
    ./cmd/api/main.go

# ==============================================================================
# STAGE 2: Final
# Esta etapa cria a imagem final mínima apenas com o necessário para execução
# ==============================================================================
FROM alpine:3.19

# Instala apenas o essencial para a aplicação rodar:
# - ca-certificates: Certificados para conexões HTTPS
# - tzdata: Dados de timezone para manipulação correta de datas
RUN apk add --no-cache ca-certificates tzdata

# Cria um usuário não-root para executar a aplicação
# Isso segue as melhores práticas de segurança para containers
RUN addgroup -g 1000 -S bookhub && \
    adduser -u 1000 -S bookhub -G bookhub

# Define o diretório de trabalho
WORKDIR /app

# Copia o binário compilado da etapa de build
COPY --from=builder /build/bookhub .

# Copia os arquivos necessários para a aplicação:
# - Migrations: Scripts SQL para criar/atualizar o banco de dados
# - OpenAPI: Especificação da API (necessário para servir o swagger)
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/api/openapi ./api/openapi

# Cria o diretório para Swagger UI
# Será populado pelo comando make swagger-ui ou montado via volume
RUN mkdir -p ./api/swagger-ui

# Define permissões corretas para os arquivos
# O usuário bookhub precisa apenas de permissão de leitura e execução
RUN chown -R bookhub:bookhub /app

# Muda para o usuário não-root
# A partir daqui, todos os comandos serão executados como 'bookhub'
USER bookhub

# Expõe a porta que a aplicação utiliza
# Esta é apenas uma documentação, o mapeamento real é feito com -p no docker run
EXPOSE 8080

# Define variáveis de ambiente padrão
# Estas podem ser sobrescritas no docker run ou docker-compose
ENV SERVER_PORT=8080 \
    GIN_MODE=release

# Comando de health check
# Verifica se a aplicação está respondendo corretamente
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando para executar a aplicação
# Usa a forma exec para que sinais sejam passados corretamente para o processo
ENTRYPOINT ["./bookhub"]
