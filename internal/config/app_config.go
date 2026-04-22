package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	routers "github.com/MarcusVNJ/GOTODO/cmd/api/router"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	Port        string `envconfig:"PORT" default:"8080"`
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
	EnableDocs  bool   `envconfig:"ENABLE_DOCS" default:"true"`
}

func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	var cfg AppConfig

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("falha ao processar variáveis de ambiente: %w", err)
	}

	return &cfg, nil
}

func InitLogger(env string) {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func InitDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("não foi possível configurar o pool de conexões: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Ping(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("banco de dados inacessível: %w", err)
	}

	slog.Info("Conexão com o banco de dados estabelecida com sucesso (pgx pool)")
	return db, nil
}

func AddRouters(router *chi.Mux, api huma.API, db *pgxpool.Pool) {
	router.Mount("/api", routers.MakeTaskRoutes(api, db))
}

func AddMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
}

func AddExternalDocs(router *chi.Mux, enableDocs bool) huma.API {
	humaConfig := huma.DefaultConfig("GOTODO API", "1.0.0")
	if enableDocs {
		humaConfig.DocsPath = "/docs"
		humaConfig.OpenAPIPath = "/openapi.json"
		humaConfig.SchemasPath = "/schemas"
	} else {
		humaConfig.DocsPath = ""
		humaConfig.OpenAPIPath = ""
		humaConfig.SchemasPath = ""
	}

	return humachi.New(router, humaConfig)
}
