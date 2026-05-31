package pipes

import (
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
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
	hashtagRepo hashtag_repo.HashtagRepository
	followRepo  follow_repo.FollowRepository
}

func NewSearchPipe(repo searchrepo.SearchRepository, likeRepo like_repo.LikeRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *SearchPipe {
	var hashtags hashtag_repo.HashtagRepository
	var follows follow_repo.FollowRepository
	for _, dep := range deps {
		switch typed := dep.(type) {
		case hashtag_repo.HashtagRepository:
			hashtags = typed
		case follow_repo.FollowRepository:
			follows = typed
		}
	}
	return &SearchPipe{repo: repo, likeRepo: likeRepo, personaRepo: personaRepo, hashtagRepo: hashtags, followRepo: follows}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "search", operation, searchmessages.Internal_Error)
}
