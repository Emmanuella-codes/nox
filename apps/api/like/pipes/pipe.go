package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/like/dtos"
	"github.com/emmanuella-codes/nox/like/messages"
	"github.com/emmanuella-codes/nox/models"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type LikePipe struct {
	likeRepo    like_repo.LikeRepository
	personaRepo persona_repo.PersonaRepository
	postRepo    post_repo.PostRepository
}

func NewLikePipe(likeRepo like_repo.LikeRepository, personaRepo persona_repo.PersonaRepository, postRepo post_repo.PostRepository) *LikePipe {
	return &LikePipe{likeRepo: likeRepo, personaRepo: personaRepo, postRepo: postRepo}
}

func (p *LikePipe) LikePostPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.LikePostDTO) *shared.PipeRes[any] {
	if res := p.validatePersonaAndPost(ctx, userID, postID, dto.PersonaID); res != nil {
		return res
	}
	if err := p.likeRepo.LikePost(ctx, dto.PersonaID, postID); err != nil {
		return pipeInternalError[any](err, "like.post")
	}
	return pipeSuccess[any](messages.Post_Liked, nil)
}

func (p *LikePipe) UnlikePostPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.LikePostDTO) *shared.PipeRes[any] {
	if res := p.validatePersonaAndPost(ctx, userID, postID, dto.PersonaID); res != nil {
		return res
	}
	if err := p.likeRepo.UnlikePost(ctx, dto.PersonaID, postID); err != nil {
		return pipeInternalError[any](err, "unlike.post")
	}
	return pipeSuccess[any](messages.Post_Unliked, nil)
}

func (p *LikePipe) validatePersonaAndPost(ctx context.Context, userID uuid.UUID, postID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[any] {
	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[any](messages.Post_Not_Found)
		}
		return pipeInternalError[any](err, "like.find_post")
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[any](messages.Persona_Not_Found)
		}
		return pipeInternalError[any](err, "like.find_persona")
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return pipeError[any](messages.Forbidden)
	}

	return nil
}

func pipeSuccess[T any](message shared.PipeMessage, data *T) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: true, Message: message, Data: data}
}

func pipeError[T any](message shared.PipeMessage) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: false, Message: message}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	if err != nil {
		log.Error().Err(err).Str("operation", operation).Msg("like internal error")
	}
	return pipeError[T](messages.Internal_Error)
}
