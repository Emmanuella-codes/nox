package pipes

import (
	"context"
	"time"

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

type FollowListResponse struct {
	Limit      int                     `json:"limit"`
	Offset     int                     `json:"offset"`
	HasMore    bool                    `json:"has_more"`
	NextOffset *int                    `json:"next_offset,omitempty"`
	Personas   []FollowPersonaResponse `json:"personas"`
}

type FollowPersonaResponse struct {
	ID             string                 `json:"id"`
	Handle         string                 `json:"handle"`
	DisplayName    string                 `json:"display_name"`
	Bio            string                 `json:"bio"`
	AvatarURL      string                 `json:"avatar_url"`
	CoverURL       string                 `json:"cover_url"`
	PersonaType    models.PersonaType     `json:"persona_type"`
	Category       models.PersonaCategory `json:"category"`
	GenreTags      []string               `json:"genre_tags"`
	FollowerCount  int                    `json:"follower_count"`
	FollowingCount int                    `json:"following_count"`
	IsFollowing    bool                   `json:"is_following"`
	PostCount      int                    `json:"post_count"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
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

func followListResponse(personas []*models.Persona, options follow_repo.ListOptions, following map[uuid.UUID]bool) FollowListResponse {
	hasMore := len(personas) > options.Limit
	if hasMore {
		personas = personas[:options.Limit]
	}

	return FollowListResponse{
		Limit:      options.Limit,
		Offset:     options.Offset,
		HasMore:    hasMore,
		NextOffset: nextOffset(options, hasMore),
		Personas:   followPersonaResponses(personas, following),
	}
}

func nextOffset(options follow_repo.ListOptions, hasMore bool) *int {
	if !hasMore {
		return nil
	}
	next := options.Offset + options.Limit
	return &next
}

func followPersonaResponses(personas []*models.Persona, following map[uuid.UUID]bool) []FollowPersonaResponse {
	responses := make([]FollowPersonaResponse, 0, len(personas))
	for _, persona := range personas {
		responses = append(responses, FollowPersonaResponse{
			ID:             persona.ID.String(),
			Handle:         persona.Handle,
			DisplayName:    persona.DisplayName,
			Bio:            persona.Bio,
			AvatarURL:      persona.AvatarURL,
			CoverURL:       persona.CoverURL,
			PersonaType:    persona.PersonaType,
			Category:       persona.Category,
			GenreTags:      persona.GenreTags,
			FollowerCount:  persona.FollowerCount,
			FollowingCount: persona.FollowingCount,
			IsFollowing:    following[persona.ID],
			PostCount:      persona.PostCount,
			CreatedAt:      persona.CreatedAt,
			UpdatedAt:      persona.UpdatedAt,
		})
	}
	return responses
}

func personaIDs(personas []*models.Persona) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(personas))
	for _, persona := range personas {
		ids = append(ids, persona.ID)
	}
	return ids
}
