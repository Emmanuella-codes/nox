package routers

import (
	"github.com/emmanuella-codes/nox/hashtag/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
)

func HashtagRoutes(controller *controllers.HashtagController) []api.RouterSchema {
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("GET"), Path: "/trending", Handler: controller.Trending},
		{RouteMethod: api.RouteMethod("GET"), Path: "/:tag", Handler: controller.GetHashtag},
		{RouteMethod: api.RouteMethod("GET"), Path: "/:tag/posts", Handler: controller.PostsByTag},
	}
}
