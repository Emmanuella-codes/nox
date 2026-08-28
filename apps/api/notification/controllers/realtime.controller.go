package controllers

import (
	"bufio"
	"context"
	"time"

	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/notification/messages"
	notification_pipes "github.com/emmanuella-codes/nox/notification/pipes"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// StreamNotifications opens one authenticated notification event stream.
func (c *NotificationController) StreamNotifications(ctx *fiber.Ctx) error {
	if c.realtimeHub == nil || c.pipe == nil {
		return pipeError(ctx, fiber.StatusInternalServerError, messages.Internal_Error)
	}
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	sub, _ := c.realtimeHub.Subscribe(userID)
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("X-Accel-Buffering", "no")
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer c.realtimeHub.Unsubscribe(sub)
		c.writeNotificationBootstrap(ctx.Context(), userID, personaID, w)
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case event, ok := <-sub.Events:
				if !ok || realtime.WriteEvent(w, event) != nil {
					return
				}
			case <-ticker.C:
				if realtime.WriteComment(w, "ping") != nil {
					return
				}
			}
		}
	})
	return nil
}

// sends the initial unread count snapshot for the stream.
func (c *NotificationController) writeNotificationBootstrap(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, w *bufio.Writer) {
	res := c.pipe.GetUnreadCountPipe(ctx, userID, personaID)
	if !res.Success || res.Data == nil {
		return
	}
	_ = realtime.WriteEvent(w, realtime.Event{
		Type: "notification.unread_count",
		Data: notification_pipes.NotificationUnreadCountResponse{UnreadCount: res.Data.UnreadCount},
	})
}
