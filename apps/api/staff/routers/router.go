package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/staff/controllers"
	"github.com/emmanuella-codes/nox/typings"
)

func Routes(c *controllers.Controller, cfg *config.Config) []api.RouterSchema {
	auth := []typings.FiberMiddleware{middleware.JWT(cfg)}
	return []api.RouterSchema{
		{RouteMethod: "POST", Path: "/events/:eventID/staff/invitations", Middlewares: auth, Handler: c.Invite},
		{RouteMethod: "GET", Path: "/events/:eventID/staff", Middlewares: auth, Handler: c.List},
		{RouteMethod: "PUT", Path: "/staff-memberships/:membershipID/roles", Middlewares: auth, Handler: c.Roles},
		{RouteMethod: "POST", Path: "/staff-memberships/:membershipID/accept", Middlewares: auth, Handler: c.Accept},
		{RouteMethod: "POST", Path: "/staff-memberships/:membershipID/revoke", Middlewares: auth, Handler: c.Revoke},
		{RouteMethod: "POST", Path: "/staff-memberships/:membershipID/shifts", Middlewares: auth, Handler: c.Shift},
	}
}
