# Visão Geral do Projeto

## GOTODO

API RESTful para gerenciamento de tarefas (estilo Kanban) desenvolvida em Go com Arquitetura Hexagonal (Ports & Adapters).

## Stack Tecnológica

- **Linguagem**: Go 1.26+
- **Router HTTP**: `go-chi/chi` + `huma v2` (OpenAPI 3.x automático)
- **Banco de Dados**: PostgreSQL com `pgx/v5` (connection pool)
- **Query Builder**: `Masterminds/squirrel` (SQL builder)
- **Migrações**: `golang-migrate/migrate`
- **Injeção de Dependência**: `uber-go/fx`
- **IDs**: `rs/xid` (20 caracteres, globalmente únicos)
- **Erros**: `samber/oops` (stack trace rica em erros de infra)
- **Logging**: `log/slog` (stdlib)
- **Testes**: `stretchr/testify` (assert + mock)
- **Configuração**: `kelseyhightower/envconfig` + `joho/godotenv`

## Estrutura do Projeto

```
cmd/api/                     # Composition Root
  main.go                    # fx.New() com todos os módulos
  di/usecases.go             # Wire dos UseCases (importa de app/)

internal/
  config/                    # Config (env vars, DB pool, logger)
  core/                      # HEXÁGONO - Domínio puro (zero imports de infra)
    base/                    # Abstrações genéricas
    enums/                   # Status, etc.
    exceptions/              # BusinessException, UnexpectedException
    models/                  # Task, Audit
    ports/                   # TaskRepository (interface) - package ports
    usecase/                 # APENAS i_usecase.go (interface IUsecase)
  app/                       # APPLICATION - UseCases + Commands/Queries
    task/
      dto.go                 # CreateTaskCommand, UpdateTaskCommand, GetTaskQuery, etc.
      create_task_uc.go      # CreateTaskUC
      update_task_uc.go      # UpdateTaskUC
      delete_task_uc.go      # DeleteTaskUC
      get_task_by_id_uc.go   # GetTaskByIdUC
      get_all_tasks_uc.go    # GetAllTasksUC
      *_test.go              # Testes unitários
  adapters/
    in/http/                 # Handlers HTTP, DTOs, Middlewares
    out/                     # Repositório PostgreSQL, Entities, Mappers

migrations/                  # SQL Up/Down
```

## Requisitos de Ambiente

Arquivo `.env` na raiz:

| Variável | Obrigatória | Default | Descrição |
|---|---|---|---|
| `ENVIRONMENT` | Não | `development` | `development` ou `production` |
| `PORT` | Não | `8080` | Porta do servidor HTTP |
| `DATABASE_URL` | **Sim** | — | URL de conexão PostgreSQL |
| `ENABLE_DOCS` | Não | `true` | Habilita `/docs`, `/openapi.json`, `/schemas` |

## Comandos

```bash
go build ./...                          # compilar todos os pacotes
go test ./internal/app/task/...         # testar use cases de task
go test ./internal/...                  # testar tudo
go run cmd/api/main.go                  # iniciar servidor (requer .env)
docker compose up -d                    # iniciar PostgreSQL
migrate -database "$DATABASE_URL" -path migrations up   # executar migrations
```

## Fluxo Principal

```
POST /api/task (JSON)
  → Huma valida contra schema OpenAPI
  → CreateTaskRequest (HTTP DTO com tags de validação)
  → Handler converte → CreateTaskCommand (struct pura, sem tags)
  → CreateTaskUC.Execute(cmd)
    → models.NewTask(cmd.Title, cmd.Description, cmd.Priority)  ← validação AQUI
    → repository.Save(task)
  → Retorna CreateTaskResult{ID: task.ID()}
  → Handler converte → OperationTaskResponse (HTTP DTO com tags JSON)
  → 201 Created
```

## Características

- **Soft Delete**: Registros nunca são removidos fisicamente (`deleted_at = now()`)
- **Tratamento de Erro Centralizado**: Middleware `HandlerException` classifica 4xx/5xx
- **Rotas Dinâmicas**: Handlers se registram via FX Value Groups, sem acoplamento ao servidor
- **Testes Unitários**: UseCases testados com mocks, sem dependência de banco
- **Configuração Fail-Fast**: App não inicia se `DATABASE_URL` estiver ausente
- **Core Imaculado**: `internal/core/` não importa `net/http`, `pgx`, `json`, `fx`, `huma`, `chi`
- **Application Layer**: UseCases + Commands/Queries puros em `internal/app/`
