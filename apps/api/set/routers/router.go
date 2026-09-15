package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/set/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func SetRoutes(controller *controllers.SetController, cfg *config.Config) []api.RouterSchema {
	auth := []typings.FiberMiddleware{middleware.JWT(cfg)}
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("GET"), Path: "/", Handler: controller.ListSets},
		{RouteMethod: api.RouteMethod("GET"), Path: "/:setID", Handler: controller.GetSet},
		{RouteMethod: api.RouteMethod("GET"), Path: "/persona/:personaID", Handler: controller.ListPersonaSets},
		{RouteMethod: api.RouteMethod("GET"), Path: "/:setID/comments", Handler: controller.ListSetComments},
		{RouteMethod: api.RouteMethod("POST"), Path: "/", Middlewares: auth, Handler: controller.CreateSet},
		{RouteMethod: api.RouteMethod("PATCH"), Path: "/:setID", Middlewares: auth, Handler: controller.UpdateSet},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/:setID", Middlewares: auth, Handler: controller.DeleteSet},
		{RouteMethod: api.RouteMethod("POST"), Path: "/:setID/likes", Middlewares: auth, Handler: controller.LikeSet},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/:setID/likes", Middlewares: auth, Handler: controller.UnlikeSet},
		{RouteMethod: api.RouteMethod("POST"), Path: "/:setID/plays", Middlewares: auth, Handler: controller.RecordSetPlay},
		{RouteMethod: api.RouteMethod("POST"), Path: "/:setID/comments", Middlewares: auth, Handler: controller.CreateSetComment},
		{RouteMethod: api.RouteMethod("PATCH"), Path: "/comments/:commentID", Middlewares: auth, Handler: controller.UpdateSetComment},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/comments/:commentID", Middlewares: auth, Handler: controller.DeleteSetComment},
		{RouteMethod: api.RouteMethod("POST"), Path: "/comments/:commentID/likes", Middlewares: auth, Handler: controller.LikeSetComment},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/comments/:commentID/likes", Middlewares: auth, Handler: controller.UnlikeSetComment},
	}
}
