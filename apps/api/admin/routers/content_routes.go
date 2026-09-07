package routers

import (
	"github.com/emmanuella-codes/nox/admin/controllers"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func adminContentRoutes(controller *controllers.AdminController, cfg *config.Config, repo adminrepo.AdminRepository) []api.RouterSchema {
	routes := []api.RouterSchema{{
		RouteMethod: api.RouteMethod("GET"),
		Path:        "/content/:entityType",
		Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, repo)},
		Handler:     controller.ListContent,
	}, {
		RouteMethod: api.RouteMethod("GET"),
		Path:        "/content/:entityType/:entityID",
		Middlewares: []typings.FiberMiddleware{middleware.AdminJWT(cfg, repo)},
		Handler:     controller.GetContent,
	}}
	entities := []struct {
		path string
		kind models.ModerationEntityType
	}{
		{"personas", models.ModerationEntityPersona}, {"posts", models.ModerationEntityPost},
		{"comments", models.ModerationEntityComment}, {"events", models.ModerationEntityEvent},
		{"sets", models.ModerationEntitySet}, {"stories", models.ModerationEntityStory},
		{"story-items", models.ModerationEntityStoryItem},
	}
	for _, entity := range entities {
		middleware := []typings.FiberMiddleware{middleware.AdminJWT(cfg, repo)}
		routes = append(routes, api.RouterSchema{RouteMethod: api.RouteMethod("GET"), Path: "/" + entity.path, Middlewares: middleware, Handler: controller.ListContentType(entity.kind)})
		routes = append(routes, api.RouterSchema{RouteMethod: api.RouteMethod("GET"), Path: "/" + entity.path + "/:entityID", Middlewares: middleware, Handler: controller.GetContentType(entity.kind)})
	}
	return routes
}
