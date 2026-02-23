package routers

import (
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type AppRouter struct {
	chi.Router
}

func (chi *AppRouter) Get(pattern string, handler middlewares.ResourceHandler) {
	chi.Router.Get(pattern, middlewares.ExceptionHandler(handler))
}

func (chi *AppRouter) Post(pattern string, handler middlewares.ResourceHandler)  {
	chi.Router.Post(pattern, middlewares.ExceptionHandler(handler))
}