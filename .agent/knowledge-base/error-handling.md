# Tratamento de Erros

## Arquitetura de Erros

O projeto GOTODO utiliza um sistema de erros tipados com duas famílias principais:

```
Exception (interface)
├── BusinessException → erros de negócio (4xx)
└── UnexpectedException → erros técnicos (5xx)
```

## Interface Exception

Definida em `internal/core/exceptions/error_code.go`:

```go
type Exception interface {
    Code() int
    Message() string
    HTTPStatus() int
}
```

**IMPORTANTE**: `HTTPStatus()` retorna raw ints (400, 404, 422, 500), **NUNCA** `http.Status*` do pacote `net/http`. O Core não importa `net/http`.

## Códigos de Erro

### BadRequestCode (4xx)

Localização: `internal/core/exceptions/codes/badrequest_code.go`

```go
type BadRequestCode int

const (
    TaskTitleEmpty      BadRequestCode = 400001
    TaskAlreadyDone     BadRequestCode = 400002
    TaskInitDone        BadRequestCode = 400003
    InvalidPriority     BadRequestCode = 400004
    UnprocessableEntity BadRequestCode = 400005
    IdInvalid           BadRequestCode = 400006
    TaskNotFound        BadRequestCode = 400007
)

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

Localização: `internal/core/exceptions/codes/unexpected_code.go`

```go
type UnexpectedCode int

const (
    UnexpectedError UnexpectedCode = 500001
)

func (code UnexpectedCode) HTTPStatus() int {
    return 500  // raw int, NÃO http.StatusInternalServerError
}
```

## BusinessException

Erro de negócio esperado — o cliente recebe detalhes do que ocorreu:

```go
type BusinessException struct {
    Code       int    `json:"code"`
    Message    string `json:"message"`
    HTTPStatus int    `json:"-"`
}

func NewBusinessException(err Exception) *BusinessException {
    return &BusinessException{
        Code:       err.Code(),
        Message:    err.Message(),
        HTTPStatus: err.HTTPStatus(),
    }
}

func (err *BusinessException) Error() string { ... }
func (err *BusinessException) GetStatus() int { return err.HTTPStatus }
```

## UnexpectedException

Erro técnico inesperado — o cliente recebe mensagem genérica, detalhes ficam no log:

```go
type UnexpectedException struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details any    `json:"-"`
}

func NewUnexpectedException(err Exception, details error) *UnexpectedException {
    code := err.HTTPStatus()
    if code == 0 {
        code = 500  // raw int, NÃO http.StatusInternalServerError
    }
    return &UnexpectedException{
        Code:    code,
        Message: err.Message(),
        Details: details,
    }
}

func (e *UnexpectedException) GetStatus() int { return e.Code }
```

## Fluxo de Tratamento

### No Handler (Adapter Inbound)

Handlers NUNCA definem status code HTTP. Apenas retornam `nil, err`:

```go
func (r *CreateTaskResource) Handler(ctx context.Context, input *request.CreateTaskRequest) (*response.OperationTaskResponse, error) {
    cmd := app_task.CreateTaskCommand{...}
    result, err := r.usecase.Execute(ctx, cmd)
    if err != nil {
        return nil, err  // ✅ retorna nil, err - sem status code
    }
    return &response.OperationTaskResponse{...}, nil
}
```

### No UseCase (Application Layer)

UseCases criam `BusinessException` para regras de negócio violadas:

```go
func (usecase *UpdateTaskUC) UpdateTask(ctx context.Context, cmd UpdateTaskCommand) (struct{}, error) {
    exist, err := usecase.repository.ExistByID(ctx, cmd.ID)
    if !exist {
        return struct{}{}, exceptions.NewBusinessException(codes.TaskNotFound)
    }
    ...
}
```

### No Repositório (Adapter Outbound)

Erros de infraestrutura são envelopados com `oops` para stack trace rica. O `HandlerException` middleware transforma `oops.OopsError` em `UnexpectedException`:

```go
func (repository *PostgresTaskRepository) Save(ctx context.Context, task *models.Task) error {
    _, err = repository.db.Exec(ctx, query, args...)
    if err != nil {
        return oops.
            In("PostgresTaskRepository").
            Tags("database", "postgres").
            Wrapf(err, "falha crítica ao tentar inserir task no banco")
    }
    return nil
}
```

### No Middleware (HandlerException)

O middleware `HandlerException` intercepta todos os erros e classifica:

```go
func handlerError(err error) error {
    // 1. BusinessException → retorna como está (4xx)
    var businessErr *exceptions.BusinessException
    if errors.As(err, &businessErr) {
        return businessErr
    }

    // 2. oops.OopsError → UnexpectedException (5xx)
    var oopsErr *oops.OopsError
    if errors.As(err, &oopsErr) {
        return exceptions.NewUnexpectedException(codes.UnexpectedError, oopsErr)
    }

    // 3. Erro genérico → UnexpectedException (5xx)
    unexpectedErr := exceptions.NewUnexpectedException(codes.UnexpectedError, oops.Wrap(err))
    unexpectedErr.Message = "Erro interno crítico no servidor."
    return unexpectedErr
}
```

## Logging

- **BusinessException** → `slog.Info` (não é erro de sistema, é fluxo normal)
- **UnexpectedException** → `slog.Error` com stack trace `oops`
- Logging é assíncrono (`go func()` separado) — não bloqueia a resposta HTTP

## Regras

| # | Regra | Motivo |
|---|---|---|
| 1 | Core NUNCA importa `net/http` | Pureza do domínio |
| 2 | HTTPStatus() retorna raw ints (400, 404, 422, 500) | Independência de protocolo |
| 3 | Handlers retornam `nil, err` — nunca definem status code | Separação de responsabilidades |
| 4 | UseCases criam `BusinessException` para regras de negócio | Erro tipado e rastreável |
| 5 | Repositórios envelopam erros com `oops` | Stack trace rica para debugging |
| 6 | Logging de erros 5xx é assíncrono | Não bloqueia a resposta |
| 7 | Erros 4xx são logados como INFO | Não poluem alertas |
