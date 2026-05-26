package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/models"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type FollowPipe struct {
	followRepo  follow_repo.FollowRepository
	personaRepo persona_repo.PersonaRepository
}

type FollowStatusResponse struct {
	IsFollowing bool `json:"is_following"`
}

func NewFollowPipe(followRepo follow_repo.FollowRepository, personaRepo persona_repo.PersonaRepository) *FollowPipe {
	return &FollowPipe{followRepo: followRepo, personaRepo: personaRepo}
}

func (p *FollowPipe) validateFollowAction(ctx context.Context, userID uuid.UUID, followerPersonaID uuid.UUID, targetPersonaID uuid.UUID) *shared.PipeRes[any] {
	followerPersona, res := p.validateOwnedVisiblePersona(ctx, userID, followerPersonaID)
	if res != nil {
		return res
	}

	targetPersona, res := p.findVisiblePersona(ctx, targetPersonaID)
	if res != nil {
		return res
	}

	if followerPersona.ID == targetPersona.ID {
		return shared.PipeError[any](messages.Self_Follow_Not_Allowed)
	}

	return nil
}

func (p *FollowPipe) validateOwnedVisiblePersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[any]) {
	persona, res := p.findVisiblePersona(ctx, personaID)
	if res != nil {
		return nil, res
	}
	if persona.UserID != userID {
		return nil, shared.PipeError[any](messages.Forbidden)
	}
	return persona, nil
}

func (p *FollowPipe) validateVisiblePersona(ctx context.Context, personaID uuid.UUID) *shared.PipeRes[any] {
	_, res := p.findVisiblePersona(ctx, personaID)
	return res
}

func (p *FollowPipe) findVisiblePersona(ctx context.Context, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[any]) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, shared.PipeError[any](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[any](err, "follow.find_persona")
	}
	if persona.PersonaType != models.VisiblePersonaType {
		return nil, shared.PipeError[any](messages.Persona_Not_Found)
	}
	return persona, nil
}

func (p *FollowPipe) mapFollowError(err error, operation string) *shared.PipeRes[any] {
	switch err {
	case follow_repo.ErrAlreadyFollowing:
		return shared.PipeError[any](messages.Already_Following)
	case follow_repo.ErrNotFollowing:
		return shared.PipeError[any](messages.Not_Following)
	case follow_repo.ErrSelfFollow:
		return shared.PipeError[any](messages.Self_Follow_Not_Allowed)
	default:
		return pipeInternalError[any](err, operation)
	}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "follow", operation, messages.Internal_Error)
}
