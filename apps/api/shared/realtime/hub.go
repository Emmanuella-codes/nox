package realtime

import (
	"sync"

	"github.com/google/uuid"
)

type Event struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type"`
	Data      any    `json:"data"`
	CreatedAt string `json:"created_at,omitempty"`
}

type Subscription struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Events chan Event
}

type Hub struct {
	mu            sync.RWMutex
	userSubs      map[uuid.UUID]map[uuid.UUID]*Subscription
	presenceCount map[uuid.UUID]int
}

// builds the in-memory realtime event hub.
func NewHub() *Hub {
	return &Hub{
		userSubs:      map[uuid.UUID]map[uuid.UUID]*Subscription{},
		presenceCount: map[uuid.UUID]int{},
	}
}

// registers one user subscription and reports whether the user just became online.
func (h *Hub) Subscribe(userID uuid.UUID) (*Subscription, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sub := &Subscription{ID: uuid.New(), UserID: userID, Events: make(chan Event, 32)}
	if h.userSubs[userID] == nil {
		h.userSubs[userID] = map[uuid.UUID]*Subscription{}
	}
	h.userSubs[userID][sub.ID] = sub
	h.presenceCount[userID]++
	return sub, h.presenceCount[userID] == 1
}

// removes one user subscription and reports whether the user just went offline.
func (h *Hub) Unsubscribe(sub *Subscription) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	subs := h.userSubs[sub.UserID]
	if subs == nil {
		return false
	}
	delete(subs, sub.ID)
	close(sub.Events)
	if len(subs) == 0 {
		delete(h.userSubs, sub.UserID)
	}
	if h.presenceCount[sub.UserID] > 0 {
		h.presenceCount[sub.UserID]--
	}
	if h.presenceCount[sub.UserID] <= 0 {
		delete(h.presenceCount, sub.UserID)
		return true
	}
	return false
}

// sends one event to all active subscriptions for the supplied users.
func (h *Hub) PublishUsers(userIDs []uuid.UUID, event Event) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, userID := range userIDs {
		for _, sub := range h.userSubs[userID] {
			select {
			case sub.Events <- event:
			default:
			}
		}
	}
	return nil
}

// returns the users that currently have at least one active subscription.
func (h *Hub) OnlineUsers(userIDs []uuid.UUID) []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	online := make([]uuid.UUID, 0, len(userIDs))
	for _, userID := range userIDs {
		if h.presenceCount[userID] > 0 {
			online = append(online, userID)
		}
	}
	return online
}
