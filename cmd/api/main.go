package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/MarcusVNJ/GOTODO/internal/config"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Erro ao carregar a configuração da aplicação", slog.Any("err", err))
		os.Exit(1)
	}

	config.InitLogger(cfg.Environment)

	slog.Info("Iniciando aplicação", slog.String("env", cfg.Environment))

	ctx := context.Background()

	db, err := config.InitDB(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Falha fatal no banco de dados", slog.Any("err", err))
		os.Exit(1)
	}
	defer db.Close()

	router := chi.NewRouter()

	config.AddMiddlewares(router)

	api := config.AddExternalDocs(router, cfg.EnableDocs)

	config.AddRouters(router, api, db)

	addr := ":" + cfg.Port
	slog.Info("Iniciando servidor", slog.String("porta", cfg.Port))
	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("Servidor parou", slog.Any("err", err))
		os.Exit(1)
	}
}
