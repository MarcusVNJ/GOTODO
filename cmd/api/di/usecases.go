package di

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	core_usecase "github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	task_usecase "github.com/MarcusVNJ/GOTODO/internal/core/usecase/task"
	"go.uber.org/fx"
)

var UsecaseModule = fx.Module("usecases",
	fx.Provide(
		fx.Annotate(
			task_usecase.NewCreateTaskUC,
			fx.As(new(core_usecase.IUsecase[*models.Task, struct{}])),
			fx.ResultTags(`name:"createTaskUC"`),
		),
		fx.Annotate(
			task_usecase.NewDeleteTaskUC,
			fx.As(new(core_usecase.IUsecase[string, struct{}])),
			fx.ResultTags(`name:"deleteTaskUC"`),
		),
		fx.Annotate(
			task_usecase.NewGetTaskByIdUC,
			fx.As(new(core_usecase.IUsecase[string, *models.Task])),
			fx.ResultTags(`name:"getTaskByIdUC"`),
		),
		fx.Annotate(
			task_usecase.NewUpdateTaskUC,
			fx.As(new(core_usecase.IUsecase[*models.Task, struct{}])),
			fx.ResultTags(`name:"updateTaskUC"`),
		),
	),
)
