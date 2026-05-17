package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/like/controllers"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func LikeRoutes(controller *controllers.LikeController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("POST"), Path: "/posts/:postID/likes", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.LikePost},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/posts/:postID/likes", Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)}, Handler: controller.UnlikePost},
	}
}
