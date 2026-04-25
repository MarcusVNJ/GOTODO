package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
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

func InitLogger(cfg *AppConfig) {
	var handler slog.Handler

	if cfg.Environment == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func InitDB(lc fx.Lifecycle, cfg *AppConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("não foi possível configurar o pool de conexões: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := db.Ping(pingCtx); err != nil {
				db.Close()
				return fmt.Errorf("banco de dados inacessível: %w", err)
			}

			slog.Info("Conexão com o banco de dados estabelecida com sucesso (pgx pool)")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Encerrando pool de conexões com banco de dados")
			db.Close()
			return nil
		},
	})

	return db, nil
}

var Module = fx.Module("config",
	fx.Provide(
		LoadConfig,
		InitDB,
	),
	fx.Invoke(InitLogger),
)
