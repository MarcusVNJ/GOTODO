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

func (chi *AppRouter) Post(pattern string, handler middlewares.ResourceHandler) {
	chi.Router.Post(pattern, middlewares.ExceptionHandler(handler))
}

func (chi *AppRouter) Delete(pattern string, handler middlewares.ResourceHandler) {
	chi.Router.Delete(pattern, middlewares.ExceptionHandler(handler))
}

func (chi *AppRouter) Put(pattern string, handler middlewares.ResourceHandler) {
	chi.Router.Put(pattern, middlewares.ExceptionHandler(handler))
}

func (chi *AppRouter) Patch(pattern string, handler middlewares.ResourceHandler) {
	chi.Router.Patch(pattern, middlewares.ExceptionHandler(handler))
}
