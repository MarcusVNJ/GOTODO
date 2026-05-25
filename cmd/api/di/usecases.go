package di

import (
	taskdto "github.com/MarcusVNJ/GOTODO/internal/app/task/dto"
	app_task "github.com/MarcusVNJ/GOTODO/internal/app/task"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"go.uber.org/fx"
)

var UsecaseModule = fx.Module("usecases",
	fx.Provide(
		fx.Annotate(
			app_task.NewCreateTaskUC,
			fx.As(new(usecase.IUsecase[taskdto.CreateTaskCommand, taskdto.CreateTaskResult])),
			fx.ResultTags(`name:"createTaskUC"`),
		),
		fx.Annotate(
			app_task.NewDeleteTaskUC,
			fx.As(new(usecase.IUsecase[string, struct{}])),
			fx.ResultTags(`name:"deleteTaskUC"`),
		),
		fx.Annotate(
			app_task.NewGetTaskByIdUC,
			fx.As(new(usecase.IUsecase[taskdto.GetTaskQuery, *models.Task])),
			fx.ResultTags(`name:"getTaskByIdUC"`),
		),
		fx.Annotate(
			app_task.NewUpdateTaskUC,
			fx.As(new(usecase.IUsecase[taskdto.UpdateTaskCommand, struct{}])),
			fx.ResultTags(`name:"updateTaskUC"`),
		),
		fx.Annotate(
			app_task.NewGetAllTasksUC,
			fx.As(new(usecase.IUsecase[taskdto.ListTasksQuery, []*models.Task])),
			fx.ResultTags(`name:"getAllTasksUC"`),
		),
	),
)
