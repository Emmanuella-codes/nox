package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PostPipe) DeletePostPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID) *shared.PipeRes[any] {
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[any](messages.Post_Not_Found)
		}
		return pipeInternalError[any](err, "post.find_for_delete")
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, post.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[any](messages.Persona_Not_Found)
		}
		return pipeInternalError[any](err, "post.find_persona_for_delete")
	}
	if persona.UserID != userID {
		return pipeError[any](messages.Forbidden)
	}

	if err := p.postRepo.DeletePost(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[any](messages.Post_Not_Found)
		}
		return pipeInternalError[any](err, "post.delete")
	}

	return pipeSuccess[any](messages.Post_Deleted, nil)
}
