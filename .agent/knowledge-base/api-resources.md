# Recursos da API (Endpoints Existentes)

> **Nota**: Este documento cataloga todos os recursos (endpoints) atualmente implementados no serviço. O agente Planner DEVE consultar esta base antes de planejar uma nova feature para evitar duplicação e entender o que já existe.

---

## Visão Geral

- **Framework**: Huma v2 (sobre chi router) — gera OpenAPI 3.x automaticamente
- **Base Path**: `/api`
- **Prefixo de tags**: Agrupados por domínio (ex: `Tasks`)
- **Documentação**: `/docs` (Swagger UI), `/openapi.json`, `/schemas`
- **Middleware**: `HandlerException` envolve todos os handlers, classificando erros em BusinessException (4xx) ou UnexpectedException (5xx)
- **Soft Delete**: Entidades usam `deleted_at` ao invés de remoção física
- **Fluxo**: Request DTO (HTTP) → Command (app) → UseCase → Domain Model → Repository
- **Handlers NUNCA criam modelos de domínio** — convertem Request → Command, chamam UseCase

---

## Recursos Disponíveis

### Domínio: Task (Tarefas)

Tag: `Tasks`

---

#### POST /api/task — Criar Tarefa

| Campo | Valor |
|---|---|
| **Operation ID** | `create-task` |
| **Summary** | Criar Tarefa |
| **Status de Sucesso** | `201 Created` |
| **UseCase** | `CreateTaskUC` (`IUsecase[CreateTaskCommand, CreateTaskResult]`) |

**Fluxo**: `CreateTaskRequest` (HTTP DTO) → `CreateTaskCommand` (app DTO puro) → `CreateTaskUC.Execute()` → `models.NewTask()` (validação) → `repository.Save()`

**Request DTO** (HTTP — `adapters/in/http/dto/request`):

| Campo | Tipo | JSON Key | Validação | Descrição |
|---|---|---|---|---|
| `Title` | `string` | `"title"` | `minLength:"1"`, `maxLength:"150"` | Titulo da tarefa |
| `Description` | `string` | `"description"` | `maxLength:"500"` | Descrição da tarefa |
| `Priority` | `int` | `"priority"` | `minimum:"1"`, `maximum:"5"` | Prioridade de 1 a 5 |

**Command DTO** (app — `internal/app/task/dto.go`):

```go
type CreateTaskCommand struct {
    Title       string
    Description string
    Priority    int
}
type CreateTaskResult struct {
    ID string
}
```

**Comportamento**:
- Handler converte `CreateTaskRequest.Body` → `CreateTaskCommand`
- UseCase chama `models.NewTask(cmd.Title, cmd.Description, cmd.Priority)` → validação de domínio
- Se válido, gera ID via `xid`, status `PENDING`, timestamps `createdAt` e `updatedAt`
- Salva via repository e retorna `CreateTaskResult{ID: task.ID()}`

**Response DTO**: `OperationTaskResponse`

| Campo | Tipo | JSON Key | Descrição |
|---|---|---|---|
| `Message` | `string` | `"message"` | `"Task criada com sucesso"` |
| `Id` | `string` | `"id"` | ID da tarefa criada |

**Erros Possíveis**:

| Código | Constante | Mensagem | HTTP Status | Condição |
|---|---|---|---|---|
| 400001 | `TaskTitleEmpty` | "task title cannot be empty" | 400 | Título vazio |
| 400004 | `InvalidPriority` | "priority must be between 1 and 5" | 400 | Prioridade fora do range |
| 400003 | `TaskInitDone` | "task cannot start with a completed state" | 400 | Status inicial COMPLETED |
| 400005 | `UnprocessableEntity` | "JSON invalido" | 422 | Payload mal formatado |
| 500001 | `UnexpectedError` | "Erro interno critico no servidor" | 500 | Erro técnico |

---

#### GET /api/task/{id} — Buscar Tarefa por ID

| Campo | Valor |
|---|---|
| **Operation ID** | `get-task-by-id` |
| **Summary** | Buscar Tarefa |
| **Status de Sucesso** | `200 OK` |
| **UseCase** | `GetTaskByIdUC` (`IUsecase[GetTaskQuery, *models.Task]`) |

**Fluxo**: `GetTaskRequest.ID` → `GetTaskQuery{ID: ...}` → `GetTaskByIdUC.Execute()` → `repository.FindByID()` → `*models.Task`

**Request DTO**: `GetTaskRequest`

| Campo | Tipo | Source | Validação | Descrição |
|---|---|---|---|---|
| `ID` | `string` | `path:"id"` | `minLength:"1"`, `maxLength:"36"` | ID da task |

**Query DTO** (app):

```go
type GetTaskQuery struct {
    ID string
}
```

**Response DTO**: `GetTaskResponse`

