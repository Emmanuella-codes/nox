package controllers

import (
	"bufio"
	"encoding/json"
	"strconv"
	"time"

	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
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
	afterID, err := parseRealtimeCursor(ctx)
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_realtime_cursor")
	}
	sub, becameOnline := c.realtimeHub.Subscribe(userID)
	relatedUserIDs, _ := c.messagingRepo.FindRelatedConversationUserIDs(ctx.Context(), userID)
	replay, err := c.messagingRepo.FindConversationEventsAfter(ctx.Context(), userID, afterID, 250)
	if err != nil {
		c.cleanupRealtimeSubscription(sub, relatedUserIDs)
		return pipeError(ctx, fiber.StatusInternalServerError, "internal_error")
	}
	if becameOnline {
		c.publishPresence(relatedUserIDs, userID.String(), true)
	}

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("X-Accel-Buffering", "no")
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer c.cleanupRealtimeSubscription(sub, relatedUserIDs)
		lastEventID := afterID
		for _, event := range replay {
			if c.writeReplayEvent(w, event) != nil {
				return
			}
			lastEventID = event.ID
		}
		c.writePresenceSnapshot(w, relatedUserIDs)
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case event, ok := <-sub.Events:
				if !ok {
					return
				}
				if eventCursor(event) <= lastEventID && eventCursor(event) != 0 {
					continue
				}
				if realtime.WriteEvent(w, event) != nil {
					return
				}
				if cursor := eventCursor(event); cursor > lastEventID {
					lastEventID = cursor
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
	_ = realtime.WriteEvent(w, realtime.Event{Type: "presence.snapshot", Data: presenceSnapshot{OnlineUserIDs: onlineIDs}})
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

// writes one stored conversation event to the realtime stream.
func (c *MessagingController) writeReplayEvent(w *bufio.Writer, event *models.ConversationEvent) error {
	return realtime.WriteEvent(w, realtime.Event{
		ID:        strconv.FormatInt(event.ID, 10),
		Type:      event.EventType,
		Data:      json.RawMessage(event.Payload),
		CreatedAt: event.CreatedAt.Format(time.RFC3339Nano),
	})
}

// resolves one reconnect cursor from query or SSE headers.
func parseRealtimeCursor(ctx *fiber.Ctx) (int64, error) {
	value := ctx.Query("after")
	if value == "" {
		value = ctx.Get("Last-Event-ID")
	}
	if value == "" {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

// extracts one numeric event cursor from a live realtime payload.
func eventCursor(event realtime.Event) int64 {
	if event.ID == "" {
		return 0
	}
	cursor, err := strconv.ParseInt(event.ID, 10, 64)
	if err != nil {
		return 0
	}
	return cursor
}
