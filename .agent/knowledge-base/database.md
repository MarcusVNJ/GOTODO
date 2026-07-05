# Banco de Dados

## PostgreSQL

- Imagem Docker: `postgres:<versão>-alpine`
- Connection pool gerenciado por `pgxpool` (`jackc/pgx/v5`)
- Migrations gerenciadas por `golang-migrate/migrate`

## Convenção de Schema

### Tipos ENUM

Tipos ENUM são criados no PostgreSQL para representar os enums do domínio:

```sql
CREATE TYPE <entity>_<field> AS ENUM ('VALUE1', 'VALUE2', ...);
```

### Tabelas

```sql
CREATE TABLE IF NOT EXISTS <entity>s (
    id          CHAR(20) PRIMARY KEY,
    -- campos de negócio
    status      <entity>_<field> NOT NULL,
    priority    INT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_<entity>s_<field> ON <entity>s(<field>);
```

### Soft Delete

A coluna `deleted_at` implementa o padrão de **soft delete**:
- Registro ativo: `deleted_at = NULL`
- Registro "deletado": `deleted_at = now()`
- A operação `Delete` no repositório executa `UPDATE SET deleted_at = now() WHERE id = $1`

### Audit Fields

Todas as entidades possuem:
- `id` — CHAR(20), gerado via xid
- `created_at` — TIMESTAMP WITH TIME ZONE, NOT NULL
- `updated_at` — TIMESTAMP WITH TIME ZONE, NOT NULL
- `deleted_at` — TIMESTAMP WITH TIME ZONE, nullable (soft delete)

## Migrations

Scripts SQL puros em `migrations/` na raiz do projeto (fora do `internal/`, pois o Core não deve saber sobre SQL).

### Convenção de nomenclatura
- `<numero>_<descricao>.up.sql` — aplica a mudança
- `<numero>_<descricao>.down.sql` — reverte a mudança

### Comando para executar
```bash
migrate -database "postgres://<user>:<pass>@localhost:5432/<db>?sslmode=disable" -path migrations up
```

## Connection Pool

```go
func InitDB(cfg *AppConfig) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }
    return pgxpool.NewWithConfig(context.Background(), config)
}
```

## Query Builder

O `QueryBuilder` usa **Squirrel** com placeholder estilo PostgreSQL (`sq.Dollar`):

| Método | SQL Gerado |
|---|---|
| `QueryInsert(entity)` | `INSERT INTO <table> (...) VALUES ($1, ...)` |
| `QueryExistsById(id)` | `SELECT 1 FROM <table> WHERE id = $1 AND deleted_at IS NULL LIMIT 1` |
| `QueryFindById(id)` | `SELECT * FROM <table> WHERE id = $1` |
| `QueryFindAll(filters)` | `SELECT * FROM <table> WHERE ... ORDER BY created_at DESC` |
| `QueryUpdate(entity)` | `UPDATE <table> SET ... WHERE id = $N` |
| `QueryDelete(id)` | `UPDATE <table> SET deleted_at = now() WHERE id = $1` |

### Filtros Dinâmicos no FindAll
- Filtros são aplicados condicionalmente baseados nos parâmetros recebidos
- Ordenação padrão: `ORDER BY created_at DESC`
- **Sempre incluir** `WHERE deleted_at IS NULL` nas queries que leem dados de entidades com soft delete
- `QueryExistsById` também filtra `deleted_at IS NULL` para que registros deletados não sejam considerados existentes

## Mappers

### Domain → Entity
`mappers.DomainToEntity(model *models.<Entity>) → entity.<Entity>Entity`
- Extrai campos via getters públicos
- Preserva todos os campos incluindo timestamps

### Entity → Domain
`mappers.EntityToDomain(entity *<Entity>Entity) → *models.<Entity>`
- Reconstrói `Audit` via `NewAuditInit()` e entidade via `New<Entity>Init()`
- Bypass de validação (dados do DB já são válidos)

## Entity Struct

Structs de mapeamento DB com tags `db`:

```go
type <Entity>Entity struct {
    ID          string       `db:"id"`
    // campos de negócio
    CreatedAt   time.Time    `db:"created_at"`
    UpdatedAt   time.Time    `db:"updated_at"`
    DeletedAt   *time.Time   `db:"deleted_at"`
}
```