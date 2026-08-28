package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// requireMutualFollow enforces that two public profiles follow each other.
func (p *MessagingPipe) requireMutualFollow(ctx context.Context, actorPersonaID uuid.UUID, targetPersonaID uuid.UUID) shared.PipeMessage {
	if p.followRepo == nil {
		return messages.Internal_Error
	}
	followsTarget, err := p.followRepo.IsFollowing(ctx, actorPersonaID, targetPersonaID)
	if err != nil {
		return messages.Internal_Error
	}
	followedByTarget, err := p.followRepo.IsFollowing(ctx, targetPersonaID, actorPersonaID)
	if err != nil {
		return messages.Internal_Error
	}
	if !followsTarget || !followedByTarget {
		return messages.Forbidden
	}
	return ""
}

// requireMutualFollows enforces mutual follow rules between the acting profile and invited group members.
func (p *MessagingPipe) requireMutualFollows(ctx context.Context, actorPersonaID uuid.UUID, members []*models.Persona) shared.PipeMessage {
	for _, member := range members {
		if member.ID == actorPersonaID {
			continue
		}
		if message := p.requireMutualFollow(ctx, actorPersonaID, member.ID); message != "" {
			return message
		}
	}
	return ""
}
