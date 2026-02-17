# CondoTrack API

API REST em Go para o sistema CondoTrack, utilizando Clean Architecture.

## Tecnologias

- **Go 1.22** - Linguagem de programação
- **Gin** - Framework HTTP
- **SQLX** - Acesso ao banco de dados
- **MySQL 8.0** - Banco de dados
- **Docker** - Containerização
- **Asaas** - Gateway de pagamentos

## Estrutura do Projeto

```
condotrack-api/
├── cmd/server/          # Entry point da aplicação
├── internal/
│   ├── config/          # Configurações
│   ├── domain/
│   │   ├── entity/      # Entidades do domínio
│   │   └── repository/  # Interfaces de repositórios
│   ├── infrastructure/
│   │   ├── database/    # Conexão MySQL
│   │   ├── repository/  # Implementações MySQL
│   │   └── external/    # Integrações externas (Asaas)
│   ├── usecase/         # Casos de uso
│   └── delivery/
│       └── http/        # Handlers e rotas HTTP
├── pkg/                 # Pacotes compartilhados
├── migrations/          # Scripts SQL
├── Dockerfile
└── docker-compose.yml
```

## Requisitos

- Go 1.22+
- Docker e Docker Compose
- MySQL 8.0 (ou via Docker)

## Instalação

### Usando Docker (Recomendado)

```bash
# Clone o repositório
cd condotrack-api

# Configure as variáveis de ambiente
cp .env.example .env
# Edite o arquivo .env com suas configurações

# Inicie os containers
docker-compose up --build
```

### Instalação Local

```bash
# Clone o repositório
cd condotrack-api

# Instale as dependências
go mod tidy

# Configure as variáveis de ambiente
cp .env.example .env

# Execute as migrations no MySQL
mysql -u root -p condotrack < migrations/001_initial_schema.sql

# Execute a aplicação
go run cmd/server/main.go
```

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| SERVER_PORT | Porta do servidor | 8000 |
| APP_ENV | Ambiente (development/production) | development |
| DB_HOST | Host do MySQL | localhost |
| DB_PORT | Porta do MySQL | 3306 |
| DB_NAME | Nome do banco | condotrack |
| DB_USER | Usuário do MySQL | root |
| DB_PASS | Senha do MySQL | - |
| ASAAS_API_KEY | Chave da API Asaas | - |
| ASAAS_API_URL | URL da API Asaas | https://sandbox.asaas.com/api/v3 |
| REVENUE_INSTRUCTOR_PERCENT | % do instrutor | 70 |
| REVENUE_PLATFORM_PERCENT | % da plataforma | 30 |

## Endpoints da API

### Health Check
- `GET /api/v1/health` - Status da aplicação

### Gestores
- `GET /api/v1/gestores` - Lista todos os gestores
- `GET /api/v1/gestores/:id` - Busca gestor por ID

### Contratos
- `GET /api/v1/contratos` - Lista todos os contratos
- `GET /api/v1/contratos?gestor_id=X` - Filtra por gestor
- `GET /api/v1/contratos/:id` - Busca contrato por ID
- `POST /api/v1/contratos` - Cria novo contrato

### Auditorias
- `GET /api/v1/audits` - Lista todas as auditorias
- `GET /api/v1/audits?contract_id=X` - Filtra por contrato
- `GET /api/v1/audits/:id` - Busca auditoria por ID
- `GET /api/v1/audits/meta?contract_id=X` - Metadados de auditoria
- `POST /api/v1/audits` - Cria nova auditoria

### Matrículas
- `GET /api/v1/enrollments` - Lista todas as matrículas
- `GET /api/v1/enrollments?student_id=X` - Filtra por aluno
- `GET /api/v1/enrollments/:id` - Busca matrícula por ID
- `POST /api/v1/enrollments` - Cria nova matrícula

### Pagamentos
- `POST /api/v1/payments/customer` - Cria cliente no Asaas
- `POST /api/v1/payments/pix` - Cria pagamento PIX
- `POST /api/v1/payments/boleto` - Cria pagamento Boleto
- `POST /api/v1/payments/card` - Cria pagamento Cartão
- `GET /api/v1/payments/:id/status` - Status do pagamento
- `GET /api/v1/payments/simulate-split` - Simula divisão de receita

### Checkout
- `POST /api/v1/checkout` - Cria checkout completo
- `GET /api/v1/checkout/:id/status` - Status do checkout

### Webhooks
- `POST /api/v1/webhooks/asaas` - Webhook do Asaas

### Certificados
- `GET /api/v1/certificados/:aluno_id` - Certificados do aluno
- `GET /api/v1/certificados/validate/:code` - Valida certificado
- `POST /api/v1/certificados/generate` - Gera certificado

### Imagens
- `GET /api/v1/images` - Lista imagens
- `POST /api/v1/images` - Upload de imagem
- `DELETE /api/v1/images/:filename` - Remove imagem

## Exemplos de Uso

### Health Check
```bash
curl http://localhost:8000/api/v1/health
```

### Listar Gestores
```bash
curl http://localhost:8000/api/v1/gestores
```

### Criar Auditoria
```bash
curl -X POST http://localhost:8000/api/v1/audits \
  -H "Content-Type: application/json" \
  -d '{
    "contract_id": "cont-001",
    "auditor_name": "João Silva",
    "score": 87.5,
    "observations": "Auditoria mensal"
  }'
```

### Criar Pagamento PIX
```bash
curl -X POST http://localhost:8000/api/v1/payments/pix \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cus_xxxxx",
    "value": 199.90,
    "description": "Curso de Gestão Condominial"
  }'
```

### Simular Divisão de Receita
```bash
curl "http://localhost:8000/api/v1/payments/simulate-split?value=100&method=pix"
```

## Taxas de Pagamento (Asaas)

| Método | Taxa |
|--------|------|
| PIX | 0,99% |
| Boleto | R$ 2,99 fixo |
| Cartão | 2,99% + R$ 0,49 |

## Divisão de Receita

Por padrão:
- **Instrutor**: 70%
- **Plataforma**: 30%

A divisão é calculada sobre o valor líquido (após taxas do gateway).

## Docker

### Build
```bash
docker build -t condotrack-api .
```

### Executar com Docker Compose
```bash
docker-compose up -d
```

### Serviços disponíveis
- **API**: http://localhost:8000
- **phpMyAdmin**: http://localhost:8080
- **MySQL**: localhost:3306

## Testes

```bash
# Executar todos os testes
go test ./...

# Executar com cobertura
go test -cover ./...
```

## Licença

Proprietário - CondoTrack © 2024
