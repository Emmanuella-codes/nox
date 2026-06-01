package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/follow/controllers"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func FollowRoutes(controller *controllers.FollowController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/personas/:personaID/follow",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.FollowPersona,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/personas/:personaID/follow",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.UnfollowPersona,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/personas/:personaID/follow-status",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.GetFollowStatus,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/personas/:personaID/followers",
			Middlewares: []typings.FiberMiddleware{middleware.OptionalJWT(cfg)},
			Handler:     controller.GetFollowers,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/personas/:personaID/following",
			Middlewares: []typings.FiberMiddleware{middleware.OptionalJWT(cfg)},
			Handler:     controller.GetFollowing,
		},
	}
}
