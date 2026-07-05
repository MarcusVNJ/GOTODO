# Modelo de Domínio

> **Nota**: Os padrões abaixo são **templates** que devem ser adaptados ao contexto de cada aplicação. Substitua `<Entity>` pelo nome da entidade de negócio do seu projeto (ex: `Task`, `Order`, `Customer`, etc.). Os exemplos concretos são meramente ilustrativos.

## Filosofia

O domínio (Core) é o coração da aplicação. Ele contém regras de negócio puras, **sem qualquer dependência de infraestrutura**. Nenhuma interface, struct ou função dentro do Core pode importar bibliotecas de infraestrutura (`net/http`, `pgx`, `json`, `fx`, `huma`, `chi`, etc.).

## Padrão de Entidade

Cada entidade do domínio segue o mesmo padrão arquitetural. Aplique este template à entidade do seu contexto:

### Audit (Base Temporal)
Struct base embutida (embedded) em todas as entidades que requerem rastreamento temporal:

| Campo | Tipo | Gerador | Descrição |
|---|---|---|---|
| `id` | `string` | `xid.New().String()` (20 chars) | ID único global |
| `createdAt` | `time.Time` | `time.Now()` | Timestamp de criação |
| `updatedAt` | `time.Time` | `time.Now()` | Timestamp da última atualização |
| `deletedAt` | `*time.Time` | `nil` (ativo) | Soft-delete (nil = ativo) |

**Factories**:
- `NewAudit()` → cria com ID e timestamps atuais
- `NewAuditInit(id, createdAt, updatedAt, deletedAt)` → reconstrói a partir do banco
- `UpdatedAudit()` → atualiza `updatedAt` para `time.Now().UTC()`

**Getters**: `ID()`, `CreatedAt()`, `UpdatedAt()`, `DeletedAt()`, `SetID(id)`

### Entidade de Negócio
Cada entidade embeds `Audit` e segue o padrão:

- **Campos privados** com **getters públicos**
- **Factories**:
  - `New<Entity>(...) (*Entity, error)` → cria com validações e Audit novo. **DEVE retornar erro** se validação falhar (nunca `nil` silencioso)
  - `New<Entity>WithoutAudit(...) (*Entity, error)` → cria sem Audit. **DEVE retornar erro**
  - `New<Entity>Init(...)` → reconstrói do banco (bypass de validação, dados já são válidos)
- **Métodos de comportamento**: transições de estado e validações
- **Método `validade()` private**: contém todas as validações de criação

**Exemplo ilustrativo** (para uma entidade `Task` com status Kanban):
```go
type Task struct {
    Audit
    title       string
    description string
    status      enums.Status
    priority    int
}

// DEVE retornar (*Task, error)
func NewTask(title, description string, priority int) (*Task, error) {
    task := &Task{
        Audit:       NewAudit(),
        title:       title,
        description: description,
        priority:    priority,
        status:      enums.Pending,
    }
    if err := task.validade(); err != nil {
        return nil, err  // NUNCA return nil sem erro
    }
    return task, nil
}

// Init: reconstrói do banco, bypassa validação
func NewTaskInit(audit Audit, title, desc string, status enums.Status, priority int) *Task {
    return &Task{Audit: audit, title: title, description: desc, status: status, priority: priority}
}
```

> Os campos, comportamentos e validações dependem inteiramente do domínio da aplicação. O que é fixo é o **padrão estrutural**: campos privados, getters, factories com retorno de erro e métodos de comportamento encapsulando invariantes.

## Quem cria o modelo de domínio?

**A criação do modelo de domínio é responsabilidade do UseCase (Application Layer), NUNCA do Handler (Adapter).**

```
❌ ERRADO:
Handler: task := models.NewTask(...)  // handler cria modelo de domínio
UseCase: uc.Execute(ctx, task)         // use case recebe modelo pronto

✅ CORRETO:
Handler: cmd := app_task.CreateTaskCommand{Title: ..., Priority: ...}
UseCase: task, err := models.NewTask(cmd.Title, cmd.Description, cmd.Priority)
```

O Handler (Adapter) converte o HTTP DTO em um **Command** (struct pura, sem lógica de domínio). O UseCase recebe o Command e internamente cria/valida o modelo de domínio.

## Enums

Enums são definidos como tipos `string` com constantes. O conteúdo do enum depende do domínio da aplicação:

**Exemplo ilustrativo** (uma aplicação de tarefas com fluxo Kanban):
```go
type Status string

const (
    Pending    Status = "PENDING"
    InProcess  Status = "IN_PROCESS"
    Completed  Status = "COMPLETED"
    Cancelled  Status = "CANCELLED"
)
```

