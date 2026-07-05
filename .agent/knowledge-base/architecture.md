# Arquitetura Hexagonal (Ports & Adapters)

## Princípio Fundamental

As regras de negócio devem ser o centro da aplicação. As dependências devem sempre apontar **de fora para dentro**. O domínio (Core) não sabe que existe um banco de dados ou que está respondendo a requisições HTTP — ele se comunica com o mundo externo puramente através de interfaces (Ports).

## Camadas da Arquitetura

```
cmd/api/                     → Composition Root (fx.New + DI wiring)
internal/
  config/                    → Configuração (env vars, DB pool, logger)
  core/                      → HEXÁGONO - Regras de negócio puras
    base/                    → Abstrações genéricas (Usecase base com Call)
    enums/                   → Enums do negócio
    exceptions/              → Exceções tipadas (NUNCA importa net/http)
    models/                  → Entidades de domínio (campos privados + getters)
    ports/                   → Interfaces (contratos de entrada/saída)
    usecase/                 → APENAS i_usecase.go (interface IUsecase)
  app/                       → APPLICATION LAYER - Casos de uso + DTOs
    <entity>/
      dto/                   → Command/Query/Result DTOs (puros, sem tags, um por arquivo)
        create_<entity>_command.go
        create_<entity>_result.go
        update_<entity>_command.go
        get_<entity>_query.go
        list_<entity>_query.go
      *_uc.go                → Implementações de UseCases
      *_test.go              → Testes unitários
      *_mock_test.go         → Mocks do repositório
  adapters/                  → ADAPTADORES - Integração com mundo externo
    in/http/
      handlers/<entity>/     → Resources HTTP + módulo FX
      dto/request/           → DTOs HTTP com tags Huma (validação)
      dto/response/          → DTOs HTTP com tags JSON
      middlewares/           → HandlerException wrapper
      router/                → Interface RouteRegister + helper AsRoute
      server/                → Setup Chi + Huma, lifecycle FX
    out/                      # Repositório, Entities, Mappers, QueryBuilders
      repository/
        entity/              → Modelos exclusivos de DB (tags db)
        mappers/             → Domain ↔ Entity
        query_builder/       → SQL gerado com Squirrel
        *_postgres_impl.go   → Implementação concreta do repositório
        module.go            → fx.Module do repositório

migrations/                  → Scripts SQL Up/Down (golang-migrate)
```

## Fronteiras Arquiteturais

### Core (Hexágono) - `internal/core/`
- **NÃO importa** nenhuma biblioteca de infraestrutura (`net/http`, `pgx`, `json`, `fx`, etc.)
- Não conhece `uber-go/fx` — a injeção é feita externamente em `cmd/api/di/`
- Contém apenas: modelos, enums, exceções, interfaces (ports)
- **O pacote `core/exceptions/codes/` NUNCA importa `net/http`** — use raw ints (400, 404, 422, 500)
- **O pacote `ports/`** declara-se como `package ports` (nome do pacote = nome do diretório)
- `core/usecase/` contém **apenas** a interface `IUsecase` — implementações estão em `app/`

### Application Layer - `internal/app/`
- Contém os **UseCases** (implementações concretas) e **Command/Query DTOs**
- DTOs são structs Go puras — **sem tags de serialização** (`json`, `db`, `query`, `path`)
- A validação de domínio ocorre **dentro do UseCase** (chamando `models.New<Entity>()`)
- Depende apenas do Core (`models`, `ports`, `exceptions`, `base`)
- **NÃO importa** `net/http`, `fx`, `pgx`, `json`

### Adaptadores Inbound - `internal/adapters/in/`
- Traduzem requisições HTTP em Commands/Queries para a Application Layer
- **NUNCA criam modelos de domínio** — convertem HTTP DTO → Command → chamam UseCase
- Huma + Chi lidam com parsing, validação de input e serialização

### Adaptadores Outbound - `internal/adapters/out/`
- Implementam as interfaces (Ports) definidas no Core
- Traduzem modelos de domínio para entidades de banco de dados
- Erros de infra são envelopados com `oops` para stack trace rica

### Composition Root - `cmd/api/`
- Ponto de entrada e orquestração de DI
- `main.go` monta o contêiner FX
- `di/usecases.go` faz o wire dos UseCases importando de `internal/app/`

## Fluxo de uma Requisição (Ex: Criar Tarefa)

```
Request HTTP POST /api/task
  │
  ▼
[Chi Router + Huma API]
  │  - Valida corpo contra o schema OpenAPI
  │  - Desserializa JSON → request.CreateTaskRequest (tags Huma)
  │
  ▼
[Handler (CreateTaskResource)]
  │  - Converte request.CreateTaskRequest → app_task.CreateTaskCommand
  │  - Command é struct pura, sem tags HTTP
  │  - HandlerException wrapper captura erros
  │
  ▼
[IUsecase[CreateTaskCommand, CreateTaskResult].Execute(ctx, cmd)]
  │
  ▼
[base.Usecase.Call(ctx, cmd, businessLogic)]
  │  - Adiciona contexto do payload ao erro
  │
  ▼
[Use Case (CreateTaskUC.createTask)]
  │  - Chama models.NewTask(cmd.Title, cmd.Description, cmd.Priority)
  │  - Validação de domínio ocorre AQUI (não no handler)
  │  - Se válido, chama repository.Save(ctx, task)
  │
  ▼
[ports.TaskRepository (interface no Core)]
  │
  ▼
[PostgresTaskRepository]
  │  - Domain → Entity (via TaskMapper)
  │  - SQL gerado pelo TaskQueryBuilder (Squirrel)
  │  - Executa no pool pgx
  │
  ▼
[PostgreSQL Database]

Resposta sobe pela pilha: Result → Response DTO → JSON HTTP
```

## Princípios Chave

1. **Dependency Inversion**: Core define interfaces (Ports), Adaptadores implementam
2. **Single Responsibility**: Cada camada tem papel único e bem definido
3. **Open/Closed**: Novos adaptadores são adicionados sem modificar o Core
4. **Fail-Fast**: Configuração ausente impede o início da aplicação
5. **Clean Error Flow**: Erros de negócio (4xx) seguem caminho diferente de erros técnicos (5xx)
6. **Core imaculado**: Zero imports de infraestrutura (`net/http`, `pgx`, `json`, `fx`)
7. **Application Layer**: UseCases + Commands/Queries puros, separados do domínio
8. **Handlers não criam modelos**: Handlers convertem HTTP DTO → Command; UseCase cria o modelo

## Anti-Padrões (PROIBIDO)

| Anti-Padrão | Exemplo | Correção |
|---|---|---|
| `net/http` no Core | `import "net/http"` em `exceptions/` | Use raw ints: `400`, `404`, `422`, `500` |
| Pacote ≠ diretório | `package repository` em `ports/` | `package ports` |
| Handler cria modelo de domínio | `task := models.NewTask(...)` no handler | Handler cria Command; UseCase cria modelo |
| Validação de negócio no adapter | `ToModel()` valida no request DTO | UseCase chama `NewTask()` internamente |
| `createdAt` zerado no update | `NewAuditInit(id, time.Time{}, ...)` | Buscar task existente e preservar `createdAt` |
| Query sem filtro soft delete | `WHERE id = $1` sem `deleted_at IS NULL` | `WHERE id = $1 AND deleted_at IS NULL` |
| Constructor sem retorno de erro | `func NewTask(...) *Task` | `func NewTask(...) (*Task, error)` |
