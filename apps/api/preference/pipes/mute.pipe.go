package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/preference/dtos"
	"github.com/emmanuella-codes/nox/preference/messages"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PreferencePipe) MuteUserPipe(ctx context.Context, userID uuid.UUID, dto dtos.PersonaTargetDTO) *shared.PipeRes[any] {
	actor, res := p.ownedPersona(ctx, userID, dto.PersonaID)
	if res != nil {
		return res
	}
	target, res := p.targetPersona(ctx, dto.TargetPersonaID)
	if res != nil {
		return res
	}
	if actor.UserID == target.UserID {
		return shared.PipeError[any](messages.Invalid_Payload)
	}
	if err := p.preferenceRepo.MuteUser(ctx, actor.UserID, target.UserID); err != nil && err != preference_repo.ErrPreferenceAlreadyExists {
		if err == preference_repo.ErrSelfTarget {
			return shared.PipeError[any](messages.Invalid_Payload)
		}
		return pipeInternalError[any](err, "preference.mute")
	}
	return shared.PipeSuccess[any](messages.Muted_Successfully, nil)
}

func (p *PreferencePipe) UnmuteUserPipe(ctx context.Context, userID uuid.UUID, dto dtos.PersonaTargetDTO) *shared.PipeRes[any] {
	actor, res := p.ownedPersona(ctx, userID, dto.PersonaID)
	if res != nil {
		return res
	}
	target, res := p.targetPersona(ctx, dto.TargetPersonaID)
	if res != nil {
		return res
	}
	if err := p.preferenceRepo.UnmuteUser(ctx, actor.UserID, target.UserID); err != nil && err != preference_repo.ErrPreferenceNotFound {
		return pipeInternalError[any](err, "preference.unmute")
	}
	return shared.PipeSuccess[any](messages.Unmuted_Successfully, nil)
}
