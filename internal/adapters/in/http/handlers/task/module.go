package handlers

import (
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/router"
	"go.uber.org/fx"
)

var Module = fx.Module("task_handlers",
	fx.Provide(
		router.AsRoute(NewCreateTaskResource, fx.ParamTags(`name:"createTaskUC"`)),
		router.AsRoute(NewDeleteTaskResource, fx.ParamTags(`name:"deleteTaskUC"`)),
		router.AsRoute(NewGetTaskByIdResource, fx.ParamTags(`name:"getTaskByIdUC"`)),
		router.AsRoute(NewUpdateTaskResource, fx.ParamTags(`name:"updateTaskUC"`)),
	),
)
