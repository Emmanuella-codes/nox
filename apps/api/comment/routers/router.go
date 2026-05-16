package routers

import (
	"github.com/emmanuella-codes/nox/comment/controllers"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func CommentRoutes(controller *controllers.CommentController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/posts/:postID/comments",
			Handler:     controller.ListComments,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/posts/:postID/comments",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.CreateComment,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/comments/:commentID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.DeleteComment,
		},
	}
}
