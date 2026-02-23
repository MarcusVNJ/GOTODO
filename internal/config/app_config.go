package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

type AppConfig struct {
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	Port        string `envconfig:"PORT" default:"8080"`
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}


func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	var cfg AppConfig

	err := envconfig.Process("", &cfg)
	if err!= nil {
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
	if err!= nil {
		return nil, fmt.Errorf("não foi possível configurar o pool de conexões: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Ping(pingCtx); err!= nil {
		db.Close()
		return nil, fmt.Errorf("banco de dados inacessível: %w", err)
	}

	slog.Info("Conexão com o banco de dados estabelecida com sucesso (pgx pool)")
	return db, nil
}