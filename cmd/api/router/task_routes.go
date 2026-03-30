package routers

import (
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/handlers/task"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/repository"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase/task"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func MakeTaskRoutes(db *pgxpool.Pool) http.Handler {

	//Repositories
	taskRepository := repository_impl.NewPostgresTaskRepository(db)

	//UseCases
	saveTaskUsecase := usecase.NewCreateTaskUC(taskRepository)
	deleteTaskUsecase := usecase.NewDeleteTaskUC(taskRepository)
	getTaskUsecase := usecase.NewGetTaskByIdUC(taskRepository)
	updateTaskUsecase := usecase.NewUpdateTaskUC(taskRepository)

	//Resources
	resourceTaskSave := handlers.NewCreateTaskResource(saveTaskUsecase)
	resourceTaskDelete := handlers.NewDeleteTaskResource(deleteTaskUsecase)
	resourceGetTask := handlers.NewGetTaskByIdResource(getTaskUsecase)
	resourceUpdateTask := handlers.NewUpdateTaskResource(updateTaskUsecase)

	router := chi.NewRouter()

	app := AppRouter{router}

	app.Post("/task", resourceTaskSave.Handler)
	app.Delete("/task/{id}", resourceTaskDelete.Handler)
	app.Get("/task/{id}", resourceGetTask.Handler)
	app.Put("/task", resourceUpdateTask.Handler)

	return router
}
