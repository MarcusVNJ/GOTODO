# Padrões de Injeção de Dependência

## Framework: uber-go/fx

Os projetos usam `uber-go/fx` para DI com **módulos descentralizados** que respeitam as fronteiras da Arquitetura Hexagonal.

## Princípios

1. **Core não importa fx**: O `internal/core` é imaculado — não conhece a biblioteca de DI
2. **App não importa fx**: O `internal/app` também é puro — não conhece DI
3. **Adaptadores são autocontidos**: Cada adaptador exporta seu próprio `fx.Module`
4. **Wire do App na borda**: A orquestração de DI dos UseCases fica em `cmd/api/di/usecases.go`

## Composição no Entry Point

```go
// cmd/api/main.go
fx.New(
    config.Module,
    repository_impl.Module,
    di.UsecaseModule,
    taskHandlers.Module,
    server.Module,
).Run()
```

## Padrões FX Utilizados

### 1. fx.Provide vs fx.Invoke

- **`fx.Provide`** (Lazy): Registra a "receita" de construção. Só executa se alguém depender do tipo.
  ```go
  fx.Provide(NewPostgresEntityRepository)
  ```

- **`fx.Invoke`** (Eager): Executa obrigatoriamente na inicialização.
  ```go
  fx.Invoke(StartHTTPServer)
  ```

### 2. Inversão de Dependência com fx.As

O adaptador fornece uma struct concreta, mas o Core exige uma interface:

```go
// repository/module.go
fx.Provide(
    fx.Annotate(
        NewPostgresEntityRepository,
        fx.As(new(ports.EntityRepository)),  // interface definida no Core
    ),
)
```

### 3. Named Instances para Generics

Quando dois UseCases retornam a mesma interface genérica, usamos tags nomeadas para evitar ambiguidade:

```go
// cmd/api/di/usecases.go
fx.Provide(
    fx.Annotate(
        app_entity.NewCreateEntityUC,
        fx.As(new(usecase.IUsecase[entitydto.CreateEntityCommand, entitydto.CreateEntityResult])),
        fx.ResultTags(`name:"createEntityUC"`),
    ),
    fx.Annotate(
        app_entity.NewUpdateEntityUC,
        fx.As(new(usecase.IUsecase[entitydto.UpdateEntityCommand, struct{}])),
        fx.ResultTags(`name:"updateEntityUC"`),
    ),
)

// handlers/module.go - consumidor
fx.Provide(
    fx.Annotate(
        NewCreateEntityResource,
        fx.ParamTags(`name:"createEntityUC"`),
    ),
)
```

> **IMPORTANTE**: Os tipos genéricos nos `fx.As` devem bater exatamente com os tipos declarados nos handlers. Se o handler usa `IUsecase[entitydto.CreateEntityCommand, entitydto.CreateEntityResult]`, o `fx.As` no provider deve usar os mesmos tipos. O import do subpacote `dto/` usa alias como `entitydto`.

### 4. Value Groups para Rotas Dinâmicas

Handlers não são acoplados ao servidor. Eles se registram em um grupo:

```go
// Handler Module
fx.Provide(
    fx.Annotate(
        NewCreateEntityResource,
        fx.ResultTags("group:\"routes\""),
    ),
)

// Server Module - injeta o grupo
type RouteParams struct {
    fx.In
    Routes []router.RouteRegister `group:"routes"`
}
```

O servidor itera sobre o slice e registra todas as rotas.

### 5. Lifecycle Hooks

O servidor HTTP usa lifecycle hooks do FX para graceful shutdown:

```go
fx.Invoke(func(lc fx.Lifecycle, r *chi.Mux, cfg *config.AppConfig) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error { /* start server */ },
        OnStop:  func(ctx context.Context) error { /* shutdown server */ },
    })
})
```

## Estrutura de Módulos

| Módulo | Localização | Fornece | Invoca |
|---|---|---|---|
| `config.Module` | `internal/config/` | `*AppConfig`, `*pgxpool.Pool` | `InitLogger` |
| `repository_impl.Module` | `internal/adapters/out/repository/` | Interface de Port (via `fx.As(new(ports.<Entity>Repository))`) | — |
| `di.UsecaseModule` | `cmd/api/di/` | `IUsecase` com Named Tags (importa de `internal/app/`) | — |
| `handlers.Module` | `internal/adapters/in/http/handlers/<entity>/` | `RouteRegister` (via Value Group) | — |
| `server.Module` | `internal/adapters/in/http/server/` | `*chi.Mux`, `huma.API` | `RegisterRoutes`, `StartHTTPServer` |

## Diagrama de Dependências FX

```
main.go (Composition Root)
  ├── config.Module
  │     ├── Provides: *AppConfig, *pgxpool.Pool
  │     └── Invokes: InitLogger
  ├── repository_impl.Module
  │     └── Provides: ports.EntityRepository (fx.As interface)
  ├── di.UsecaseModule
  │     └── Provides: IUsecase[Command, Result] com Named Tags
  │                    (importa construtores de internal/app/<entity>/)
  ├── handlers.Module
  │     └── Provides: RouteRegister (group:"routes")
  │                    (importa construtores do próprio pacote handlers)
  └── server.Module
        ├── Provides: *chi.Mux, huma.API
        └── Invokes: RegisterRoutes, StartHTTPServer
```

## Como adicionar uma nova entidade

Para cada nova entidade, adicione um bloco `fx.Annotate` em `cmd/api/di/usecases.go`:

```go
fx.Annotate(
    app_entity.NewCreateEntityUC,
    fx.As(new(usecase.IUsecase[entitydto.CreateEntityCommand, entitydto.CreateEntityResult])),
    fx.ResultTags(`name:"createEntityUC"`),
),
```

E no `handlers/<entity>/module.go`:

```go
fx.Provide(
    router.AsRoute(NewCreateEntityResource, fx.ParamTags(`name:"createEntityUC"`)),
)
```

> **Nota**: Os tipos `entitydto.CreateEntityCommand` vêm de `internal/app/<entity>/dto/` (import alias `entitydto`). Os UseCases (`app_entity.NewCreateEntityUC`) vêm de `internal/app/<entity>/`.
