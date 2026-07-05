# Estratégia de Testes

## Filosofia

Testes focam na **camada de Application** (UseCases) porque é onde as regras de negócio são orquestradas. Como os UseCases dependem de interfaces (Ports), injetamos Mocks no lugar de repositórios reais, garantindo isolamento total e execução em milissegundos.

## Localização dos Testes

Testes de UseCase ficam em `internal/app/<entity>/` com sufixo `_test.go`:

```
internal/app/task/
├── create_task_uc.go
├── create_task_uc_test.go          # Testes do CreateTaskUC
├── update_task_uc.go
├── update_task_uc_test.go          # Testes do UpdateTaskUC
├── delete_task_uc.go
├── delete_task_uc_test.go          # Testes do DeleteTaskUC
├── get_task_by_id_uc.go
├── get_task_by_id_uc_test.go       # Testes do GetTaskByIdUC
├── task_repository_mock_test.go    # Mock do repositório
└── ...
```

## Padrão de Teste

### AAA (Arrange - Act - Assert)

```go
package task_test  // black-box testing

func TestCreateTaskUC_Execute_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockTaskRepository)
    uc := app_task.NewCreateTaskUC(mockRepo)
    cmd := app_task.CreateTaskCommand{
        Title: "Test Task", Description: "Desc", Priority: 1,
    }
    mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(t interface{}) bool {
        return true
    })).Return(nil)

    // Act
    result, err := uc.Execute(context.Background(), cmd)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ID)
    mockRepo.AssertExpectations(t)
}
```

### Mock do Repositório

```go
type MockTaskRepository struct {
    mock.Mock
}

func (m *MockTaskRepository) Save(ctx context.Context, task *models.Task) error {
    args := m.Called(ctx, task)
    return args.Error(0)
}

func (m *MockTaskRepository) ExistByID(ctx context.Context, id string) (bool, error) {
    args := m.Called(ctx, id)
    return args.Bool(0), args.Error(1)
}
// ... demais métodos
```

## Cenários de Teste Obrigatórios por UseCase

### Create

| Cenário | O que testar |
|---|---|
| **Success** | Command válido → retorna Result com ID, sem erro |
| **ValidationError** | Campo obrigatório vazio → `BusinessException` + código específico |
| **InvalidField** | Campo fora do range → `BusinessException` + código específico |
| **RepositoryError** | Repositório falha → erro propagado |

### GetById

| Cenário | O que testar |
|---|---|
| **Success** | ID existe → retorna modelo preenchido |
| **NotFound** | ID não existe → erro `TaskNotFound` |

### Update

| Cenário | O que testar |
|---|---|
| **Success** | Command válido → `ExistByID` true, `FindByID` retorna existente, `Update` sem erro |
| **NotFound** | ID não existe → `BusinessException(TaskNotFound)` |
| **ExistError** | `ExistByID` falha → erro propagado |
| **UpdateError** | `Update` falha → erro propagado |
| **PreservesCreatedAt** | `createdAt` original é preservado após update |

### Delete

| Cenário | O que testar |
|---|---|
| **Success** | ID existe → `ExistByID` true, `Delete` sem erro |
| **NotFound** | ID não existe → `BusinessException(TaskNotFound)` |
| **ExistError** | `ExistByID` falha → erro propagado |
| **DeleteError** | `Delete` falha → erro propagado |

### List/GetAll

| Cenário | O que testar |
|---|---|
| **Success** | Query válida → retorna slice de modelos |
| **EmptyResult** | Sem resultados → slice vazio, sem erro |
| **RepositoryError** | Repositório falha → erro propagado |

## Ferramentas

| Ferramenta | Uso |
|---|---|
| `github.com/stretchr/testify/assert` | Asserções (`NoError`, `Error`, `ErrorIs`, `ErrorAs`, `Equal`) |
| `github.com/stretchr/testify/mock` | Mocks (`Mock.On(...).Return(...)`, `AssertExpectations`) |

## Nomenclatura

Testes seguem o padrão: `Test<UseCase>_<Scenario>`

```
TestCreateTaskUC_Execute_Success
TestCreateTaskUC_Execute_ValidationError
TestCreateTaskUC_Execute_InvalidPriority
TestCreateTaskUC_Execute_RepositoryError
TestUpdateTaskUC_Execute_Success
TestUpdateTaskUC_Execute_NotFound
TestUpdateTaskUC_PreservesCreatedAt
```

## Boas Práticas

1. **Black-box testing**: use `package <entity>_test` (não `package <entity>`)
2. **Um assertion lógico por teste** (exceto mock expectations)
3. **Testes independentes**: cada teste cria seu próprio mock e use case
4. **Use `mock.MatchedBy`** para asserts flexíveis em parâmetros complexos
5. **Verifique `mockRepo.AssertExpectations(t)`** no fim de cada teste
6. **Para erros de negócio**, use `assert.ErrorAs` para verificar o tipo e `assert.Equal` para o código
7. **Para erros de infra**, use `assert.ErrorIs`

## Comandos

```bash
go test ./internal/app/task/...           # testar use cases de task
go test ./internal/...                     # testar tudo
go test -v -run TestCreate ./internal/app/task/...  # teste específico
```
