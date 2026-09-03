package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/like/messages"
	"github.com/emmanuella-codes/nox/models"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type LikePipe struct {
	likeRepo              like_repo.LikeRepository
	personaRepo           persona_repo.PersonaRepository
	postRepo              post_repo.PostRepository
	notificationRepo      notification_repo.NotificationRepository
	notificationPublisher interface {
		PublishCreatedNotification(ctx context.Context, notification *models.Notification)
	}
}

// NewLikePipe builds the like orchestration layer from repositories.
func NewLikePipe(likeRepo like_repo.LikeRepository, personaRepo persona_repo.PersonaRepository, postRepo post_repo.PostRepository, deps ...any) *LikePipe {
	pipe := &LikePipe{likeRepo: likeRepo, personaRepo: personaRepo, postRepo: postRepo}
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

// validatePersonaAndPost checks the acting profile and target post before a like action.
func (p *LikePipe) validatePersonaAndPost(ctx context.Context, userID uuid.UUID, postID uuid.UUID, personaID uuid.UUID) (*models.Persona, *models.Post, *shared.PipeRes[any]) {
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return nil, nil, shared.PipeError[any](messages.Post_Not_Found)
		}
		return nil, nil, pipeInternalError[any](err, "like.find_post")
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, nil, shared.PipeError[any](messages.Persona_Not_Found)
		}
		return nil, nil, pipeInternalError[any](err, "like.find_persona")
	}
	if persona.UserID != userID {
		return nil, nil, shared.PipeError[any](messages.Forbidden)
	}

	return persona, post, nil
}

// pipeInternalError maps internal like errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "like", operation, messages.Internal_Error)
}
