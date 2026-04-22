package routers

import (
	"net/http"

	handlers "github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/handlers/task"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/middlewares"
	repository_impl "github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/repository"
	usecase "github.com/MarcusVNJ/GOTODO/internal/core/usecase/task"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MakeTaskRoutes(api huma.API, db *pgxpool.Pool) http.Handler {

	// Repositories
	taskRepository := repository_impl.NewPostgresTaskRepository(db)

	// UseCases
	saveTaskUsecase := usecase.NewCreateTaskUC(taskRepository)
	deleteTaskUsecase := usecase.NewDeleteTaskUC(taskRepository)
	getTaskUsecase := usecase.NewGetTaskByIdUC(taskRepository)
	updateTaskUsecase := usecase.NewUpdateTaskUC(taskRepository)

	// Resources
	resourceTaskSave := handlers.NewCreateTaskResource(saveTaskUsecase)
	resourceTaskDelete := handlers.NewDeleteTaskResource(deleteTaskUsecase)
	resourceGetTask := handlers.NewGetTaskByIdResource(getTaskUsecase)
	resourceUpdateTask := handlers.NewUpdateTaskResource(updateTaskUsecase)

	huma.Register(api, huma.Operation{
		OperationID:   "create-task",
		Method:        http.MethodPost,
		Path:          "/api/task",
		Summary:       "Criar Tarefa",
		Description:   "Cria uma nova tarefa",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusCreated,
	}, middlewares.HandlerException(resourceTaskSave.Handler))

	huma.Register(api, huma.Operation{
		OperationID: "get-task-by-id",
		Method:      http.MethodGet,
		Path:        "/api/task/{id}",
		Summary:     "Buscar Tarefa",
		Description: "Busca uma tarefa por ID",
		Tags:        []string{"Tasks"},
	}, middlewares.HandlerException(resourceGetTask.Handler))

	huma.Register(api, huma.Operation{
		OperationID: "update-task",
		Method:      http.MethodPut,
		Path:        "/api/task",
		Summary:     "Atualizar Tarefa",
		Description: "Atualiza os atributos de uma tarefa",
		Tags:        []string{"Tasks"},
	}, middlewares.HandlerException(resourceUpdateTask.Handler))

	huma.Register(api, huma.Operation{
		OperationID:   "delete-task",
		Method:        http.MethodDelete,
		Path:          "/api/task/{id}",
		Summary:       "Excluir Tarefa",
		Description:   "Remove uma tarefa por ID",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusNoContent,
	}, middlewares.HandlerException(resourceTaskDelete.Handler))

	return chi.NewRouter()
}
