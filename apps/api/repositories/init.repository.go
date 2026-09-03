package repositories

import (
	"github.com/emmanuella-codes/nox/repositories/comment"
	"github.com/emmanuella-codes/nox/repositories/crew"
	"github.com/emmanuella-codes/nox/repositories/event"
	"github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/emmanuella-codes/nox/repositories/like"
	"github.com/emmanuella-codes/nox/repositories/media"
	"github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/repositories/preference"
	"github.com/emmanuella-codes/nox/repositories/search"
	"github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User         user.UserRepository
	Persona      persona.PersonaRepository
	Post         post.PostRepository
	Preference   preference.PreferenceRepository
	Comment      comment.CommentRepository
	Crew         crew.CrewRepository
	Like         like.LikeRepository
	Event        event.EventRepository
	Media        media.MediaRepository
	Messaging    messaging.MessagingRepository
	Notification notification.NotificationRepository
	Set          set.SetRepository
	Story        story.StoryRepository
	Search       search.SearchRepository
	Follow       follow.FollowRepository
	Hashtag      hashtag.HashtagRepository
}

func Init(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:         user.NewUserRepository(pool),
		Persona:      persona.NewPersonaRepository(pool),
		Post:         post.NewPostRepository(pool),
		Preference:   preference.NewPreferenceRepository(pool),
		Comment:      comment.NewCommentRepository(pool),
		Crew:         crew.NewCrewRepository(pool),
		Like:         like.NewLikeRepository(pool),
		Event:        event.NewEventRepository(pool),
		Media:        media.NewMediaRepository(pool),
		Messaging:    messaging.NewMessagingRepository(pool),
		Notification: notification.NewNotificationRepository(pool),
		Set:          set.NewSetRepository(pool),
		Story:        story.NewStoryRepository(pool),
		Search:       search.NewSearchRepository(pool),
		Follow:       follow.NewFollowRepository(pool),
		Hashtag:      hashtag.NewHashtagRepository(pool),
	}
}
