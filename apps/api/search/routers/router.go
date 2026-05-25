package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/search/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func SearchRoutes(controller *controllers.SearchController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/",
			Middlewares: []typings.FiberMiddleware{middleware.OptionalJWT(cfg)},
			Handler:     controller.Search,
		},
	}
}
