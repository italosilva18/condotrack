# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run

```bash
# Verify compilation (fast check)
go build ./...

# Run locally (requires MySQL + MinIO running)
go run ./cmd/server/main.go

# Docker (from parent directory)
docker compose up --build api -d

# Health check
curl http://localhost:8000/api/v1/health
```

## Architecture

Go 1.22 + Gin + SQLX + MySQL. Clean Architecture:

```
cmd/server/main.go         → entry point
internal/config/            → env vars → Config struct
internal/domain/entity/     → domain structs (db: + json: tags)
internal/domain/repository/ → repository interfaces
internal/domain/gateway/    → PaymentGateway interface + canonical types
internal/infrastructure/    → implementations (MySQL repos, JWT, Asaas, MercadoPago, MinIO)
internal/usecase/           → business logic (one package per module)
internal/delivery/http/     → handlers, middleware, router.go
pkg/response/               → standard JSON response helpers
```

**router.go is the central file.** `NewRouter()` wires all dependencies. `Setup()` registers all routes with middleware. Start here for any feature work.

## Key Patterns

### Adding a Module
entity → repository interface → MySQL implementation → use case → handler → wire in router.go

### Error Handling
```go
// Internal errors: logs real error, returns generic message to client
response.SafeInternalError(c, "Failed to create X", err)

// Known errors: safe static message
response.BadRequest(c, "Invalid email format")
response.NotFound(c, "User not found")
```
Never concatenate `err.Error()` into client-facing responses.

### Entity Struct Tags
```go
type Entity struct {
    ID        string     `db:"id" json:"id"`              // db: must match column exactly
    Name      string     `db:"name" json:"name"`
    Phone     *string    `db:"phone" json:"phone,omitempty"` // pointer = nullable
    CreatedAt time.Time  `db:"created_at" json:"created_at"`
}
```

### Handler Structure
```go
func (h *Handler) Create(c *gin.Context) {
    var req entity.CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request: "+err.Error())
        return
    }
    result, err := h.usecase.Create(c.Request.Context(), req)
    if err != nil {
        response.SafeInternalError(c, "Failed to create", err)
        return
    }
    response.Created(c, result)
}
```

### Import Aliases
```go
authUseCase "github.com/condotrack/api/internal/usecase/auth"
infraRepo "github.com/condotrack/api/internal/infrastructure/repository"
infraAuth "github.com/condotrack/api/internal/infrastructure/auth"
```

## Security

- **Auth**: `middleware.AuthMiddleware(jwtManager)` on all data routes (17 groups)
- **Optional**: `middleware.OptionalAuth(jwtManager)` for legacy/public-mixed routes
- **Roles**: `middleware.RequireRole("admin")` for role gating
- **Token blacklist**: Logout blacklists JWT in-memory; `jwtManager.IsBlacklisted(token)`
- **Rate limits**: Login 10/min, Register 5/min, AI 20/min, Global 100/min per IP
- **Body size**: Global `MaxBodySize` + webhook `io.LimitReader(1MB)`
- **Registration**: Always forces `entity.RoleStudent` regardless of request payload
- **CORS**: Whitelist mode via `CORS_ALLOWED_ORIGINS` env var

Public endpoints (no auth): health, login, register, webhooks, coupon validate, portal images GET, certificate validate.

## Payment Gateway

Factory pattern in `infrastructure/external/`. Asaas and MercadoPago both implement `gateway.PaymentGateway` interface. Webhooks parse to canonical `gateway.WebhookEvent` type. Event types: `EventPaymentConfirmed`, `EventPaymentOverdue`, `EventPaymentRefunded`, `EventPaymentDeleted`, `EventPaymentChargeback`.

Revenue splits are created inside DB transactions on payment confirmation. Financial calculations use `roundCents()` (math.Round to 2 decimals).

## DB Roles

MySQL ENUM: `admin, gestor, supervisor, zelador, manutencao, asg, student, instructor`. Go constants: `RoleAdmin`, `RoleManager`, `RoleInstructor`, `RoleStudent`, `RoleUser`. Not all Go constants exist in the DB ENUM - verify before using.

## Troubleshooting

- **"no rows in result set"**: Check `db:` tags match column names exactly
- **"package redeclared"**: Use import aliases (see above)
- **API won't start**: `docker compose logs api` - usually DB connection or missing env var
- **Compilation errors after changes**: Run `go build ./...` before Docker rebuild
