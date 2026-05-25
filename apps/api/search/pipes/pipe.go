package pipes

import (
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
	searchmessages "github.com/emmanuella-codes/nox/search/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/rs/zerolog/log"
)

type SearchPipe struct {
	repo        searchrepo.Repository
	likeRepo    like_repo.LikeRepository
	personaRepo persona_repo.PersonaRepository
}

func NewSearchPipe(repo searchrepo.Repository, likeRepo like_repo.LikeRepository, personaRepo persona_repo.PersonaRepository) *SearchPipe {
	return &SearchPipe{repo: repo, likeRepo: likeRepo, personaRepo: personaRepo}
}

func pipeSuccess[T any](message shared.PipeMessage, data *T) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: true, Message: message, Data: data}
}

func pipeError[T any](message shared.PipeMessage) *shared.PipeRes[T] {
	return &shared.PipeRes[T]{Success: false, Message: message}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	if err != nil {
		log.Error().Err(err).Str("operation", operation).Msg("search internal error")
	}
	return pipeError[T](searchmessages.Internal_Error)
}
