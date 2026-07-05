# Convenções e Padrões

## Nomenclatura

### Pacotes
- **Core**: nomes curtos e descritivos (`models`, `ports`, `enums`, `exceptions`, `base`, `usecase`)
- **Application**: `app/<entity>/` (UseCases + Commands/Queries)
- **Adaptadores**: nomes descritivos da tecnologia (`http`, `repository`, `mappers`, `entity`)
- **Testes**: package `<pacote>_test` (black-box testing) nos testes de UseCase
- **Pacote `ports/` = `package ports`**: o nome do pacote DEVE ser igual ao nome do diretório

### Arquivos
- UseCases: `<acao>_<entity>_uc.go` em `internal/app/<entity>/`
- Commands/Queries DTOs: `dto.go` em `internal/app/<entity>/`
- Implementações de interfaces: `<interface_name>_postgres_impl.go`
- Módulos FX: `module.go` em cada pacote de adaptador
- Testes: `<nome>_test.go` no mesmo diretório do código testado

### Structs e Interfaces
- **Interfaces de Port**: definidas no Core (`ports.<Entity>Repository`, `usecase.IUsecase`)
- **Implementações concretas**: `Postgres<Entity>Repository` (no adapter out)
- **Handlers**: `<Acao><Entity>Resource`
- **Use Cases**: `<Acao><Entity>UC` (em `app/<entity>/`)
- **Command DTOs**: `<Acao><Entity>Command` (ex: `CreateTaskCommand`, `UpdateTaskCommand`)
- **Query DTOs**: `<Acao><Entity>Query` (ex: `GetTaskQuery`, `ListTasksQuery`)
- **Result DTOs**: `<Acao><Entity>Result` (ex: `CreateTaskResult`)
- **Request DTOs (HTTP)**: `<Acao><Entity>Request` (ex: `CreateTaskRequest`)
- **Response DTOs (HTTP)**: `<Acao><Entity>Response` (ex: `GetTaskResponse`)
- **Entities de DB**: `<Entity>Entity`

## Encapsulamento

- Campos de structs de domínio são **privados** (minúscula)
- Acesso via **getters** públicos (ex: `entity.Title()`, `entity.ID()`)
- Construtores como funções públicas (`New<Entity>()`, `NewAudit()`)
- **Construtores com validação DEVEM retornar `(*Entity, error)`**

## Separação de Modelos (Request → Command → Domain → Entity → Response)

O dado atravessa **5 representações** distintas:

1. **Request DTO** (`adapters/in/http/dto/request`): Parsing e validação HTTP (tags `json`, `required`, `maxLength`, `query`, `path`). Tecnologia-específico.
2. **Command / Query** (`app/<entity>/dto.go`): Struct Go pura, **sem tags de serialização**. Tecnologia-agnóstico. Handler converte Request → Command.
3. **Domain** (`core/models`): Lógica de negócio pura. **Sem** tags de serialização. UseCase cria o modelo a partir do Command.
4. **Entity** (`adapters/out/entity`): Mapeamento DB (tags `db`). Mapper converte Domain ↔ Entity.
5. **Response DTO** (`adapters/in/http/dto/response`): Tags JSON para serialização HTTP. Handler converte Domain → Response.

```
[HTTP Request] → Request DTO → Command → Domain Model → Entity → DB
                                                         ↑
[HTTP Response] ← Response DTO ← Domain Model ←──────────┘
```

Conversões:
- `request → Command`: inline no Handler (mapeamento simples de campos)
- `Command → Domain`: `models.New<Entity>(cmd.Title, ...)` dentro do UseCase
- **Proibido**: `Request → Domain` direto (pula o Command)
- `Domain → Entity`: `mappers.DomainToEntity()`
- `Entity → Domain`: `mappers.EntityToDomain()` via `New<Entity>Init()`
- `Domain → Response`: `response.New<Entity>Response(model)`

## Erros

- **Regras de negócio**: `BusinessException` com `BadRequestCode` (4xx)
- **Erros inesperados**: `UnexpectedException` com `UnexpectedCode` (5xx)
- **Erros de infraestrutura**: Envelopados com `oops` no repositório (stack trace rica)
- **Handlers nunca definem status code**: Retornam `nil, err` e o `HandlerException` classifica
- **Códigos de erro NUNCA importam `net/http`**: use raw ints (400, 404, 422, 500)

## Validação

- **Input HTTP**: Validado pelo Huma via tags nas structs Request DTO (`required`, `maxLength`, `minLength`, `maximum`, `minimum`)
- **Regras de negócio**: Validadas nos modelos de domínio (`<Entity>.validade()`), chamado dentro do UseCase
- **Validação NUNCA ocorre no Handler**: Handler só converte Request → Command; UseCase valida

## Passagem de Parâmetros

- **Por valor (padrão)**: Para structs pequenas/médias como Commands, Queries, DTOs
- **Por ponteiro**: Apenas quando precisa mutar a struct original ou para structs grandes (ex: `*models.Task`)

## Injeção de Dependência (uber-go/fx)

### Módulos Descentralizados
Cada adaptador exporta `var Module = fx.Module(...)` autocontido:
- `config.Module` → fornece `*AppConfig`, `*pgxpool.Pool`, invoca `InitLogger`
- `repository_impl.Module` → fornece interface de Port (via `fx.As(new(ports.TaskRepository))`)
- `di.UsecaseModule` → fornece UseCases instanciados de `app/<entity>/` com named instances
- `handlers.Module` → fornece handlers como `RouteRegister` (via `fx.ResultTags("group:\"routes\"")`)
- `server.Module` → fornece `*chi.Mux`, `huma.API`, invoca registro de rotas e start do servidor

### Value Groups para Rotas
- Handlers registram rotas via `fx.ResultTags("group:\"routes\"")`
- Servidor injeta `Routes []router.RouteRegister` via `fx.In`
- Rotas são registradas iterativamente sem acoplamento

### Named Instances para Generics
- UseCases genéricos do mesmo tipo usam `fx.ResultTags("name:\"<name>\"")` para evitar colisão
- Consumers usam `fx.ParamTags("name:\"<name>\"")` para injetar a instância específica

## Logging

- `slog` (stdlib) para logging estruturado
- BusinessException → `slog.Info` (não é erro de sistema)
- UnexpectedException → `slog.Error` com stack trace `oops`
- Logging assíncrono em goroutine separada (não bloqueia a requisição)

## Configuração

- `envconfig` para mapeamento tipado de variáveis de ambiente
- Fail-Fast: variável obrigatória ausente = aplicação não inicia

## Regras de Ouro

| # | Regra | Motivo |
|---|---|---|
| 1 | Core nunca importa `net/http`, `pgx`, `json`, `fx` | Pureza do domínio |
| 2 | `ports/` declara `package ports` | Convenção Go: pacote = diretório |
| 3 | Constructors com validação retornam `(*Entity, error)` | Nunca `nil` silencioso |
| 4 | Handler nunca cria modelo de domínio | Separação de responsabilidades |
| 5 | UseCase recebe Command/Query, cria modelo internamente | Validação no lugar certo |
| 6 | Update preserva `createdAt` original | Integridade dos dados |
| 7 | Queries de leitura filtram `deleted_at IS NULL` | Consistência do soft delete |
| 8 | Códigos de erro usam raw ints, não `http.Status*` | Core não conhece HTTP |
