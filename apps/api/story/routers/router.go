package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/story/controllers"
	"github.com/emmanuella-codes/nox/typings"
)

func StoryRoutes(controller *controllers.StoryController, cfg *config.Config) []api.RouterSchema {
	auth := []typings.FiberMiddleware{middleware.JWT(cfg)}
	optionalAuth := []typings.FiberMiddleware{middleware.OptionalJWT(cfg)}
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("GET"), Path: "/stories/:storyID", Middlewares: optionalAuth, Handler: controller.GetStory},
		{RouteMethod: api.RouteMethod("GET"), Path: "/stories/:storyID/items", Middlewares: optionalAuth, Handler: controller.ListStoryItems},
		{RouteMethod: api.RouteMethod("GET"), Path: "/events/:eventID/stories", Middlewares: optionalAuth, Handler: controller.ListEventStories},
		{RouteMethod: api.RouteMethod("GET"), Path: "/events/:eventID/highlight-stories", Middlewares: optionalAuth, Handler: controller.ListEventHighlightStories},
		{RouteMethod: api.RouteMethod("GET"), Path: "/personas/:personaID/stories", Middlewares: optionalAuth, Handler: controller.ListPersonaStories},
		{RouteMethod: api.RouteMethod("POST"), Path: "/stories", Middlewares: auth, Handler: controller.CreateStory},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/stories/:storyID", Middlewares: auth, Handler: controller.DeleteStory},
		{RouteMethod: api.RouteMethod("POST"), Path: "/stories/:storyID/items", Middlewares: auth, Handler: controller.AddStoryItem},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/stories/:storyID/items/:itemID", Middlewares: auth, Handler: controller.DeleteStoryItem},
		{RouteMethod: api.RouteMethod("PATCH"), Path: "/stories/:storyID/items/:itemID/position", Middlewares: auth, Handler: controller.ReorderStoryItem},
		{RouteMethod: api.RouteMethod("POST"), Path: "/events/:eventID/highlight-stories", Middlewares: auth, Handler: controller.AddEventHighlightStory},
		{RouteMethod: api.RouteMethod("DELETE"), Path: "/events/:eventID/highlight-stories/:storyID", Middlewares: auth, Handler: controller.RemoveEventHighlightStory},
		{RouteMethod: api.RouteMethod("PATCH"), Path: "/events/:eventID/highlight-stories/:storyID/position", Middlewares: auth, Handler: controller.ReorderEventHighlightStory},
	}
}
