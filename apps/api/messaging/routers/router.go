package routers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/messaging/controllers"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/typings"
)

// MessagingRoutes returns the messaging HTTP route definitions.
func MessagingRoutes(controller *controllers.MessagingController, cfg *config.Config) []api.RouterSchema {
	return []api.RouterSchema{
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/conversations/direct",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.CreateDirectConversation,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/conversations/group",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.CreateGroupConversation,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/conversations",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.ListConversations,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/conversations/:conversationID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.GetConversation,
		},
		{
			RouteMethod: api.RouteMethod("GET"),
			Path:        "/conversations/:conversationID/messages",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.ListMessages,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/conversations/:conversationID/messages",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.SendMessage,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/conversations/:conversationID/read",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.MarkRead,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/conversations/:conversationID/members",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.AddMembers,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/conversations/:conversationID/members/:personaID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.RemoveMember,
		},
		{
			RouteMethod: api.RouteMethod("POST"),
			Path:        "/conversations/:conversationID/leave",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.LeaveConversation,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/conversations/:conversationID/members/:personaID/role",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.UpdateMemberRole,
		},
		{
			RouteMethod: api.RouteMethod("PATCH"),
			Path:        "/messages/:messageID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.EditMessage,
		},
		{
			RouteMethod: api.RouteMethod("DELETE"),
			Path:        "/messages/:messageID",
			Middlewares: []typings.FiberMiddleware{middleware.JWT(cfg)},
			Handler:     controller.DeleteMessage,
		},
	}
}
