package router

import (
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

type RouteRegister interface {
	Register(api huma.API)
}

func AsRoute(f any, annotations ...fx.Annotation) any {
	opts := []fx.Annotation{
		fx.As(new(RouteRegister)),
		fx.ResultTags(`group:"routes"`),
	}
	opts = append(opts, annotations...)
	return fx.Annotate(f, opts...)
}
