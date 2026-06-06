# AGENTS.md

## Project Overview

Go microservice with Hexagonal Architecture (Ports & Adapters). Kanban task management domain. Module: `github.com/MarcusVNJ/GOTODO`.

## Commands

```bash
go build ./...          # build all packages
go test ./...           # run all tests
go run cmd/api/main.go  # start server (requires .env)
```

No Makefile, no linter config, no CI workflows. `gofmt` is the only formatter.

## Environment

Requires `.env` at root with `ENVIRONMENT`, `PORT`, `DATABASE_URL`, `ENABLE_DOCS`. App fails fast if `DATABASE_URL` is missing. PostgreSQL must be running before startup (`docker compose up -d`).

## Architecture

```
cmd/api/main.go              # Composition Root (fx.New)
cmd/api/di/usecases.go       # UseCase wiring (ResultTags/ParamTags, imports from app/)
internal/core/                # Domain — NO infra imports (models, enums, exceptions, ports)
internal/app/task/            # Application — UseCases + Commands/Queries (also no infra imports)
internal/adapters/in/http/   # Handlers, DTOs, Middlewares
internal/adapters/out/        # Repository impls, Entities, Mappers, QueryBuilders
migrations/                   # SQL Up/Down (golang-migrate)
```

Core (`internal/core/`) and App (`internal/app/`) must never import `net/http`, `pgx`, `json`, `fx`, or any infra library.

## Key Patterns

- **DI**: `uber-go/fx` with named instances for generic UseCases. New entity = new `fx.Annotate` block in `di/usecases.go` (importing from `internal/app/<entity>/`) + new `AsRoute` in handler `module.go`.
- **Handlers**: Huma v2 resources implementing `RouteRegister`, registered via `AsRoute()` with `fx.ResultTags("group:routes")`. Handlers convert HTTP DTO → Command/Query → call UseCase. **Handlers never create domain models**.
- **Commands/Queries**: Pure Go structs in `internal/app/<entity>/dto/` (one file per concept: `create_<entity>_command.go`, `create_<entity>_result.go`, etc.). No serialization tags. Technology-agnostic transport between Handler and UseCase.
- **Error handling**: Handlers return `nil, err`. `HandlerException` classifies: `BusinessException` → 4xx, `oops.OopsError`/generic → 5xx. Never set HTTP status in handlers. Error codes use raw ints (400, 404, 422, 500), never `http.Status*`.
- **Models**: Private fields, public getters. `Audit` struct embedded for `id`, `createdAt`, `updatedAt`, `deletedAt`. Constructors: `New*() (*Entity, error)`, `New*WithoutAudit() (*Entity, error)`, `New*Init()`.
- **Soft delete**: `deleted_at` column. `QueryDelete` sets `deleted_at = NOW()`. All read queries filter `WHERE deleted_at IS NULL`.
- **IDs**: `xid.New().String()` (20 chars), stored as `CHAR(20)`.

## Testing

Tests in `internal/app/task/*_test.go`. Pattern: AAA with `stretchr/testify`. Package `<entity>_test` for black-box testing. Mock: `MockTaskRepository` using `mock.Mock`. Tests create Command/Query objects, not domain models.

```bash
go test ./internal/app/task/...   # test task use cases
```

## Migrations

```bash
migrate -database "postgres://postgres:root@localhost:5432/todo_db?sslmode=disable" -path migrations up
```

## Custom Agents

This project has five custom OpenCode agents. Switch with **Tab** or invoke with `@`:

| Agent | Mode | Purpose |
|---|---|---|
| `planner` | primary | Generates task plans from user stories/bugs following Hexagonal patterns |
| `knowledge-base-creator` | primary | Analyzes codebase and generates 12 KB files in `.agent/knowledge-base/` |
| `software-engineer` | primary | Implements tasks from Planner plans (Core → Outbound → App → Inbound → Orchestrator) |
| `test-engineer` | primary | Creates and maintains unit tests focused on business rules |
| `doc-updater` | primary | Updates knowledge bases after code changes (endpoints, models, error codes, etc.) |

**Workflow** (`@planner` → `@software-engineer` → `@test-engineer` → `@doc-updater`):

1. Run `knowledge-base-creator` first on new projects to generate the KB
2. `@planner` — Transform user stories, bugs, specs into actionable task plans in `.agent/tasks/output/`
3. `@software-engineer` — Implement the plan (Core → Outbound → App → Inbound → Orchestrator)
4. `@test-engineer` — Create unit tests for the business rules identified in the plan
5. `@doc-updater` — Update knowledge bases and project docs (`AGENTS.md`, `ARCHITECTURE.md`, `ONBOARDING.md`, `README.md`) based on code changes (`git diff`)

Full pipeline documented in `.agent/workflows/pipeline.md`.

The KB files are in `.agent/knowledge-base/` (gitignored — regenerated per project). Agent specs are in `.agent/agents/` (also gitignored — referenced by `.opencode/agents/`).

## Common Pitfalls

- Adding `fx` imports inside `internal/core/` or `internal/app/` — DI wiring lives in `cmd/api/di/` only
- Adding `net/http` imports in `internal/core/exceptions/codes/` — use raw ints (400, 404, 422, 500)
- Forgetting `ResultTags`/`ParamTags` naming when two UseCases share the same `IUsecase<REQ, RES>` — causes fx collision
- Returning HTTP status codes in handlers — let `HandlerException` handle classification
- Handlers creating domain models (`models.NewTask(...)` in a handler) — handler converts Request → Command, UseCase creates the model
- Using `New*()` in update flows — use `New*Init()` to bypass validation for DB reconstruction
- Omitting `WHERE deleted_at IS NULL` in read queries (`ExistByID`, `FindByID`, `FindAll`) — soft-deleted rows will leak
- Factories returning `nil` without error on validation failure — factories must return `(*Entity, error)`
- Update zeroing `createdAt` — must fetch existing entity and preserve original `createdAt` in `NewAuditInit`
- Response DTOs use unexported inner struct (`taskResponse`) with exported wrapper (`GetTaskResponse`)