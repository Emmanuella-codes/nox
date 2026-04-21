package api

import "github.com/emmanuella-codes/nox/typings"

type RouteMethod string

type RouterSchema struct {
	RouteMethod RouteMethod
	Path        string
	Middlewares []typings.FiberMiddleware
	Handler     typings.FiberMiddleware
}
