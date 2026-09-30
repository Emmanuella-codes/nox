package routers

import (
	"github.com/emmanuella-codes/nox/admission/controllers"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func Routes(c *controllers.Controller, cfg *config.Config) []api.RouterSchema {
	auth := []typings.FiberMiddleware{middleware.JWT(cfg)}
	return []api.RouterSchema{
		{RouteMethod: "POST", Path: "/events/:eventID/admissions", Middlewares: auth, Handler: c.Create},
		{RouteMethod: "POST", Path: "/events/:eventID/admissions/check-in", Middlewares: auth, Handler: c.QR},
		{RouteMethod: "GET", Path: "/events/:eventID/admissions/search", Middlewares: auth, Handler: c.Search},
		{RouteMethod: "POST", Path: "/events/:eventID/admissions/manual-check-in", Middlewares: auth, Handler: c.Manual},
		{RouteMethod: "POST", Path: "/admissions/:admissionID/override-check-in", Middlewares: auth, Handler: c.Override},
		{RouteMethod: "GET", Path: "/events/:eventID/admissions/reconciliation", Middlewares: auth, Handler: c.Reconcile},
	}
}
