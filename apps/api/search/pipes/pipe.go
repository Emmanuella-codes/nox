package pipes

import (
	"fmt"

	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
	searchmessages "github.com/emmanuella-codes/nox/search/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type SearchPipe struct {
	repo        searchrepo.SearchRepository
	likeRepo    like_repo.LikeRepository
	personaRepo persona_repo.PersonaRepository
	hashtagRepo hashtag_repo.HashtagRepository
	followRepo  follow_repo.FollowRepository
	postRepo    post_repo.PostRepository
}

func NewSearchPipe(repo searchrepo.SearchRepository, likeRepo like_repo.LikeRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *SearchPipe {
	var hashtags hashtag_repo.HashtagRepository
	var follows follow_repo.FollowRepository
	var posts post_repo.PostRepository
	for _, dep := range deps {
		switch typed := dep.(type) {
		case hashtag_repo.HashtagRepository:
			hashtags = typed
		case follow_repo.FollowRepository:
			follows = typed
		case post_repo.PostRepository:
			posts = typed
		}
	}
	return &SearchPipe{repo: repo, likeRepo: likeRepo, personaRepo: personaRepo, hashtagRepo: hashtags, followRepo: follows, postRepo: posts}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "search", operation, searchmessages.Internal_Error)
}

func anonymousHandle() string {
	id := uuid.NewString()
	return fmt.Sprintf("ghost_%s", id[:8])
}
