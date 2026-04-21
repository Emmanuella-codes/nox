package routers

import (
	"github.com/emmanuella-codes/nox/auth/controllers"
	"github.com/emmanuella-codes/nox/auth/pipes"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(controller *controllers.AuthController, jwt fiber.Handler) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/register",
			Handler:     controller.Signup,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/login",
			Handler:     controller.Login,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/refresh",
			Handler:     controller.Refresh,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/logout",
			Middlewares: []typings.FiberMiddleware{jwt},
			Handler:     controller.Logout,
		},
	}
}

func RegisterRoutes(router fiber.Router, pipe *pipes.AuthPipe, jwt fiber.Handler) {
	controller := controllers.NewAuthController(pipe)
	api.BaseRouter(router.Group("/auth"), AuthRoutes(controller, jwt))
}
