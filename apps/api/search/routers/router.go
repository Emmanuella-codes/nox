package routers

import (
	"github.com/emmanuella-codes/nox/search/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
)

func SearchRoutes(controller *controllers.SearchController) []api.RouterSchema {
	return []api.RouterSchema{
		{RouteMethod: api.RouteMethod("GET"), Path: "/", Handler: controller.Search},
	}
}
