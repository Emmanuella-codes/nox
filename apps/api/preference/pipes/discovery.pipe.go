package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/preference/dtos"
	"github.com/emmanuella-codes/nox/preference/messages"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PreferencePipe) AddDiscoverySuppressionPipe(ctx context.Context, userID uuid.UUID, dto dtos.DiscoverySuppressionDTO) *shared.PipeRes[any] {
	actor, res := p.ownedPersona(ctx, userID, dto.PersonaID)
	if res != nil {
		return res
	}
	if !validDiscoveryTargetType(dto.TargetType) {
		return shared.PipeError[any](messages.Invalid_Payload)
	}
	if targetErr := p.ensureDiscoveryTargetExists(ctx, actor.UserID, dto.TargetType, dto.TargetID); targetErr != nil {
		return targetErr
	}
	if err := p.preferenceRepo.AddDiscoverySuppression(ctx, actor.UserID, dto.TargetType, dto.TargetID); err != nil && err != preference_repo.ErrPreferenceAlreadyExists {
		return pipeInternalError[any](err, "preference.add_discovery_suppression")
	}
	return shared.PipeSuccess[any](messages.Discovery_Suppression_Added, nil)
}

func (p *PreferencePipe) RemoveDiscoverySuppressionPipe(ctx context.Context, userID uuid.UUID, dto dtos.DiscoverySuppressionDTO) *shared.PipeRes[any] {
	actor, res := p.ownedPersona(ctx, userID, dto.PersonaID)
	if res != nil {
		return res
	}
	if !validDiscoveryTargetType(dto.TargetType) {
		return shared.PipeError[any](messages.Invalid_Payload)
	}
	if err := p.preferenceRepo.RemoveDiscoverySuppression(ctx, actor.UserID, dto.TargetType, dto.TargetID); err != nil && err != preference_repo.ErrPreferenceNotFound {
		return pipeInternalError[any](err, "preference.remove_discovery_suppression")
	}
	return shared.PipeSuccess[any](messages.Discovery_Suppression_Removed, nil)
}

func (p *PreferencePipe) ensureDiscoveryTargetExists(ctx context.Context, viewerUserID uuid.UUID, targetType models.DiscoverySuppressionTargetType, targetID uuid.UUID) *shared.PipeRes[any] {
	switch targetType {
	case models.PersonaSuppressionTargetType:
		persona, res := p.targetPersona(ctx, targetID)
		if res != nil {
			if res.Message == messages.Persona_Not_Found {
				return shared.PipeError[any](messages.Target_Not_Found)
			}
			return res
		}
		if persona.UserID == viewerUserID {
			return shared.PipeError[any](messages.Invalid_Payload)
		}
	case models.PostSuppressionTargetType:
		if p.postRepo == nil {
			return pipeInternalError[any](post_repo.ErrPostNotFound, "preference.discovery_post_repo")
		}
		post, err := p.postRepo.FindPostByID(ctx, targetID)
		if err != nil {
			if err == post_repo.ErrPostNotFound {
				return shared.PipeError[any](messages.Target_Not_Found)
			}
			return pipeInternalError[any](err, "preference.discovery_post_target")
		}
		if post.AuthorUserID == viewerUserID {
			return shared.PipeError[any](messages.Invalid_Payload)
		}
	case models.EventSuppressionTargetType:
		if p.eventRepo == nil {
			return pipeInternalError[any](event_repo.ErrEventNotFound, "preference.discovery_event_repo")
		}
		event, err := p.eventRepo.FindEventByID(ctx, targetID)
		if err != nil {
			if err == event_repo.ErrEventNotFound {
				return shared.PipeError[any](messages.Target_Not_Found)
			}
			return pipeInternalError[any](err, "preference.discovery_event_target")
		}
		organizer, res := p.targetPersona(ctx, event.OrganizerID)
		if res != nil {
			return shared.PipeError[any](messages.Target_Not_Found)
		}
		if organizer.UserID == viewerUserID {
			return shared.PipeError[any](messages.Invalid_Payload)
		}
	case models.SetSuppressionTargetType:
		if p.setRepo == nil {
			return pipeInternalError[any](set_repo.ErrSetNotFound, "preference.discovery_set_repo")
		}
		set, err := p.setRepo.FindSetByID(ctx, targetID)
		if err != nil {
			if err == set_repo.ErrSetNotFound {
				return shared.PipeError[any](messages.Target_Not_Found)
			}
			return pipeInternalError[any](err, "preference.discovery_set_target")
		}
		if set.AuthorUserID == viewerUserID {
			return shared.PipeError[any](messages.Invalid_Payload)
		}
	default:
		return shared.PipeError[any](messages.Invalid_Payload)
	}
	return nil
}
