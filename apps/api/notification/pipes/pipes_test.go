package pipes

import (
	"context"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/models"
	notification_messages "github.com/emmanuella-codes/nox/notification/messages"
	persona_dtos "github.com/emmanuella-codes/nox/persona/dtos"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/google/uuid"
)

type notificationTestRepo struct {
	notifications []*models.Notification
	markRead      *models.Notification
	markAllCalls  int
	unreadCount   int
}

type notificationPersonaRepo struct {
	personas map[uuid.UUID]*models.Persona
}

// CreateNotifications is unused in these tests.
func (r *notificationTestRepo) CreateNotifications(ctx context.Context, inputs []notification_repo.CreateNotificationInput) ([]*models.Notification, error) {
	return nil, nil
}

// FindPersonaNotifications returns configured notification fixtures.
func (r *notificationTestRepo) FindPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.Notification, error) {
	return r.notifications, nil
}

// CountUnreadPersonaNotifications returns the configured unread count fixture.
func (r *notificationTestRepo) CountUnreadPersonaNotifications(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int, error) {
	return r.unreadCount, nil
}

// MarkNotificationRead returns the configured read notification fixture.
func (r *notificationTestRepo) MarkNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID) (*models.Notification, error) {
	if r.markRead == nil {
		return nil, notification_repo.ErrNotificationNotFound
	}
	return r.markRead, nil
}

// MarkPersonaNotificationsRead records one mark-all call.
func (r *notificationTestRepo) MarkPersonaNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (int64, error) {
	r.markAllCalls++
	return 1, nil
}

// MarkConversationNotificationsRead is unused in these tests.
func (r *notificationTestRepo) MarkConversationNotificationsRead(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) (int64, error) {
	return 0, nil
}

// DeleteMessageNotifications is unused in these tests.
func (r *notificationTestRepo) DeleteMessageNotifications(ctx context.Context, messageID uuid.UUID) error {
	return nil
}

// FindPersonaByID returns the configured persona fixture.
func (r *notificationPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona, ok := r.personas[personaID]
	if !ok {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

// CreatePersona is unused in these tests.
func (r *notificationPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto persona_dtos.CreatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

// FindPersonasByUserID is unused in these tests.
func (r *notificationPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

// FindPersonaByHandle is unused in these tests.
func (r *notificationPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

// UpdatePersona is unused in these tests.
func (r *notificationPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto persona_dtos.UpdatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

// TestListNotificationsPipeListsOwnedPersonaNotifications verifies notifications can be listed for an owned persona.
func TestListNotificationsPipeListsOwnedPersonaNotifications(t *testing.T) {
	userID, personaID, actorID := uuid.New(), uuid.New(), uuid.New()
	repo := &notificationTestRepo{
		notifications: []*models.Notification{
			{
				ID:                 uuid.New(),
				RecipientUserID:    userID,
				RecipientPersonaID: personaID,
				ActorPersonaID:     &actorID,
				ActorPostingMode:   models.PublicPostingMode,
				NotificationType:   models.DirectMessageNotificationType,
				CreatedAt:          time.Now(),
			},
		},
		unreadCount: 1,
	}
	pipe := NewNotificationPipe(repo, &notificationPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID},
			actorID:   {ID: actorID, UserID: uuid.New(), Handle: "actor", DisplayName: "Actor"},
		},
	})

	res := pipe.ListNotificationsPipe(context.Background(), userID, personaID, 20, 0)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data == nil || len(res.Data.Notifications) != 1 {
		t.Fatal("expected one notification response")
	}
	if res.Data.UnreadCount != 1 {
		t.Fatalf("expected unread count 1, got %d", res.Data.UnreadCount)
	}
	if res.Data.Notifications[0].NotificationType != models.DirectMessageNotificationType {
		t.Fatalf("expected direct message notification, got %q", res.Data.Notifications[0].NotificationType)
	}
}

// TestMarkNotificationReadPipeMarksOwnedNotification verifies one notification can be marked read.
func TestMarkNotificationReadPipeMarksOwnedNotification(t *testing.T) {
	userID, personaID, actorID, notificationID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	now := time.Now()
	repo := &notificationTestRepo{
		markRead: &models.Notification{
			ID:                 notificationID,
			RecipientUserID:    userID,
			RecipientPersonaID: personaID,
			ActorPersonaID:     &actorID,
			ActorPostingMode:   models.PublicPostingMode,
			IsRead:             true,
			ReadAt:             &now,
			NotificationType:   models.GroupMessageNotificationType,
			CreatedAt:          now,
		},
	}
	pipe := NewNotificationPipe(repo, &notificationPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID},
			actorID:   {ID: actorID, UserID: uuid.New(), Handle: "actor", DisplayName: "Actor"},
		},
	})

	res := pipe.MarkNotificationReadPipe(context.Background(), userID, notificationID, personaID)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data == nil || !res.Data.IsRead {
		t.Fatal("expected read notification response")
	}
}

