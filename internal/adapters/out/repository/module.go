package repository_impl

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/ports"
	"go.uber.org/fx"
)

var Module = fx.Module("repositories",
	fx.Provide(
		fx.Annotate(
			NewPostgresTaskRepository,
			fx.As(new(ports.TaskRepository)),
		),
	),
)
