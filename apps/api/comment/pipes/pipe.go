package pipes

import (
	"github.com/emmanuella-codes/nox/comment/messages"
	comment_repo "github.com/emmanuella-codes/nox/repositories/comment"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/rs/zerolog/log"
)

type CommentPipe struct {
	commentRepo comment_repo.CommentRepository
	personaRepo persona_repo.PersonaRepository
	postRepo    post_repo.PostRepository
}

func NewCommentPipe(commentRepo comment_repo.CommentRepository, personaRepo persona_repo.PersonaRepository, postRepo post_repo.PostRepository) *CommentPipe {
	return &CommentPipe{commentRepo: commentRepo, personaRepo: personaRepo, postRepo: postRepo}
}

func pipeSuccess[T any](message shared.PipeMessage, data *T) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: true, Message: message, Data: data}
}

func pipeError[T any](message shared.PipeMessage) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: false, Message: message}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	if err != nil {
		log.Error().Err(err).Str("operation", operation).Msg("comment internal error")
	}
	return pipeError[T](messages.Internal_Error)
}