// TestMarkAllNotificationsReadPipeMarksOwnedPersonaNotifications verifies one persona can mark all notifications as read.
func TestMarkAllNotificationsReadPipeMarksOwnedPersonaNotifications(t *testing.T) {
	userID, personaID := uuid.New(), uuid.New()
	repo := &notificationTestRepo{}
	pipe := NewNotificationPipe(repo, &notificationPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID},
		},
	})

	res := pipe.MarkAllNotificationsReadPipe(context.Background(), userID, personaID)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if repo.markAllCalls != 1 {
		t.Fatalf("expected one mark-all call, got %d", repo.markAllCalls)
	}
}

// TestMarkNotificationReadPipeReturnsNotFound verifies missing notifications map to the expected pipe message.
func TestMarkNotificationReadPipeReturnsNotFound(t *testing.T) {
	userID, personaID := uuid.New(), uuid.New()
	repo := &notificationTestRepo{}
	pipe := NewNotificationPipe(repo, &notificationPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID},
		},
	})

	res := pipe.MarkNotificationReadPipe(context.Background(), userID, uuid.New(), personaID)
	if res.Message != notification_messages.Notification_Not_Found {
		t.Fatalf("expected %q, got %q", notification_messages.Notification_Not_Found, res.Message)
	}
}

// TestGetUnreadCountPipeReturnsCount verifies unread count retrieval for one owned persona.
func TestGetUnreadCountPipeReturnsCount(t *testing.T) {
	userID, personaID := uuid.New(), uuid.New()
	pipe := NewNotificationPipe(&notificationTestRepo{unreadCount: 4}, &notificationPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID},
		},
	})

	res := pipe.GetUnreadCountPipe(context.Background(), userID, personaID)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.UnreadCount != 4 {
		t.Fatal("expected unread count response")
	}
}

// TestListNotificationsPipePreservesAnonymousActor verifies anonymous comment notifications stay anonymous.
func TestListNotificationsPipePreservesAnonymousActor(t *testing.T) {
	userID, personaID := uuid.New(), uuid.New()
	repo := &notificationTestRepo{
		notifications: []*models.Notification{{
			ID:                      uuid.New(),
			RecipientUserID:         userID,
			RecipientPersonaID:      personaID,
			ActorPostingMode:        models.AnonymousPostingMode,
			ActorAnonymousHandle:    "ghost_1234",
			ActorAnonymousAvatarKey: "avatar_a",
			NotificationType:        models.CommentNotificationType,
			CreatedAt:               time.Now(),
		}},
	}
	pipe := NewNotificationPipe(repo, &notificationPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID},
		},
	})

	res := pipe.ListNotificationsPipe(context.Background(), userID, personaID, 20, 0)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.Notifications[0].Actor.Anonymous == nil || res.Data.Notifications[0].Actor.Persona != nil {
		t.Fatal("expected anonymous actor response without public persona")
	}
}
