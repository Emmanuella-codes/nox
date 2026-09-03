package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/preference/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func PreferenceRoutes(controller *controllers.PreferenceController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("POST"), Path: "/blocks", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.BlockUser},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/blocks", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.UnblockUser},
		{RouteMethod: api.RouteMethod("POST"), Path: "/mutes", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.MuteUser},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/mutes", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.UnmuteUser},
		{RouteMethod: api.RouteMethod("POST"), Path: "/discovery-suppressions", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.AddDiscoverySuppression},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/discovery-suppressions", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.RemoveDiscoverySuppression},
	}
}
