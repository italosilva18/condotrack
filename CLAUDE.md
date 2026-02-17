# CLAUDE.md - CondoTrack API

Este arquivo fornece contexto para futuras instâncias do Claude Code operando neste repositório.

## Visão Geral do Projeto

CondoTrack API é uma aplicação backend em Go para gestão de condomínios, cursos e auditorias. Implementa Clean Architecture com Gin framework, SQLX para MySQL, e integração com gateway de pagamentos Asaas.

## Comandos Essenciais

### Build e Execução
```bash
# Build local
cd condotrack-api
go mod tidy
go build ./cmd/server

# Executar servidor (porta 8000)
go run ./cmd/server/main.go

# Docker Compose (inclui MySQL, API, MinIO, phpMyAdmin)
docker-compose up --build

# Apenas rebuild da API
docker-compose up --build api
```

### Testes e Validação
```bash
# Verificar compilação
go build ./...

# Health check
curl http://localhost:8000/api/v1/health

# Testar autenticação (credenciais padrão)
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@condotrack.com","password":"admin123"}'
```

### Banco de Dados
```bash
# Executar migrations via Docker
docker exec -i condotrack-mysql mysql -ucondotrack_user -pCondo@2024Docker condotrack < migrations/001_initial_schema.sql
docker exec -i condotrack-mysql mysql -ucondotrack_user -pCondo@2024Docker condotrack < migrations/002_additional_tables.sql

# Acesso direto ao MySQL
docker exec -it condotrack-mysql mysql -ucondotrack_user -pCondo@2024Docker condotrack
```

## Arquitetura do Projeto

```
condotrack-api/
├── cmd/server/main.go           # Entry point da aplicação
├── internal/
│   ├── config/                  # Configuração via env vars
│   ├── domain/
│   │   ├── entity/              # Entidades do domínio (structs)
│   │   └── repository/          # Interfaces dos repositórios
│   ├── infrastructure/
│   │   ├── auth/                # JWT e hash de senhas
│   │   ├── database/            # Conexão MySQL/SQLX
│   │   ├── external/asaas/      # Cliente API Asaas
│   │   └── repository/          # Implementações MySQL
│   ├── usecase/                 # Lógica de negócio
│   └── delivery/http/
│       ├── handler/             # Handlers HTTP (controllers)
│       ├── middleware/          # CORS, Auth, Logger
│       └── router.go            # Configuração de rotas
├── migrations/                  # Scripts SQL
└── docker-compose.yml           # Orquestração containers
```

## Padrões Importantes

### Estrutura de Entidades
```go
type Entity struct {
    ID        string     `db:"id" json:"id"`
    Name      string     `db:"name" json:"name"`
    CreatedAt time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
```
- Tags `db:` devem corresponder EXATAMENTE às colunas do banco
- Use ponteiros para campos nullable (`*string`, `*time.Time`)
- JSON tags podem ter nomes diferentes das colunas DB

### Estrutura de Handlers
```go
func (h *Handler) Action(c *gin.Context) {
    // 1. Parse request
    var req RequestStruct
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"success": false, "error": err.Error()})
        return
    }

    // 2. Call usecase
    result, err := h.usecase.Method(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"success": false, "error": err.Error()})
        return
    }

    // 3. Return response
    c.JSON(200, gin.H{"success": true, "data": result})
}
```

### Convenção de Respostas
Todas as APIs retornam JSON com estrutura:
```json
{"success": true, "data": {...}}
{"success": false, "error": "mensagem"}
```

## Módulos Implementados

| Módulo | Arquivos | Endpoints |
|--------|----------|-----------|
| Auth | `usecase/auth/`, `handler/auth_handler.go` | `/auth/*` |
| Users | `entity/user.go`, `repository/user_mysql.go` | `/users/*` |
| Gestores | `usecase/gestor/`, `handler/gestor_handler.go` | `/gestores` |
| Contratos | `usecase/contrato/`, `handler/contrato_handler.go` | `/contratos/*` |
| Audits | `usecase/audit/`, `handler/audit_handler.go` | `/audits/*` |
| Matrículas | `usecase/matricula/`, `handler/matricula_handler.go` | `/enrollments/*` |
| Courses | `usecase/courses/`, `handler/course_handler.go` | `/courses/*` |
| Tasks | `usecase/tasks/`, `handler/task_handler.go` | `/tasks/*` |
| Suppliers | `usecase/suppliers/`, `handler/supplier_handler.go` | `/suppliers/*` |
| Team | `usecase/team/`, `handler/team_handler.go` | `/team/*` |
| Agenda | `usecase/agenda/`, `handler/agenda_handler.go` | `/agenda/*` |
| Inspections | `usecase/inspection/`, `handler/inspection_handler.go` | `/inspections/*` |
| Payments | `usecase/payment/`, `handler/payment_handler.go` | `/payments/*` |
| Checkout | `usecase/checkout/`, `handler/checkout_handler.go` | `/checkout/*` |
| Settings | `usecase/setting/`, `handler/setting_handler.go` | `/settings/*` |

