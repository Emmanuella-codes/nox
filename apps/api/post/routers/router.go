package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/post/controllers"
	"github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func PostRoutes(controller *controllers.PostController, cfg *config.Config, personaRepo persona.PersonaRepository) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg), middleware.RequirePersonaOwner(personaRepo)},
			Handler:     controller.CreatePost,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/:postID",
			Handler:     controller.GetPost,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/:postID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.DeletePost,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/persona/:personaID",
			Handler:     controller.GetPersonaPosts,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/persona/:personaID/feed",
			Handler:     controller.GetFeed,
		},
	}
}
