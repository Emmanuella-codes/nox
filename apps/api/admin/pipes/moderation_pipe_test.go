package pipes

import (
	"context"
	"testing"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

func TestModerateContentPipeRequiresReasonForRemoval(t *testing.T) {
	actorID := uuid.New()
	pipe := moderationTestPipe(actorID)

	res := pipe.ModerateContentPipe(context.Background(), actorID, models.ModerationEntityPost, uuid.New(), dtos.ModerateContentDTO{
		Status: models.ModerationStatusRemoved,
	})
	if res.Success || res.Message != messages.Invalid_Payload {
		t.Fatalf("expected required reason error, got %#v", res)
	}
}

func TestModerateContentPipeUpdatesContentForActiveAdmin(t *testing.T) {
	actorID := uuid.New()
	pipe := moderationTestPipe(actorID)

	res := pipe.ModerateContentPipe(context.Background(), actorID, models.ModerationEntityPost, uuid.New(), dtos.ModerateContentDTO{
		Status: models.ModerationStatusHidden,
		Reason: "spam",
	})
	if !res.Success || res.Data == nil || res.Data.Status != models.ModerationStatusHidden {
		t.Fatalf("expected moderation update, got %#v", res)
	}
}

func TestModerateContentPipeRejectsInactiveAdmin(t *testing.T) {
	actorID := uuid.New()
	repo := &adminTestRepo{membership: &models.AdminMembership{UserID: actorID, IsActive: false}}
	pipe := &AdminPipe{adminRepo: repo}

	res := pipe.ModerateContentPipe(context.Background(), actorID, models.ModerationEntityPost, uuid.New(), dtos.ModerateContentDTO{
		Status: models.ModerationStatusActive,
	})
	if res.Success || res.Message != messages.Admin_Access_Denied {
		t.Fatalf("expected access denial, got %#v", res)
	}
}

func moderationTestPipe(actorID uuid.UUID) *AdminPipe {
	return &AdminPipe{
		adminRepo: &adminTestRepo{
			membership: &models.AdminMembership{UserID: actorID, IsActive: true},
		},
	}
}
