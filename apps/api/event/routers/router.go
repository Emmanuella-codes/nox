package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/event/controllers"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func EventRoutes(controller *controllers.EventController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("GET"), Path: "/", Handler: controller.ListEvents},
		{RouteMethod: api.RouteMethod("GET"), Path: "/:eventID", Handler: controller.GetEvent},
		{RouteMethod: api.RouteMethod("POST"), Path: "/", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.CreateEvent},
	}
}
