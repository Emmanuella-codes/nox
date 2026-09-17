package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/ticket/controllers"
	"github.com/emmanuella-codes/nox/typings"
)

func Routes(c *controllers.Controller, cfg *config.Config) []api.RouterSchema {
	auth := []typings.FiberMiddleware{middleware.JWT(cfg)}
	return []api.RouterSchema{
		{RouteMethod: "GET", Path: "/events/:eventID/ticket-products", Handler: c.Products},
		{RouteMethod: "POST", Path: "/events/:eventID/ticket-waves", Middlewares: auth, Handler: c.CreateWave},
		{RouteMethod: "POST", Path: "/ticket-waves/:waveID/release", Middlewares: auth, Handler: c.ReleaseWave},
		{RouteMethod: "POST", Path: "/ticket-waves/:waveID/ticket-tiers", Middlewares: auth, Handler: c.CreateTier},
		{RouteMethod: "POST", Path: "/ticket-tiers/:tierID/reservations", Middlewares: auth, Handler: c.Reserve},
	}
}
