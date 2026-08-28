package controllers

import (
	"bufio"
	"encoding/json"
	"time"

	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type presenceSnapshot struct {
	OnlineUserIDs []string `json:"online_user_ids"`
}

type presenceUpdate struct {
	UserID string `json:"user_id"`
	Online bool   `json:"online"`
}

// opens one authenticated server-sent event stream for messaging updates.
func (c *MessagingController) StreamRealtime(ctx *fiber.Ctx) error {
	if c.realtimeHub == nil || c.messagingRepo == nil {
		return pipeError(ctx, fiber.StatusInternalServerError, "internal_error")
	}
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	sub, becameOnline := c.realtimeHub.Subscribe(userID)
	relatedUserIDs, _ := c.messagingRepo.FindRelatedConversationUserIDs(ctx.Context(), userID)
	if becameOnline {
		c.publishPresence(relatedUserIDs, userID.String(), true)
	}

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("X-Accel-Buffering", "no")
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer c.cleanupRealtimeSubscription(sub, relatedUserIDs)
		c.writePresenceSnapshot(w, relatedUserIDs)
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case payload, ok := <-sub.Events:
				if !ok || realtime.WriteSSE(w, payload) != nil {
					return
				}
			case <-ticker.C:
				if realtime.WriteSSEComment(w, "ping") != nil {
					return
				}
			}
		}
	})
	return nil
}

// removes one active realtime subscription and publishes offline presence if needed.
func (c *MessagingController) cleanupRealtimeSubscription(sub *realtime.Subscription, relatedUserIDs []uuid.UUID) {
	if c.realtimeHub.Unsubscribe(sub) {
		c.publishPresence(relatedUserIDs, sub.UserID.String(), false)
	}
}

// sends the current online peers to the connected user.
func (c *MessagingController) writePresenceSnapshot(w *bufio.Writer, relatedUserIDs []uuid.UUID) {
	online := c.realtimeHub.OnlineUsers(relatedUserIDs)
	onlineIDs := make([]string, 0, len(online))
	for _, userID := range online {
		onlineIDs = append(onlineIDs, userID.String())
	}
	payload, err := json.Marshal(realtime.Event{Type: "presence.snapshot", Data: presenceSnapshot{OnlineUserIDs: onlineIDs}})
	if err != nil {
		return
	}
	_ = realtime.WriteSSE(w, payload)
}

// broadcasts one presence update to the supplied users.
func (c *MessagingController) publishPresence(userIDs []uuid.UUID, actorUserID string, online bool) {
	if len(userIDs) == 0 {
		return
	}
	_ = c.realtimeHub.PublishUsers(userIDs, realtime.Event{
		Type: "presence.updated",
		Data: presenceUpdate{UserID: actorUserID, Online: online},
	})
}
