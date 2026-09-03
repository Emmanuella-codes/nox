package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/notification/controllers"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

func NotificationRoutes(controller *controllers.NotificationController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/notification-devices",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.ListNotificationDevices,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/notification-devices",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.UpsertNotificationDevice,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/notification-devices/:deviceID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.RemoveNotificationDevice,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/notification-preferences",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.ListNotificationPreferences,
		},
		{
			RouteMethod: api.RouteMethod("PUT"),
			Path:        "/notification-preferences",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.UpdateNotificationPreference,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/notifications",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.ListNotifications,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/notifications/unread-count",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.GetUnreadCount,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/notifications/read-all",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.MarkAllNotificationsRead,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/notifications/:notificationID/read",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.MarkNotificationRead,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/notifications/realtime/stream",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.StreamNotifications,
		},
	}
}
