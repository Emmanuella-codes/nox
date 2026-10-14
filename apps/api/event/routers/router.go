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
		{RouteMethod: api.RouteMethod("GET"), Path: "/", Middlewares: []typings.FiberMiddleware{middleware.OptionalJWT(cfg)}, Handler: controller.ListEvents},
		{RouteMethod: api.RouteMethod("GET"), Path: "/:eventID", Middlewares: []typings.FiberMiddleware{middleware.OptionalJWT(cfg)}, Handler: controller.GetEvent},
		{RouteMethod: api.RouteMethod("POST"), Path: "/", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.CreateEvent},
		{RouteMethod: api.RouteMethod("PATCH"), Path: "/:eventID", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.UpdateEvent},
		{RouteMethod: api.RouteMethod("POST"), Path: "/:eventID/publish", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.PublishEvent},
		{RouteMethod: api.RouteMethod("POST"), Path: "/:eventID/close", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.CloseEvent},
		{RouteMethod: api.RouteMethod("POST"), Path: "/:eventID/archive", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.ArchiveEvent},
	}
}
