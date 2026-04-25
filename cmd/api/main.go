package main

import (
	"github.com/MarcusVNJ/GOTODO/cmd/api/di"
	taskHandlers "github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/handlers/task"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/server"
	repository_impl "github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/repository"
	"github.com/MarcusVNJ/GOTODO/internal/config"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		repository_impl.Module,
		di.UsecaseModule,
		taskHandlers.Module,
		server.Module,
	).Run()
}
