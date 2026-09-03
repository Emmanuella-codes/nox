package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/preference/dtos"
	"github.com/emmanuella-codes/nox/preference/messages"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PreferencePipe) BlockUserPipe(ctx context.Context, userID uuid.UUID, dto dtos.PersonaTargetDTO) *shared.PipeRes[any] {
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
	if err := p.preferenceRepo.BlockUser(ctx, actor.UserID, target.UserID); err != nil && err != preference_repo.ErrPreferenceAlreadyExists {
		if err == preference_repo.ErrSelfTarget {
			return shared.PipeError[any](messages.Invalid_Payload)
		}
		return pipeInternalError[any](err, "preference.block")
	}
	if err := p.preferenceRepo.DeleteFollowRelationsBetweenUsers(ctx, actor.UserID, target.UserID); err != nil {
		return pipeInternalError[any](err, "preference.block_cleanup_follows")
	}
	return shared.PipeSuccess[any](messages.Blocked_Successfully, nil)
}

func (p *PreferencePipe) UnblockUserPipe(ctx context.Context, userID uuid.UUID, dto dtos.PersonaTargetDTO) *shared.PipeRes[any] {
	actor, res := p.ownedPersona(ctx, userID, dto.PersonaID)
	if res != nil {
		return res
	}
	target, res := p.targetPersona(ctx, dto.TargetPersonaID)
	if res != nil {
		return res
	}
	if err := p.preferenceRepo.UnblockUser(ctx, actor.UserID, target.UserID); err != nil && err != preference_repo.ErrPreferenceNotFound {
		return pipeInternalError[any](err, "preference.unblock")
	}
	return shared.PipeSuccess[any](messages.Unblocked_Successfully, nil)
}
