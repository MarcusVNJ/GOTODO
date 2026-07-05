# Mapa do Codebase

> **Nota**: A estrutura abaixo é um **template** genérico. Substitua `<entity>` pela entidade de negócio do seu projeto. Os nomes de pacotes são fixos, os nomes de arquivos seguem o padrão `<entity>`.

## Estrutura de Pacotes

### `cmd/api/`
- **main.go**: Composition Root — monta o contêiner `fx.New()` com todos os módulos
- **di/usecases.go**: Wire dos UseCases (borda externa, mantém o Core puro de fx). Importa de `internal/app/<entity>/`

### `internal/config/`
- **app_config.go**: Struct `AppConfig` com env tags, funções `LoadConfig()`, `InitDB()`, `InitLogger()`
- **module.go**: Exporta `fx.Module` do config

### `internal/core/base/`
- **usecase.go**: Struct genérica `Usecase[REQ any, RES any]` com método `Call()` que envolve lógica de negócio e adiciona contexto ao erro

### `internal/core/enums/`
- Definições de tipos enum do negócio (ex: `Status`)

### `internal/core/exceptions/`
- **business_exception.go**: `BusinessException` com `Code`, `Message`, `HTTPStatus`; implementa `error` e `huma.StatusError`
- **unexpected_exception.go**: `UnexpectedException` com `Code`, `Message`, `Details`; implementa `error` e `huma.StatusError`
- **error_code.go**: Interface `Exception` com métodos `Code()`, `Message()`, `HTTPStatus()`
- **codes/**: Constantes tipadas organizadas por categoria (BadRequest 4xxxx, Unexpected 5xxxx). **NUNCA importa `net/http`** — usa raw ints (400, 404, 422, 500)

### `internal/core/models/`
- Entidades de negócio com campos privados, getters públicos e validações
- Structs auxiliares como `Audit` para rastreamento temporal (id, createdAt, updatedAt, deletedAt)
- Factories: `New<Entity>() (*Entity, error)`, `New<Entity>WithoutAudit() (*Entity, error)`, `New<Entity>Init()` (reconstrução do DB)
- **Factories de criação com validação DEVEM retornar `(*Entity, error)`**

### `internal/core/ports/`
- **Pacote declarado como `package ports`** (nome do pacote = nome do diretório)
- Interfaces de contrato definidas pelo domínio
- Repository: métodos de CRUD conforme necessidade do domínio (Save, FindByID, ExistByID, FindAll, Update, Delete, etc.)
- `import "github.com/.../core/ports"` → usa-se `ports.TaskRepository`

### `internal/core/usecase/`
- **APENAS `i_usecase.go`**: Interface genérica `IUsecase[REQ any, RES any]` com método `Execute(ctx, REQ) (RES, error)`
- Implementações concretas foram movidas para `internal/app/`

### `internal/app/<entity>/` — NOVA CAMADA
- **dto/**: Subpacote com Command/Query/Result DTOs — structs Go puras, **sem tags de serialização**, uma por arquivo
  - `create_<entity>_command.go`: `Create<Entity>Command`
  - `create_<entity>_result.go`: `Create<Entity>Result`
  - `update_<entity>_command.go`: `Update<Entity>Command`
  - `get_<entity>_query.go`: `Get<Entity>Query`
  - `list_<entity>_query.go`: `List<Entity>Query`
- **create_<entity>_uc.go**: `Create<Entity>UC` — recebe Command, cria modelo, salva
- **update_<entity>_uc.go**: `Update<Entity>UC` — recebe Command, busca existente, preserva createdAt, atualiza
- **delete_<entity>_uc.go**: `Delete<Entity>UC` — recebe ID string, verifica existência, soft delete
- **get_<entity>_by_id_uc.go**: `Get<Entity>ByIdUC` — recebe Query, busca por ID
- **get_all_<entity>_uc.go**: `GetAll<Entity>UC` — recebe Query com filtros, lista
- **_test.go**: Testes unitários (black-box: `package <entity>_test`)
- **_mock_test.go**: Mock do repositório para testes

### `internal/adapters/in/http/server/`
- **server.go**: `NewRouter()` (Chi com middlewares), `NewHumaAPI()`, `RegisterRoutes()`, `StartHTTPServer()`

### `internal/adapters/in/http/router/`
- **route.go**: Interface `RouteRegister`, struct `RouteParams` para FX injection, helper `AsRoute()`

### `internal/adapters/in/http/handlers/<entity>/`
- Resources HTTP (handlers) + `module.go` com `fx.Module`
- Cada handler: `<Acao><Entity>Resource`
- **Handlers convertem HTTP DTO → Command/Query e chamam UseCase.Execute()**
- **Handlers NUNCA criam modelos de domínio** (`models.NewTask()` é proibido aqui)

### `internal/adapters/in/http/dto/request/`
- DTOs de entrada com tags de validação Huma (`required`, `maxLength`, etc.)
- Structs: `<Acao><Entity>Request` (ex: `CreateTaskRequest`)
- **NÃO contêm métodos que criam modelos de domínio** (sem `ToModel()`)

### `internal/adapters/in/http/dto/response/`
- DTOs de saída com tags JSON
- Funções mapper para converter Domain → Response (ex: `NewTaskResponse(task)`)

### `internal/adapters/in/http/middlewares/`
- **handler_exception.go**: Wrapper genérico `HandlerException[I, O]` para classificação e logging de erros

### `internal/adapters/out/repository/`
- **module.go**: `fx.Module` provendo implementação concreta como interface de Port (via `fx.As`)
- **<entity>_repository_postgres_impl.go**: Implementação PostgreSQL com `pgxpool` e QueryBuilder

### `internal/adapters/out/repository/entity/`
- Structs de mapeamento DB com tags `db`

### `internal/adapters/out/repository/mappers/`
- `DomainToEntity()` e `EntityToDomain()` — conversão bidirecional

### `internal/adapters/out/repository/query_builder/`
- Query builder com Squirrel gerando INSERT, SELECT, UPDATE, DELETE (soft delete)
- **Todas as queries de leitura (SELECT) devem filtrar `deleted_at IS NULL` para entidades com soft delete**

### `migrations/`
- Scripts SQL Up/Down versionados (formato: `<numero>_<descricao>.up.sql` / `.down.sql`)

## Convenção de Nomenclatura

| Tipo | Padrão | Exemplo (entidade `Task`) |
|---|---|---|
| Command DTO | `<Acao><Entity>Command` | `CreateTaskCommand`, `UpdateTaskCommand` |
| Query DTO | `<Acao><Entity>Query` | `GetTaskQuery`, `ListTasksQuery` |
| Result DTO | `<Acao><Entity>Result` | `CreateTaskResult` |
| Use Case | `<Acao><Entity>UC` | `CreateTaskUC` |
| Handler Resource | `<Acao><Entity>Resource` | `CreateTaskResource` |
| Request DTO (HTTP) | `<Acao><Entity>Request` | `CreateTaskRequest` |
| Response DTO (HTTP) | `<Acao><Entity>Response` | `GetTaskResponse` |
| Entity DB | `<Entity>Entity` | `TaskEntity` |
| Repository Impl | `Postgres<Entity>Repository` | `PostgresTaskRepository` |
| Mapper | `DomainToEntity` / `EntityToDomain` | — |
| Módulo FX | `module.go` por pacote | — |
| Test file | `<nome>_test.go` (black-box: `package <entity>_test`) | `create_task_uc_test.go` |
| Mock | `Mock<Entity>Repository` | `MockTaskRepository` |

## Regras de Importação

| Camada | Pode importar | NÃO pode importar |
|---|---|---|
| `core/` | stdlib (`context`, `time`, `fmt`), `xid` | `net/http`, `pgx`, `json`, `fx`, `huma`, `chi` |
| `app/` | `core/` (models, ports, exceptions, base, enums) | `net/http`, `pgx`, `json`, `fx`, `huma`, `chi` |
| `adapters/in/` | `core/`, `app/`, `huma`, `chi`, `fx` | — (tudo permitido) |
| `adapters/out/` | `core/`, `pgx`, `squirrel`, `fx` | `huma`, `chi` |
| `cmd/api/` | todos os módulos | — (composition root) |
