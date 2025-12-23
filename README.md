# BookHub API

API REST em Go para gerenciamento de empréstimos de livros, desenvolvida seguindo os princípios de Clean Architecture.

## Índice

- [Funcionalidades](#funcionalidades)
- [Arquitetura](#arquitetura)
- [Tecnologias](#tecnologias)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Pré-requisitos](#pré-requisitos)
- [Instalação e Execução](#instalação-e-execução)
- [Endpoints da API](#endpoints-da-api)
- [Autenticação](#autenticação)
- [Banco de Dados](#banco-de-dados)
- [Testes](#testes)
- [Docker](#docker)
- [Decisões Técnicas](#decisões-técnicas)

## Funcionalidades

### Usuários

- Listar todos os usuários (com paginação)
- Criar novo usuário
- Buscar usuário por ID
- Editar usuário
- Desabilitar usuário

### Livros

- Listar todos os livros (com paginação e filtro de disponibilidade)
- Criar novo livro
- Buscar livro por ID
- Status de disponibilidade automático (mostra "Indisponível - todas as cópias emprestadas" quando não há cópias disponíveis)

### Empréstimos

- Emprestar livro para usuário
- Devolver livro
- Listar empréstimos (com filtros por usuário e status)

### Autenticação

- Login com JWT
- Proteção de rotas autenticadas

## Arquitetura

O projeto segue os princípios de **Clean Architecture**, separando as responsabilidades em camadas:

```
┌─────────────────────────────────────────────────────────────┐
│                      Infraestrutura                         │
│  (HTTP Handlers, Repositórios, Database, JWT, Middleware)   │
├─────────────────────────────────────────────────────────────┤
│                       Casos de Uso                          │
│        (UserUseCase, BookUseCase, LoanUseCase)              │
├─────────────────────────────────────────────────────────────┤
│                         Domínio                             │
│    (Entidades: User, Book, Loan | Interfaces: Repository)   │
└─────────────────────────────────────────────────────────────┘
```

### Princípios Aplicados

- **Dependency Inversion**: Camadas internas não dependem de camadas externas
- **Single Responsibility**: Cada componente tem uma única responsabilidade
- **Interface Segregation**: Interfaces específicas para cada caso de uso
- **Open/Closed**: Código aberto para extensão, fechado para modificação

## Tecnologias

| Tecnologia           | Descrição                | Justificativa                                               |
| -------------------- | ------------------------ | ----------------------------------------------------------- |
| **Go 1.25**          | Linguagem de programação | Performance, simplicidade e forte suporte a concorrência    |
| **Gin**              | Framework HTTP           | Framework web mais popular e performático do ecossistema Go |
| **golang-jwt/jwt**   | Biblioteca JWT           | Biblioteca JWT mais utilizada pela comunidade Go            |
| **PostgreSQL**       | Banco de dados           | Robusto, confiável e com excelente suporte a JSON e UUID    |
| **MongoDB**          | Banco de dados NoSQL     | Flexibilidade de esquema e escalabilidade horizontal        |
| **database/sql**     | Driver de banco          | Biblioteca padrão do Go, sem dependências externas          |
| **lib/pq**           | Driver PostgreSQL        | Driver PostgreSQL mais estável para Go                      |
| **mongo-driver**     | Driver MongoDB           | Driver oficial da MongoDB para Go                           |
| **SQLC**             | Gerador de código SQL    | Type-safe, gera código Go a partir de queries SQL           |
| **oapi-codegen**     | Gerador OpenAPI          | Gera handlers e types a partir da especificação OpenAPI     |
| **mockgen**          | Gerador de mocks         | Gera mocks para interfaces, facilitando testes unitários    |
| **testcontainers**   | Testes de integração     | Containers Docker para testes de integração confiáveis      |
| **bcrypt**           | Hash de senhas           | Algoritmo padrão e seguro para hash de senhas               |
| **UUID**             | Identificadores          | IDs universalmente únicos para entidades                    |

## Estrutura do Projeto

```
go-bookhub/
├── api/
│   ├── generated/                 # Código gerado pelo oapi-codegen
│   │   └── openapi.gen.go         # Types e handlers gerados
│   ├── openapi/
│   │   └── openapi.yaml           # Especificação OpenAPI 3.0
│   └── swagger-ui/                # Arquivos do Swagger UI
├── cmd/
│   ├── api/
│   │   └── main.go                # Entry point da aplicação (PostgreSQL)
│   └── api-mongo/
│       └── main.go                # Entry point da aplicação (MongoDB)
├── internal/
│   ├── config/
│   │   └── config.go              # Configuração via variáveis de ambiente
│   ├── domain/
│   │   ├── entity/                # Entidades de domínio
│   │   │   ├── user.go            # Entidade User
│   │   │   ├── user_test.go       # Testes da entidade User
│   │   │   ├── book.go            # Entidade Book
│   │   │   ├── book_test.go       # Testes da entidade Book
│   │   │   ├── loan.go            # Entidade Loan
│   │   │   └── loan_test.go       # Testes da entidade Loan
│   │   └── repository/            # Interfaces dos repositórios
│   │       ├── user_repository.go
│   │       ├── book_repository.go
│   │       └── loan_repository.go
│   ├── infrastructure/
│   │   ├── auth/
│   │   │   ├── jwt.go             # Serviço JWT
│   │   │   └── jwt_test.go        # Testes do serviço JWT
│   │   ├── database/
│   │   │   ├── postgres.go        # Conexão PostgreSQL
│   │   │   ├── mongo.go           # Conexão MongoDB
│   │   │   └── sqlc/              # Configuração e código SQLC
│   │   │       ├── sqlc.yaml      # Configuração do SQLC
│   │   │       ├── queries/       # Queries SQL
│   │   │       ├── db.go          # Interface gerada
│   │   │       ├── models.go      # Models gerados
│   │   │       └── *.sql.go       # Código gerado
│   │   ├── http/
│   │   │   ├── router.go          # Configuração de rotas
│   │   │   ├── handler/
│   │   │   │   ├── handler.go     # Implementação dos handlers
│   │   │   │   ├── auth.go        # Handler de autenticação
│   │   │   │   ├── user.go        # Handler de usuários
│   │   │   │   ├── book.go        # Handler de livros
│   │   │   │   ├── loan.go        # Handler de empréstimos
│   │   │   │   ├── helpers.go     # Funções auxiliares
│   │   │   │   └── *_test.go      # Testes dos handlers
│   │   │   └── middleware/
│   │   │       └── auth.go        # Middleware de autenticação
│   │   └── repository/            # Implementação dos repositórios
│   │       ├── user_repository_postgres.go
│   │       ├── book_repository_postgres.go
│   │       ├── loan_repository_postgres.go
│   │       ├── user_repository_mongo.go
│   │       ├── book_repository_mongo.go
│   │       ├── loan_repository_mongo.go
│   │       ├── mongo_models.go    # Models para MongoDB
│   │       └── *_integration_test.go  # Testes de integração
│   ├── mocks/                     # Mocks gerados pelo mockgen
│   │   ├── mock_user_usecase.go
│   │   ├── mock_book_usecase.go
│   │   ├── mock_loan_usecase.go
│   │   └── mock_jwt_service.go
│   └── usecase/                   # Casos de uso
│       ├── user_usecase.go
│       ├── user_usecase_test.go
│       ├── book_usecase.go
│       ├── book_usecase_test.go
│       ├── loan_usecase.go
│       └── loan_usecase_test.go
├── migrations/                    # Migrações
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_books.up.sql
│   ├── 000002_create_books.down.sql
│   ├── 000003_create_loans.up.sql
│   ├── 000003_create_loans.down.sql
│   └── mongo/
│       └── init-db.js             # Script de inicialização MongoDB
├── .dockerignore
├── .env.example                   # Exemplo de variáveis de ambiente
├── .gitignore
├── docker-compose.yaml            # Docker Compose para desenvolvimento
├── Dockerfile                     # Dockerfile PostgreSQL
├── Dockerfile.mongo               # Dockerfile MongoDB
├── go.mod
├── go.sum
├── Makefile                       # Automação de tarefas
└── README.md
```

## Pré-requisitos

- Go 1.25 ou superior
- PostgreSQL 18 ou superior (para backend PostgreSQL)
- MongoDB 7 ou superior (para backend MongoDB)
- Docker e Docker Compose (opcional, mas recomendado para testes de integração)
- Make (opcional, mas recomendado)

## Instalação e Execução

### Usando Make (Recomendado)

#### Backend PostgreSQL

```bash
# Setup completo do projeto (instala ferramentas, dependências e gera código)
make setup

# Iniciar banco de dados PostgreSQL
docker compose up -d postgres

# Executar migrações
make migrate-up

# Executar aplicação
make run
```

#### Backend MongoDB

```bash
# Setup completo do projeto
make setup

# Iniciar banco de dados MongoDB
docker compose up -d mongodb

# Executar aplicação MongoDB
make run-mongo
```

### Manualmente

```bash
# Instalar ferramentas
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install go.uber.org/mock/mockgen@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Baixar dependências
go mod download

# Gerar código OpenAPI
mkdir -p api/generated
oapi-codegen -generate types,gin,spec -package generated -o api/generated/openapi.gen.go api/openapi/openapi.yaml

# Gerar código SQLC
cd internal/infrastructure/database/sqlc && sqlc generate && cd -

# Gerar mocks
make generate-mocks

# Criar banco de dados PostgreSQL
createdb bookhub

# Executar migrações
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/bookhub?sslmode=disable" up

# Executar aplicação PostgreSQL
go run cmd/api/main.go

# Ou executar aplicação MongoDB
MONGO_URI=mongodb://mongo:mongo@localhost:27017 MONGO_DATABASE=bookhub go run cmd/api-mongo/main.go
```

### Usando Docker Compose

```bash
# Executar todas as aplicações (PostgreSQL + MongoDB)
docker-compose up --build

# Ou em background
docker-compose up -d --build

# Executar apenas PostgreSQL
docker-compose up -d postgres api

# Executar apenas MongoDB
docker-compose up -d mongodb api-mongo
```

### Variáveis de Ambiente

Copie o arquivo `.env.example` para `.env` e ajuste conforme necessário:

```bash
cp .env.example .env
```

#### Variáveis Comuns

| Variável               | Descrição             | Padrão      |
| ---------------------- | --------------------- | ----------- |
| `SERVER_PORT`          | Porta do servidor     | `8080`      |
| `SERVER_READ_TIMEOUT`  | Timeout de leitura    | `15s`       |
| `SERVER_WRITE_TIMEOUT` | Timeout de escrita    | `15s`       |
| `JWT_SECRET_KEY`       | Chave secreta JWT     | -           |
| `JWT_TOKEN_DURATION`   | Duração do token      | `24h`       |
| `JWT_ISSUER`           | Emissor do token      | `bookhub`   |

#### PostgreSQL

| Variável      | Descrição             | Padrão      |
| ------------- | --------------------- | ----------- |
| `DB_HOST`     | Host do PostgreSQL    | `localhost` |
| `DB_PORT`     | Porta do PostgreSQL   | `5432`      |
| `DB_USER`     | Usuário do PostgreSQL | `postgres`  |
| `DB_PASSWORD` | Senha do PostgreSQL   | `postgres`  |
| `DB_NAME`     | Nome do banco         | `bookhub`   |
| `DB_SSL_MODE` | Modo SSL              | `disable`   |

#### MongoDB

| Variável         | Descrição         | Padrão                                |
| ---------------- | ----------------- | ------------------------------------- |
| `MONGO_URI`      | URI de conexão    | `mongodb://mongo:mongo@localhost:27017` |
| `MONGO_DATABASE` | Nome do banco     | `bookhub`                             |

## Endpoints da API

### Autenticação

| Método | Endpoint             | Descrição          | Autenticação |
| ------ | -------------------- | ------------------ | ------------ |
| POST   | `/api/v1/auth/login` | Autenticar usuário | Não          |

### Usuários

| Método | Endpoint                     | Descrição             | Autenticação |
| ------ | ---------------------------- | --------------------- | ------------ |
| GET    | `/api/v1/users`              | Listar usuários       | Sim          |
| POST   | `/api/v1/users`              | Criar usuário         | Sim          |
| GET    | `/api/v1/users/{id}`         | Buscar usuário por ID | Sim          |
| PUT    | `/api/v1/users/{id}`         | Atualizar usuário     | Sim          |
| PATCH  | `/api/v1/users/{id}/disable` | Desabilitar usuário   | Sim          |

### Livros

| Método | Endpoint             | Descrição           | Autenticação |
| ------ | -------------------- | ------------------- | ------------ |
| GET    | `/api/v1/books`      | Listar livros       | Sim          |
| POST   | `/api/v1/books`      | Criar livro         | Sim          |
| GET    | `/api/v1/books/{id}` | Buscar livro por ID | Sim          |

### Empréstimos

| Método | Endpoint                    | Descrição          | Autenticação |
| ------ | --------------------------- | ------------------ | ------------ |
| GET    | `/api/v1/loans`             | Listar empréstimos | Sim          |
| POST   | `/api/v1/loans/borrow`      | Emprestar livro    | Sim          |
| PATCH  | `/api/v1/loans/{id}/return` | Devolver livro     | Sim          |

### Swagger UI

Após iniciar a aplicação, acesse a documentação interativa:

- **Swagger UI**: <http://localhost:8080/docs>
- **OpenAPI Spec**: <http://localhost:8080/swagger.yaml>

## Autenticação

A API utiliza autenticação JWT (JSON Web Token). Para acessar endpoints protegidos:

1. Faça login no endpoint `/api/v1/auth/login`
2. Use o token retornado no header `Authorization: Bearer <token>`

### Usuário Admin Padrão

As migrações criam um usuário administrador padrão:

| Campo | Valor               |
| ----- | ------------------- |
| Email | `admin@bookhub.com` |
| Senha | `admin123`          |

### Exemplo de Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@bookhub.com", "password": "admin123"}'
```

Resposta:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2024-01-02T15:04:05Z",
  "user": {
    "id": "uuid",
    "name": "Admin",
    "email": "admin@bookhub.com",
    "active": true
  }
}
```

### Exemplo de Requisição Autenticada

```bash
curl -X GET http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Banco de Dados

### Diagrama ER

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     users       │       │     loans       │       │     books       │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │       │ id (PK)         │
│ name            │       │ user_id (FK)    │───────│ title           │
│ email (UNIQUE)  │───────│ book_id (FK)    │       │ author          │
│ password_hash   │       │ borrowed_at     │       │ isbn (UNIQUE)   │
│ active          │       │ due_date        │       │ published_year  │
│ created_at      │       │ returned_at     │       │ total_copies    │
│ updated_at      │       │ status          │       │ available_copies│
└─────────────────┘       └─────────────────┘       │ created_at      │
                                                    │ updated_at      │
                                                    └─────────────────┘
```

### Migrações

As migrações estão em `migrations/` e são executadas automaticamente pelo comando `make migrate-up`.

```bash
# Executar migrações
make migrate-up

# Reverter última migração
make migrate-down

# Criar nova migração
make migrate-create name=nome_da_migracao
```

## Testes

O projeto possui testes em todas as camadas, incluindo testes unitários e de integração com testcontainers.

### Testes Unitários

```bash
# Executar todos os testes unitários
make test

# Executar testes com cobertura
make test-coverage

# Executar testes de um pacote específico
go test ./internal/domain/entity/... -v
go test ./internal/usecase/... -v
go test ./internal/infrastructure/auth/... -v
go test ./internal/infrastructure/http/handler/... -v
```

### Testes de Integração

Os testes de integração utilizam **testcontainers** para criar containers Docker temporários dos bancos de dados:

```bash
# Executar testes de integração (requer Docker)
make test-integration

# Executar todos os testes (unitários + integração)
make test-all

# Executar todos os testes com cobertura
make test-coverage-all
```

### Gerando Mocks

Os mocks são gerados automaticamente com **mockgen**:

```bash
# Gerar todos os mocks
make generate-mocks
```

### Cobertura de Testes

- **Domínio**: Entidades User, Book, Loan
- **Casos de Uso**: UserUseCase, BookUseCase, LoanUseCase
- **Infraestrutura**: JWTService, Handlers HTTP
- **Repositórios**: Testes de integração para PostgreSQL e MongoDB

## Docker

### Dockerfiles

O projeto possui dois Dockerfiles, um para cada backend:

- `Dockerfile` - Backend PostgreSQL
- `Dockerfile.mongo` - Backend MongoDB

Ambos utilizam multi-stage build para criar imagens otimizadas:

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
# Compila o binário

# Stage 2: Runtime
FROM alpine:3.19
# Imagem final mínima (~15MB)
```

Características:

- Imagem final baseada em Alpine (~15MB)
- Binário estático sem dependências externas
- Usuário não-root para segurança
- Health check configurado
- Multi-stage build para menor tamanho

### Comandos Docker

```bash
# Build da imagem PostgreSQL
make docker-build

# Build da imagem MongoDB
make docker-build-mongo

# Build de ambas as imagens
make docker-build-all

# Executar container PostgreSQL
make docker-run

# Executar container MongoDB
make docker-run-mongo

# Docker Compose (API + PostgreSQL + MongoDB)
make docker-compose

# Parar containers
make docker-compose-down
```

## Decisões Técnicas

### 1. Clean Architecture

Escolhida para garantir separação de responsabilidades, facilitar testes e permitir evolução independente das camadas.

### 2. Gin Framework

Framework HTTP mais popular e performático do ecossistema Go, com excelente documentação e comunidade ativa.

### 3. golang-jwt/jwt

Biblioteca JWT mais utilizada pela comunidade Go, com suporte completo às especificações JWT e manutenção ativa.

### 4. SQLC

Gera código Go type-safe a partir de queries SQL, evitando erros em runtime e melhorando a experiência de desenvolvimento.

### 5. oapi-codegen

Gera handlers e types a partir da especificação OpenAPI, garantindo que a API sempre esteja em sincronia com a documentação.

### 6. Suporte a Múltiplos Bancos de Dados

O projeto suporta tanto PostgreSQL quanto MongoDB, demonstrando a flexibilidade da Clean Architecture:

- **PostgreSQL**: Utiliza `database/sql` com driver `lib/pq` e SQLC para queries type-safe
- **MongoDB**: Utiliza o driver oficial `mongo-driver` com models específicos

### 7. Testcontainers para Testes de Integração

Testes de integração confiáveis utilizando containers Docker temporários, garantindo ambiente isolado e reproduzível.

### 8. Mockgen para Testes Unitários

Geração automática de mocks a partir de interfaces, facilitando testes unitários isolados.

### 9. UUID para IDs

Identificadores universalmente únicos evitam colisões e permitem geração distribuída de IDs.

### 10. bcrypt para Senhas

Algoritmo padrão da indústria para hash de senhas, com salt automático e custo configurável.

### 11. Variáveis de Ambiente

Configuração via variáveis de ambiente seguindo os princípios do 12-Factor App.

### 12. Multi-stage Docker Build

Reduz o tamanho da imagem final e melhora a segurança ao não incluir ferramentas de build.

## Comandos Make Disponíveis

```bash
make help              # Mostra todos os comandos disponíveis

# Build
make build             # Compila o binário PostgreSQL
make build-mongo       # Compila o binário MongoDB
make build-all         # Compila ambos os binários

# Execução
make run               # Executa a aplicação PostgreSQL
make run-mongo         # Executa a aplicação MongoDB

# Testes
make test              # Executa testes unitários
make test-integration  # Executa testes de integração (requer Docker)
make test-all          # Executa todos os testes
make test-coverage     # Executa testes unitários com cobertura
make test-coverage-all # Executa todos os testes com cobertura

# Geração de Código
make generate          # Gera código (OpenAPI + SQLC)
make generate-api      # Gera apenas código OpenAPI
make generate-sqlc     # Gera apenas código SQLC
make generate-mocks    # Gera mocks com mockgen

# Migrações
make migrate-up        # Executa migrações
make migrate-down      # Reverte última migração
make migrate-create    # Cria nova migração (name=nome_da_migracao)

# Docker
make docker-build       # Build da imagem PostgreSQL
make docker-build-mongo # Build da imagem MongoDB
make docker-build-all   # Build de ambas as imagens
make docker-run         # Executa container PostgreSQL
make docker-run-mongo   # Executa container MongoDB
make docker-compose     # Executa com Docker Compose
make docker-compose-down # Para containers Docker Compose

# Utilitários
make lint              # Executa linter
make deps              # Baixa dependências
make tools             # Instala ferramentas necessárias
make swagger-ui        # Baixa arquivos do Swagger UI
make setup             # Setup completo do projeto
make clean             # Remove artefatos de build
```

## Licença

Este projeto é para fins educacionais.
