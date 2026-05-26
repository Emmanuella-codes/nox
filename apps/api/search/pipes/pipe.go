package pipes

import (
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
	searchmessages "github.com/emmanuella-codes/nox/search/messages"
	"github.com/emmanuella-codes/nox/shared"
)

type SearchPipe struct {
	repo        searchrepo.SearchRepository
	likeRepo    like_repo.LikeRepository
	personaRepo persona_repo.PersonaRepository
}

func NewSearchPipe(repo searchrepo.SearchRepository, likeRepo like_repo.LikeRepository, personaRepo persona_repo.PersonaRepository) *SearchPipe {
	return &SearchPipe{repo: repo, likeRepo: likeRepo, personaRepo: personaRepo}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "search", operation, searchmessages.Internal_Error)
}
