package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/report/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func ReportRoutes(controller *controllers.ReportController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{{
		RouteMethod: api.RouteMethod("POST"),
		Path:        "/",
		Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
		Handler:     controller.CreateReport,
	}}
}
