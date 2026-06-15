package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/media/controllers"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func MediaRoutes(controller *controllers.MediaController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/set-video/uploads",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.InitiateSetVideoUpload,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/story-video/uploads",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.InitiateStoryVideoUpload,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/:mediaAssetID/processing/ready",
			Handler:     controller.CompleteMediaProcessing,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/:mediaAssetID/processing/story-ready",
			Handler:     controller.CompleteStoryMediaProcessing,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/:mediaAssetID/processing/failed",
			Handler:     controller.FailMediaProcessing,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/orphans",
			Handler:     controller.CleanupOrphanedMediaAssets,
		},
	}
}
