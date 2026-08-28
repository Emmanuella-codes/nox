package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/like/messages"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type LikePipe struct {
	likeRepo    like_repo.LikeRepository
	personaRepo persona_repo.PersonaRepository
	postRepo    post_repo.PostRepository
}

// NewLikePipe builds the like orchestration layer from repositories.
func NewLikePipe(likeRepo like_repo.LikeRepository, personaRepo persona_repo.PersonaRepository, postRepo post_repo.PostRepository) *LikePipe {
	return &LikePipe{likeRepo: likeRepo, personaRepo: personaRepo, postRepo: postRepo}
}

// validatePersonaAndPost checks the acting profile and target post before a like action.
func (p *LikePipe) validatePersonaAndPost(ctx context.Context, userID uuid.UUID, postID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[any] {
	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[any](messages.Post_Not_Found)
		}
		return pipeInternalError[any](err, "like.find_post")
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[any](messages.Persona_Not_Found)
		}
		return pipeInternalError[any](err, "like.find_persona")
	}
	if persona.UserID != userID {
		return shared.PipeError[any](messages.Forbidden)
	}

	return nil
}

// pipeInternalError maps internal like errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "like", operation, messages.Internal_Error)
}
