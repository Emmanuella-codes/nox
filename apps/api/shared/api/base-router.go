package api

import (
	"github.com/emmanuella-codes/nox/typings"
	"github.com/gofiber/fiber/v2"
)

func mergeMiddlewares(fMiddlewares []typings.FiberMiddleware, handler typings.FiberMiddleware) []func(*fiber.Ctx) error {
	// Convert the slice of FiberMiddleware to a slice of func(*fiber.Ctx) error
	var middlewares []func(*fiber.Ctx) error
	for _, mw := range fMiddlewares {
		middlewares = append(middlewares, mw)
	}

	// Append the handler function to the middlewares slice
	middlewares = append(middlewares, handler)
	return middlewares
}

func BaseRouter(router fiber.Router, schema []RouterSchema) {
	for _, route := range schema {
		router.Add(
			string(route.RouteMethod),
			route.Path,
			mergeMiddlewares(route.Middlewares, route.Handler)...,
		)
	}
}