| Campo | Tipo | JSON Key | Descrição |
|---|---|---|---|
| `Title` | `string` | `"title"` | Título da tarefa |
| `Description` | `string` | `"description"` | Descrição (omitido se vazio) |
| `Status` | `enums.Status` | `"status"` | PENDING, IN_PROCESS, COMPLETED, CANCELLED |
| `Priority` | `int` | `"priority"` | Prioridade 1-5 |

**Erros Possíveis**:

| Código | Constante | Mensagem | HTTP Status | Condição |
|---|---|---|---|---|
| 400007 | `TaskNotFound` | "Task nao existe" | 404 | ID não encontrado |
| 500001 | `UnexpectedError` | "Erro interno critico no servidor" | 500 | Erro técnico |

---

#### PUT /api/task — Atualizar Tarefa

| Campo | Valor |
|---|---|
| **Operation ID** | `update-task` |
| **Summary** | Atualizar Tarefa |
| **Status de Sucesso** | `200 OK` |
| **UseCase** | `UpdateTaskUC` (`IUsecase[UpdateTaskCommand, struct{}]`) |

**Fluxo**: `UpdateTaskRequest.Body` → `UpdateTaskCommand` → `UpdateTaskUC.Execute()` → verifica existência → busca task atual para preservar `createdAt` → `models.NewTaskInit()` → `repository.Update()`

**Request DTO**: `UpdateTaskRequest`

| Campo | Tipo | JSON Key | Validação | Descrição |
|---|---|---|---|---|
| `Id` | `string` | `"id"` | `minLength:"1"` | ID da task |
| `Title` | `string` | `"title"` | `minLength:"1"`, `maxLength:"150"` | Título |
| `Description` | `string` | `"description"` | `maxLength:"500"` | Descrição |
| `Status` | `enums.Status` | `"status"` | (enum) | Status: PENDING, IN_PROCESS, COMPLETED, CANCELLED |
| `Priority` | `int` | `"priority"` | `minimum:"1"`, `maximum:"5"` | Prioridade |

**Command DTO** (app):

```go
type UpdateTaskCommand struct {
    ID          string
    Title       string
    Description string
    Status      enums.Status
    Priority    int
}
```

**Comportamento**:
- Handler converte `UpdateTaskRequest.Body` → `UpdateTaskCommand`
- UseCase verifica existência por `ExistByID(cmd.ID)` (com filtro `deleted_at IS NULL`)
- Busca task existente via `FindByID(cmd.ID)` para **preservar `createdAt` original**
- Cria `NewTaskInit(auditComCreatedAtOriginal, cmd.Title, cmd.Description, cmd.Status, cmd.Priority)`
- Atualiza via `repository.Update()`

**Response DTO**: `OperationTaskResponse`

| Campo | Tipo | JSON Key | Descrição |
|---|---|---|---|
| `Message` | `string` | `"message"` | `"Task atualizada com sucesso"` |
| `Id` | `string` | `"id"` | Omitempty |

**Erros Possíveis**:

| Código | Constante | Mensagem | HTTP Status | Condição |
|---|---|---|---|---|
| 400007 | `TaskNotFound` | "Task nao existe" | 404 | ID não encontrado |
| 500001 | `UnexpectedError` | "Erro interno critico no servidor" | 500 | Erro técnico |

---

#### DELETE /api/task/{id} — Excluir Tarefa (Soft Delete)

| Campo | Valor |
|---|---|
| **Operation ID** | `delete-task` |
| **Summary** | Excluir Tarefa |
| **Status de Sucesso** | `204 No Content` |
| **UseCase** | `DeleteTaskUC` (`IUsecase[string, struct{}]`) |

**Request DTO**: `DeleteTaskRequest`

| Campo | Tipo | Source | Validação | Descrição |
|---|---|---|---|---|
| `ID` | `string` | `path:"id"` | `minLength:"1"`, `maxLength:"36"` | ID da task |

**Comportamento**:
- Handler valida ID não vazio (`validateId()`)
- UseCase verifica existência por ID (`ExistByID` com filtro `deleted_at IS NULL`)
- Soft delete: `UPDATE SET deleted_at = now() WHERE id = $1`

**Response DTO**: `OperationTaskResponse`

| Campo | Tipo | JSON Key | Descrição |
|---|---|---|---|
| `Message` | `string` | `"message"` | `"Task deletada com sucesso"` |
| `Id` | `string` | `"id"` | Omitempty |

**Erros Possíveis**:

| Código | Constante | Mensagem | HTTP Status | Condição |
|---|---|---|---|---|
| 400006 | `IdInvalid` | "ID invalido" | 400 | ID vazio |
| 400007 | `TaskNotFound` | "Task nao existe" | 404 | ID não encontrado |
| 500001 | `UnexpectedError` | "Erro interno critico no servidor" | 500 | Erro técnico |

---

#### GET /api/task — Listar Tarefas

| Campo | Valor |
|---|---|
| **Operation ID** | `list-tasks` |
| **Summary** | Listar Tarefas |
| **Status de Sucesso** | `200 OK` |
| **UseCase** | `GetAllTasksUC` (`IUsecase[ListTasksQuery, []*models.Task]`) |

