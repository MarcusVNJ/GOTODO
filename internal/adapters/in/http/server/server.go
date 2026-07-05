package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/router"
	"github.com/MarcusVNJ/GOTODO/internal/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	return r
}

func NewHumaAPI(r *chi.Mux, cfg *config.AppConfig) huma.API {
	humaConfig := huma.DefaultConfig("GOTODO API", "1.0.0")
	if cfg.EnableDocs {
		humaConfig.DocsPath = "/docs"
		humaConfig.OpenAPIPath = "/openapi.json"
		humaConfig.SchemasPath = "/schemas"
	} else {
		humaConfig.DocsPath = ""
		humaConfig.OpenAPIPath = ""
		humaConfig.SchemasPath = ""
	}

	return humachi.New(r, humaConfig)
}

type RouteParams struct {
	fx.In
	Routes []router.RouteRegister `group:"routes"`
}

func RegisterRoutes(api huma.API, params RouteParams) {
	for _, route := range params.Routes {
		route.Register(api)
	}
}

func StartHTTPServer(lc fx.Lifecycle, r *chi.Mux, cfg *config.AppConfig) {
	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Iniciando servidor HTTP", slog.String("porta", cfg.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("Servidor parou com erro", slog.Any("err", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Desligando servidor HTTP")
			return srv.Shutdown(ctx)
		},
	})
}

var Module = fx.Module("server",
	fx.Provide(
		NewRouter,
		NewHumaAPI,
	),
	fx.Invoke(
		RegisterRoutes,
		StartHTTPServer,
	),
)
