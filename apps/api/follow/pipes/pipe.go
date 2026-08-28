package pipes

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/models"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type FollowPipe struct {
	followRepo            follow_repo.FollowRepository
	personaRepo           persona_repo.PersonaRepository
	notificationRepo      notification_repo.NotificationRepository
	notificationPublisher interface {
		PublishCreatedNotification(ctx context.Context, notification *models.Notification)
	}
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

// NewFollowPipe builds the follow orchestration layer from repositories.
func NewFollowPipe(followRepo follow_repo.FollowRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *FollowPipe {
	pipe := &FollowPipe{followRepo: followRepo, personaRepo: personaRepo}
	for _, dep := range deps {
		if repo, ok := dep.(notification_repo.NotificationRepository); ok {
			pipe.notificationRepo = repo
		}
		if publisher, ok := dep.(interface {
			PublishCreatedNotification(ctx context.Context, notification *models.Notification)
		}); ok {
			pipe.notificationPublisher = publisher
		}
	}
	return pipe
}

// validateFollowAction checks the follower and target public profiles before a follow action.
func (p *FollowPipe) validateFollowAction(ctx context.Context, userID uuid.UUID, followerPersonaID uuid.UUID, targetPersonaID uuid.UUID) (*models.Persona, *models.Persona, *shared.PipeRes[any]) {
	followerPersona, res := p.validateOwnedProfile(ctx, userID, followerPersonaID)
	if res != nil {
		return nil, nil, res
	}

	targetPersona, res := p.findProfile(ctx, targetPersonaID)
	if res != nil {
		return nil, nil, res
	}

	if followerPersona.ID == targetPersona.ID {
		return nil, nil, shared.PipeError[any](messages.Self_Follow_Not_Allowed)
	}

	return followerPersona, targetPersona, nil
}

// validateOwnedProfile checks that the current user owns the acting public profile.
func (p *FollowPipe) validateOwnedProfile(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[any]) {
	persona, res := p.findProfile(ctx, personaID)
	if res != nil {
		return nil, res
	}
	if persona.UserID != userID {
		return nil, shared.PipeError[any](messages.Forbidden)
	}
	return persona, nil
}

// validateProfile checks that the target public profile exists.
func (p *FollowPipe) validateProfile(ctx context.Context, personaID uuid.UUID) *shared.PipeRes[any] {
	_, res := p.findProfile(ctx, personaID)
	return res
}

// findProfile fetches one public profile by id.
func (p *FollowPipe) findProfile(ctx context.Context, personaID uuid.UUID) (*models.Persona, *shared.PipeRes[any]) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, shared.PipeError[any](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[any](err, "follow.find_persona")
	}
	return persona, nil
}

// mapFollowError maps repository follow errors to pipe responses.
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

// pipeInternalError maps internal follow errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "follow", operation, messages.Internal_Error)
}

// followListResponse maps follow list pagination and profiles into the API shape.
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

// nextOffset calculates the next offset when more results are available.
func nextOffset(options follow_repo.ListOptions, hasMore bool) *int {
	if !hasMore {
		return nil
	}
	next := options.Offset + options.Limit
	return &next
}

// followPersonaResponses maps public profiles into the follow list response shape.
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

// personaIDs extracts ids from a slice of public profiles.
func personaIDs(personas []*models.Persona) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(personas))
	for _, persona := range personas {
		ids = append(ids, persona.ID)
	}
	return ids
}
