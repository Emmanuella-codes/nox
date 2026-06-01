package repositories

import (
	"github.com/emmanuella-codes/nox/repositories/comment"
	"github.com/emmanuella-codes/nox/repositories/event"
	"github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/emmanuella-codes/nox/repositories/like"
	"github.com/emmanuella-codes/nox/repositories/media"
	"github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/repositories/search"
	"github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User    user.UserRepository
	Persona persona.PersonaRepository
	Post    post.PostRepository
	Comment comment.CommentRepository
	Like    like.LikeRepository
	Event   event.EventRepository
	Media   media.MediaRepository
	Set     set.SetRepository
	Story   story.StoryRepository
	Search  search.SearchRepository
	Follow  follow.FollowRepository
	Hashtag hashtag.HashtagRepository
}

func Init(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:    user.NewUserRepository(pool),
		Persona: persona.NewPersonaRepository(pool),
		Post:    post.NewPostRepository(pool),
		Comment: comment.NewCommentRepository(pool),
		Like:    like.NewLikeRepository(pool),
		Event:   event.NewEventRepository(pool),
		Media:   media.NewMediaRepository(pool),
		Set:     set.NewSetRepository(pool),
		Story:   story.NewStoryRepository(pool),
		Search:  search.NewSearchRepository(pool),
		Follow:  follow.NewFollowRepository(pool),
		Hashtag: hashtag.NewHashtagRepository(pool),
	}
}
