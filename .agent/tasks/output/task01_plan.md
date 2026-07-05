# [FEAT]: Listar tasks e estender busca por ID

---

## 1. Contexto e Objetivo

**História de Usuário Original**:
> Como usuário, eu quero poder listar todas as tasks no meu sistema, apresentando apenas título, prioridade e id no retorno. Como usuário, quero também poder fazer uma busca detalhada da minha task com o id, trazendo assim todas as informações necessárias da task.

**Problema**: Não existe endpoint de listagem — o método `FindAll` do repository está implementado mas não exposto na API. O endpoint de busca por ID (`GET /api/task/{id}`) existe mas retorna apenas `title, description, status, priority`, sem `id, created_at, updated_at`.

## 2. Regras de Negócio

- **RN01**: A listagem retorna apenas `id`, `title` e `priority` — campos como `description`, `status`, `created_at`, `updated_at` são omitidos *(inferido)*
- **RN02**: A listagem suporta filtros opcionais por `status` e `priority` mínima via query params *(inferido — método `FindAll(statusFilter, minPriority)` já existe)*
- **RN03**: A listagem exclui tasks deletadas (soft delete — `deleted_at IS NULL`) *(inferido — gap atual no `FindAll`)*
- **RN04**: A busca detalhada por ID retorna todos os campos: `id`, `title`, `description`, `status`, `priority`, `created_at`, `updated_at`
- **RN05**: A listagem é ordenada por `created_at DESC` *(já implementado no QueryBuilder)*

## 3. Especificações Técnicas

### 3.1 Abordagem Sugerida

**Core**:
- UseCase: `GetAllTasksUC` — recebe `TaskFilter{Status string, MinPriority int}`, retorna `[]*models.Task`, chama `FindAll` do repository
- Tipo `TaskFilter` — struct de input para filtros opcionais

**Outbound**:
- QueryBuilder: `QueryFindAllTasks` — adicionar `WHERE deleted_at IS NULL` como filtro fixo (hoje não filtra soft delete)

**Inbound**:
- Request DTO: `ListTaskRequest` — query params `status` e `min_priority`
- Response DTO: `taskListItem` — `{id, title, priority}`, `NewListTaskResponse(tasks []*models.Task) []taskListItem`
- Response DTO: **estender** `taskResponse` existente para incluir `id`, `created_at`, `updated_at`
- Handler: `ListTaskResource` — `GET /api/task`, operationID `list-tasks`

**Orchestrator**:
- Wire `GetAllTasksUC` em `cmd/api/di/usecases.go` com `ResultTags("name:getAllTasksUC")`
- Registrar `NewListTaskResource` em `handlers/module.go` com `ParamTags("name:getAllTasksUC")`

### 3.2 Dependências

- Não depende de outras tasks
- Usa `FindAll` já existente no repository

### 3.3 Restrições

- Nenhuma migration necessária — tabela `tasks` já tem todos os campos
- Nenhum novo error code necessário — `TaskNotFound` já existe para busca por ID

## 4. Endpoints

| Método | Rota | Operation ID | Status | Descrição | Tipo |
|---|---|---|---|---|---|
| GET | `/api/task` | `list-tasks` | 200 | Listar tarefas (resposta simplificada) | NOVO |
| GET | `/api/task/{id}` | `get-task-by-id` | 200 | Buscar tarefa por ID (resposta completa) | ALTERADO |

### DTOs de Request

```go
type ListTaskRequest struct {
    Status      string `query:"status" description:"Filtrar por status (PENDING, IN_PROCESS, COMPLETED, CANCELLED)"`
    MinPriority int    `query:"min_priority" description:"Filtrar por prioridade mínima (1-5)"`
}
```

### DTOs de Response

```go
// Novo — resposta simplificada da listagem
type taskListItem struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Priority int    `json:"priority"`
}

type ListTaskResponse struct {
    Body []taskListItem
}

// Alterado — estender taskResponse existente para incluir id, created_at, updated_at
type taskResponse struct {
    ID          string       `json:"id"`           // NOVO
    Title       string       `json:"title"`
    Description string       `json:"description,omitempty"`
    Status      enums.Status `json:"status"`
    Priority    int          `json:"priority"`
    CreatedAt   string       `json:"created_at"`   // NOVO
    UpdatedAt   string       `json:"updated_at"`   // NOVO
}
```

## 7. Critérios de Aceite

- [ ] **AC01**: `GET /api/task` retorna lista de tasks ativas com campos `id`, `title`, `priority`
- [ ] **AC02**: `GET /api/task` retorna array vazio `[]` quando não há tasks (não retorna erro)
- [ ] **AC03**: `GET /api/task?status=PENDING` filtra por status
- [ ] **AC04**: `GET /api/task?min_priority=3` filtra por prioridade mínima (>= 3)
- [ ] **AC05**: `GET /api/task` não retorna tasks deletadas (soft delete)
- [ ] **AC06**: `GET /api/task` retorna tasks ordenadas por `created_at DESC`
- [ ] **AC07**: `GET /api/task/{id}` retorna `id`, `title`, `description`, `status`, `priority`, `created_at`, `updated_at`
- [ ] **AC08**: `GET /api/task/{id}` com ID inexistente retorna BusinessException(TaskNotFound) com status 404

## 9. Ordem de Implementação

1. **Core**: Criar `TaskFilter` e `GetAllTasksUC`
2. **Outbound**: Ajustar `QueryFindAllTasks` para adicionar `WHERE deleted_at IS NULL`
3. **Inbound**: Criar `ListTaskRequest`, `taskListItem`/`ListTaskResponse`, `ListTaskResource`. Estender `taskResponse` com `id`, `created_at`, `updated_at`
4. **Orchestrator**: Wire `GetAllTasksUC` em `di/usecases.go` e registrar `ListTaskResource` em `handlers/module.go`

## 10. Referências

- [ ] História de usuário original sobre listagem e busca detalhada de tasks
- [ ] [Recursos da API](../knowledge-base/api-resources.md) — gap: `FindAll` sem endpoint, `GetByID` sem campos de auditoria
- [ ] [Arquitetura](../knowledge-base/architecture.md), [DI Patterns](../knowledge-base/di-patterns.md)