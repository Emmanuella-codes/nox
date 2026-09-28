package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/guest/controllers"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func Routes(c *controllers.Controller, cfg *config.Config) []api.RouterSchema {
	auth := []typings.FiberMiddleware{middleware.JWT(cfg)}
	optional := []typings.FiberMiddleware{middleware.OptionalJWT(cfg)}
	return []api.RouterSchema{
		{
			RouteMethod: "POST",
			Path:        "/events/:eventID/guest-list",
			Middlewares: auth,
			Handler:     c.Configure,
		},
		{
			RouteMethod: "POST",
			Path:        "/events/:eventID/guest-list/quotas",
			Middlewares: auth,
			Handler:     c.Quota,
		},
		{
			RouteMethod: "POST",
			Path:        "/events/:eventID/guest-list/quotas/:quotaID/guests",
			Middlewares: auth,
			Handler:     c.Add,
		},
		{
			RouteMethod: "POST",
			Path:        "/guest-invitations/claim",
			Middlewares: optional,
			Handler:     c.Claim,
		},
		{
			RouteMethod: "POST",
			Path:        "/guest-invitations/decline",
			Middlewares: optional,
			Handler:     c.Decline,
		},
		{
			RouteMethod: "POST",
			Path:        "/guest-passes/:entryID/attach",
			Middlewares: auth,
			Handler:     c.Attach,
		},
		{
			RouteMethod: "POST",
			Path:        "/guest-passes/:entryID/revoke",
			Middlewares: auth,
			Handler:     c.Revoke,
		},
		{
			RouteMethod: "POST",
			Path:        "/events/:eventID/guest-list/check-in",
			Middlewares: auth,
			Handler:     c.CheckIn,
		},
		{
			RouteMethod: "GET",
			Path:        "/events/:eventID/guest-list/report",
			Middlewares: auth,
			Handler:     c.Report,
		},
	}
}
