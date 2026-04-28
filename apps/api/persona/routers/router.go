package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/persona/controllers"
	"github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func PersonaRoutes(controller *controllers.PersonaController, cfg *config.Config, repo persona.PersonaRepository) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.CreatePersona,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/me",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.GetMyPersonas,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/:personaID",
			Handler:     controller.GetPersona,
		},
		{
			RouteMethod: api.RouteMethod("PUT"),
			Path:        "/:personaID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg), middleware.RequirePersonaOwner(repo)},
			Handler:     controller.UpdatePersona,
		},
	}
}