## Autenticação JWT

```bash
# Registro
POST /api/v1/auth/register
{"email": "...", "password": "...", "nome": "..."}

# Login (retorna token)
POST /api/v1/auth/login
{"email": "...", "password": "..."}

# Usar token em rotas protegidas
Authorization: Bearer <token>
```

Middleware de autenticação: `internal/delivery/http/middleware/auth.go`

## Variáveis de Ambiente

```env
# Server
SERVER_PORT=8000
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=3306
DB_NAME=condotrack
DB_USER=condotrack_user
DB_PASS=Condo@2024Docker

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Asaas
ASAAS_API_KEY=your-api-key
ASAAS_API_URL=https://sandbox.asaas.com/api/v3

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=condotrack
```

## Integração Asaas (Pagamentos)

Cliente em `internal/infrastructure/external/asaas/`:
- `client.go` - HTTP client base
- `customer.go` - Criar/buscar clientes
- `payment.go` - PIX, Boleto, Cartão
- `types.go` - Structs de request/response

Taxas configuradas:
- PIX: 0.99%
- Boleto: R$ 2.99
- Cartão: 2.99% + R$ 0.49

## Tabelas do Banco de Dados

### Principais
- `users` - Usuários do sistema
- `gestores` - Gestores de condomínios
- `contratos` - Contratos de gestão
- `audits` / `audit_items` - Auditorias
- `enrollments` / `enrollment_payments` - Matrículas
- `courses` / `course_modules` / `course_lessons` - Cursos

### Adicionais
- `agenda` - Eventos do calendário
- `tasks` - Tarefas
- `suppliers` - Fornecedores
- `team_members` - Membros da equipe
- `inspections` - Vistorias
- `routine_plans` / `routine_items` - Planos de rotina

## Dicas de Produtividade

1. **Adicionar novo módulo**: Criar na ordem entity → repository interface → repository impl → usecase → handler → router

2. **Conflitos de import**: Quando dois packages têm o mesmo nome, use alias:
   ```go
   authUseCase "github.com/condotrack/api/internal/usecase/auth"
   ```

3. **Mismatch DB/Entity**: Tags `db:` DEVEM corresponder às colunas. Se erro "no rows", verifique nomes das colunas.

4. **Router**: Todas as rotas definidas em `internal/delivery/http/router.go`

5. **Testes locais**: Use `go build ./...` para verificar compilação antes de commitar

## Docker Services

| Service | Porta | Descrição |
|---------|-------|-----------|
| api | 8000 | API Go |
| mysql | 3306 | MySQL 8.0 |
| phpmyadmin | 8080 | Interface DB |
| minio | 9000/9001 | Object storage |
| frontend | 3000 | React app |
| portal | 3001 | Portal alunos |

## Troubleshooting

### Erro: "no rows in result set"
- Verificar se dados existem no banco
- Verificar tags `db:` correspondem às colunas

### Erro: package redeclared
- Usar alias para imports com mesmo nome

### API não responde
```bash
docker-compose logs api
docker-compose restart api
```

### MySQL não conecta
```bash
docker-compose logs mysql
# Verificar health do container
docker-compose ps
```

## Sistema de Configurações

### Página de Configurações (Admin)
Acesse: `http://localhost:3000/settings`

A página de configurações permite gerenciar:
- **Geral**: Nome do sistema, email, tamanho de upload
- **Pagamentos (Asaas)**: API Key, URL, ambiente, webhook token
- **IA (Gemini)**: API Key, habilitar/desabilitar
- **Divisão de Receita**: Percentuais instrutor/plataforma
- **Armazenamento (MinIO)**: Endpoint, access key, secret key
- **Email (SMTP)**: Servidor, porta, credenciais

### API de Configurações
```bash
# Listar todas as configurações (requer token admin)
GET /api/v1/settings

# Atualizar configuração individual
PUT /api/v1/settings/:key
{"value": "novo-valor"}

# Atualizar múltiplas configurações
PUT /api/v1/settings
{"settings": {"key1": "value1", "key2": "value2"}}
```

### Tabela de Configurações (settings)
| Coluna | Tipo | Descrição |
|--------|------|-----------|
| setting_key | VARCHAR(100) | Chave única da configuração |
| setting_value | TEXT | Valor da configuração |
| setting_type | VARCHAR(20) | Tipo: string, number, boolean, secret |
| category | VARCHAR(50) | Categoria: general, payment, ai, revenue, storage, email |
| is_secret | BOOLEAN | Se é um valor sensível (mascarado na UI) |
| is_required | BOOLEAN | Se é obrigatório |

### Notas de Segurança
- Apenas usuários com role "admin" podem acessar as configurações
- Valores secretos (API keys, senhas) são mascarados na resposta da API
- Campos secretos vazios mantêm o valor anterior ao salvar
