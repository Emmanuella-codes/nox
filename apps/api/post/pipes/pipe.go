package pipes

import (
	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/rs/zerolog/log"
)

type PostPipe struct {
	postRepo    post_repo.PostRepository
	personaRepo persona_repo.PersonaRepository
}

func NewPostPipe(postRepo post_repo.PostRepository, personaRepo persona_repo.PersonaRepository) *PostPipe {
	return &PostPipe{postRepo: postRepo, personaRepo: personaRepo}
}

func pipeSuccess[T any](message shared.PipeMessage, data *T) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: true, Message: message, Data: data}
}

func pipeError[T any](message shared.PipeMessage) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: false, Message: message}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	if err != nil {
		log.Error().Err(err).Str("operation", operation).Msg("post internal error")
	}
	return pipeError[T](messages.Internal_Error)
}