> Substitua pelo domínio do seu projeto. Uma aplicação de pedidos poderia ter `OrderStatus` (`PLACED`, `PAID`, `SHIPPED`, `DELIVERED`), uma de usuários teria `UserRole` (`ADMIN`, `MEMBER`, `GUEST`), etc.

Enums do domínio devem corresponder ao tipo ENUM no PostgreSQL para consistência.

## Regras de Negócio

Cada entidade possui validações internas no método `validade()` (private). O conteúdo é específico do domínio, mas o padrão é fixo:

- Campos obrigatórios não podem ser vazios
- Valores devem estar em ranges válidos
- Transições de estado seguem regras restritas
- Status inicial não pode ser terminal

> **O agente deve**: Ao criar uma nova entidade, identificar as regras de negócio do domínio e implementá-las dentro do modelo, nunca fora dele.

## Sistema de Códigos de Erro

### BadRequestCode (4xx)
Constantes tipadas com valores a partir de 400001, uma por regra de negócio violada. O conteúdo é específico do domínio.

**Exemplo de padrão de nomenclatura**:
```go
type BadRequestCode int

const (
    <Entity>TitleEmpty    BadRequestCode = 400001 + iota
    <Entity>Already<State>
    <Entity>Init<State>
    Invalid<Campo>
    UnprocessableEntity
    IdInvalid
    <Entity>NotFound
)
```

**IMPORTANTE**: O método `HTTPStatus()` retorna **raw ints** (400, 404, 422), **NUNCA** `http.Status*`. O pacote `codes/` **não importa `net/http`**.

```go
func (code BadRequestCode) HTTPStatus() int {
    switch code {
    case UnprocessableEntity:
        return 422  // raw int, NÃO http.StatusUnprocessableEntity
    case TaskNotFound:
        return 404
    default:
        return 400
    }
}
```

### UnexpectedCode (5xx)
Constantes tipadas com valores a partir de 500001, para erros técnicos inesperados.

```go
type UnexpectedCode int

const (
    UnexpectedError UnexpectedCode = 500001
)

func (code UnexpectedCode) HTTPStatus() int {
    return 500  // raw int, NÃO http.StatusInternalServerError
}
```

## Ports (Interfaces)

### Repository
Interface genérica de persistência definida no Core (`package ports`). Os métodos dependem do domínio:

```go
package ports

type <Entity>Repository interface {
    Save(ctx context.Context, entity *models.<Entity>) error
    FindByID(ctx context.Context, id string) (*models.<Entity>, error)
    ExistByID(ctx context.Context, id string) (bool, error)
    FindAll(ctx context.Context, statusFilter string, minPriority int) ([]*models.<Entity>, error)
    Update(ctx context.Context, entity *models.<Entity>) error
    Delete(ctx context.Context, id string) error
}
```

> O pacote `ports/` declara-se como `package ports` (nome = diretório). Importa-se como `"github.com/.../core/ports"` e usa-se `ports.TaskRepository`.

### IUsecase
Interface genérica para casos de uso (único arquivo em `core/usecase/i_usecase.go`):

```go
type IUsecase[REQ any, RES any] interface {
    Execute(ctx context.Context, req REQ) (RES, error)
}
```

> Implementações concretas de UseCases estão em `internal/app/<entity>/`, NÃO em `core/usecase/`.

## Application Layer - Commands, Queries e Results

Commands, Queries e Results são DTOs puros definidos em `internal/app/<entity>/dto/`, **um arquivo por conceito** (Single Responsibility):

```
internal/app/task/dto/
├── create_task_command.go   # CreateTaskCommand
├── create_task_result.go    # CreateTaskResult
├── update_task_command.go   # UpdateTaskCommand
├── get_task_query.go        # GetTaskQuery
└── list_tasks_query.go      # ListTasksQuery
```

```go
// internal/app/task/dto/create_task_command.go
package dto

type CreateTaskCommand struct {
    Title       string
    Description string
    Priority    int
}

// internal/app/task/dto/create_task_result.go
package dto

type CreateTaskResult struct {
    ID string
}

// internal/app/task/dto/get_task_query.go
package dto

type GetTaskQuery struct {
    ID string
}
```

> **Regra**: Cada arquivo contém **um único DTO**. Command, Query e Result são separados. Nenhum possui tags `json`, `query`, `path`, `db` ou qualquer tag de infraestrutura. São objetos de transporte puros entre o Adapter e o UseCase. Se trocar de HTTP para gRPC, apenas o Adapter muda — os DTOs permanecem idênticos.

> **Sobre `CreateTaskResult`**: É o output do UseCase — faz parte do contrato `IUsecase[CreateTaskCommand, CreateTaskResult]`. Pertence à Application Layer junto com o Command porque ambos formam o par de entrada/saída do port. O Handler recebe o Result e o converte em Response DTO (HTTP).
