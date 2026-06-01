package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/set/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func SetRoutes(controller *controllers.SetController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("GET"), Path: "/", Handler: controller.ListSets},
		{RouteMethod: api.RouteMethod("GET"), Path: "/:setID", Handler: controller.GetSet},
		{RouteMethod: api.RouteMethod("GET"), Path: "/persona/:personaID", Handler: controller.ListPersonaSets},
		{RouteMethod: api.RouteMethod("POST"), Path: "/", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.CreateSet},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/:setID", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.DeleteSet},
	}
}