**Fluxo**: `ListTaskRequest` → `ListTasksQuery` → `GetAllTasksUC.Execute()` → `repository.FindAll()`

**Request DTO**: `ListTaskRequest`

| Campo | Tipo | Source | Validação | Descrição |
|---|---|---|---|---|
| `Status` | `string` | `query:"status"` | — | Filtrar por status |
| `MinPriority` | `int` | `query:"min_priority"` | — | Filtrar por prioridade mínima |

**Query DTO** (app):

```go
type ListTasksQuery struct {
    Status      string
    MinPriority int
}
```

**Response DTO**: `ListTaskResponse` — array de `{id, title, priority}`

---

## Modelo de Domínio: Task

### Enums

```go
type Status string
const (
    Pending    Status = "PENDING"
    InProcess  Status = "IN_PROCESS"
    Completed  Status = "COMPLETED"
    Cancelled  Status = "CANCELLED"
)
```

### Campos da Entidade

| Campo | Tipo Go | Tipo DB | Obrigatório | Descrição |
|---|---|---|---|---|
| `id` | `string` | `CHAR(20)` | Auto (xid) | Identificador único |
| `title` | `string` | `VARCHAR(150)` | Sim | Título da tarefa (1-150 chars) |
| `description` | `string` | `TEXT` | Não | Descrição (max 500 chars) |
| `status` | `enums.Status` | `task_status` ENUM | Sim | PENDING, IN_PROCESS, COMPLETED, CANCELLED |
| `priority` | `int` | `INT` | Sim | Prioridade (1-5) |
| `created_at` | `time.Time` | `TIMESTAMPTZ` | Auto | Data de criação |
| `updated_at` | `time.Time` | `TIMESTAMPTZ` | Auto | Data de atualização |
| `deleted_at` | `*time.Time` | `TIMESTAMPTZ` | Não | Soft delete (null = ativo) |

### Regras de Negócio Atuais

- **RN01**: Title não pode ser vazio (TaskTitleEmpty)
- **RN02**: Priority deve ser entre 1 e 5 (InvalidPriority)
- **RN03**: Status inicial não pode ser COMPLETED (TaskInitDone)
- **RN04**: Task já completa não pode ser completada novamente (TaskAlreadyDone — método `Complete()` existe mas não é exposto via endpoint)
- **RN05**: Soft delete — exclusão lógica via `deleted_at`

### Repository Methods Disponíveis

| Método | Descrição |
|---|---|
| `Save(ctx, *Task)` | Cria nova tarefa |
| `FindByID(ctx, id)` | Busca por ID |
| `ExistByID(ctx, id)` | Verifica existência (filtra `deleted_at IS NULL`) |
| `FindAll(ctx, statusFilter, minPriority)` | Lista com filtros (filtra `deleted_at IS NULL`) |
| `Update(ctx, *Task)` | Atualiza tarefa |
| `Delete(ctx, id)` | Soft delete |
| `FindByStatus(ctx, status)` | Lista por status (NÃO EXPOSTO como endpoint) |

---

## Códigos de Erro — Catálogo Completo

### BadRequestCode (4xx)

| Código | Constante | Mensagem | HTTP Status |
|---|---|---|---|
| 400001 | `TaskTitleEmpty` | "task title cannot be empty" | 400 |
| 400002 | `TaskAlreadyDone` | "task is already completed" | 400 |
| 400003 | `TaskInitDone` | "task cannot start with a completed state" | 400 |
| 400004 | `InvalidPriority` | "priority must be between 1 and 5" | 400 |
| 400005 | `UnprocessableEntity` | "JSON invalido" | 422 |
| 400006 | `IdInvalid` | "ID invalido" | 400 |
| 400007 | `TaskNotFound` | "Task nao existe" | 404 |

### UnexpectedCode (5xx)

| Código | Constante | Mensagem | HTTP Status |
|---|---|---|---|
| 500001 | `UnexpectedError` | "Unexpected Error" | 500 |

> **Nota**: HTTP Status são raw ints (400, 404, 422, 500). O pacote `exceptions/codes/` NÃO importa `net/http`.

---

## Gaps e Endpoints Ausentes

| Método | Endpoint Sugerido | Descrição |
|---|---|---|
| `FindByStatus(ctx, status)` | `GET /api/task?status=PENDING` | Listar por status (já coberto por `FindAll` com filtro) |

O método **`Complete()`** existe no modelo mas **não é invocado** por nenhum endpoint atualmente.

---

## Observações

1. **Update usa PUT em `/api/task`** com ID no body — convencionalmente seria `/api/task/{id}`
2. **Update preserva `createdAt`** — busca a task existente antes de atualizar (corrigido)
3. **`ExistByID` filtra `deleted_at IS NULL`** — tasks deletadas não são consideradas existentes (corrigido)
4. **Validação de domínio ocorre no UseCase** — `models.NewTask()` é chamado dentro do UseCase, não no handler
5. **Handlers convertem Request → Command** — nunca criam modelos de domínio diretamente
