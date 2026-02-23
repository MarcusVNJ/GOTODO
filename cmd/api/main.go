package main

import (
	"context"
	"github.com/MarcusVNJ/GOTODO/cmd/api/router"
	"github.com/MarcusVNJ/GOTODO/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
)



func main() {

	cfg, err := config.LoadConfig()
	if err!= nil {
		slog.Error("Erro ao carregar a configuração da aplicação", slog.Any("err", err))
		os.Exit(1)
	}

	config.InitLogger(cfg.Environment)

	slog.Info("Iniciando aplicação", slog.String("env", cfg.Environment))

	ctx := context.Background()

	db, err := config.InitDB(ctx, cfg.DatabaseURL)
	if err!= nil {
		slog.Error("Falha fatal no banco de dados", slog.Any("err", err))
		os.Exit(1)
	}
	defer db.Close()


	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Mount("/api", routers.MakeTaskRoutes(db))


	addr := ":" + cfg.Port
	slog.Info("Iniciando servidor", slog.String("porta", cfg.Port))
	if err := http.ListenAndServe(addr, r); err!= nil {
		slog.Error("Servidor parou", slog.Any("err", err))
		os.Exit(1)
	}
}