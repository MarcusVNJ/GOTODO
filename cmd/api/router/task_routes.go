package routers

import (
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/handlers"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/repository"
    "github.com/MarcusVNJ/GOTODO/internal/core/usecase"
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "net/http"
)

func MakeTaskRoutes(db *pgxpool.Pool) http.Handler {

	//Repositories
	taskRepository := repository_impl.NewPostgresTaskRepository(db)

	//UseCases
	saveTaskUsecase := usecase.NewCreateTaskUC(taskRepository)

	//Resources
	resourceTaskSave := handlers.NewCreateTaskResource(saveTaskUsecase)


    router := chi.NewRouter()

    app := AppRouter{router}

    app.Post("/task", resourceTaskSave.Handler)

	return router
}
