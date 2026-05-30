package pipes

import (
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/messages"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type PostPipe struct {
	postRepo    post_repo.PostRepository
	personaRepo persona_repo.PersonaRepository
	likeRepo    like_repo.LikeRepository
	hashtagRepo hashtag_repo.HashtagRepository
}

func NewPostPipe(postRepo post_repo.PostRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *PostPipe {
	var likes like_repo.LikeRepository
	var hashtags hashtag_repo.HashtagRepository
	for _, dep := range deps {
		switch typed := dep.(type) {
		case like_repo.LikeRepository:
			likes = typed
		case hashtag_repo.HashtagRepository:
			hashtags = typed
		}
	}
	return &PostPipe{postRepo: postRepo, personaRepo: personaRepo, likeRepo: likes, hashtagRepo: hashtags}
}

type PostResponse struct {
	ID           string           `json:"id"`
	Author       PostAuthor       `json:"author"`
	EventID      *string          `json:"event_id,omitempty"`
	Body         string           `json:"body"`
	PostType     models.PostType  `json:"post_type"`
	MediaURL     string           `json:"media_url,omitempty"`
	MediaType    models.MediaType `json:"media_type,omitempty"`
	Location     string           `json:"location,omitempty"`
	LikeCount    int              `json:"like_count"`
	CommentCount int              `json:"comment_count"`
	RepostCount  int              `json:"repost_count"`
	IsLiked      bool             `json:"is_liked"`
	IsRepost     bool             `json:"is_repost"`
	RepostOf     *string          `json:"repost_of,omitempty"`
	Hashtags     []string         `json:"hashtags"`
	CreatedAt    time.Time        `json:"created_at"`
}

type PostAuthor struct {
	Mode    models.PostingMode `json:"mode"`
	Persona *PostPersonaAuthor `json:"persona,omitempty"`
}

type PostPersonaAuthor struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "post", operation, messages.Internal_Error)
}

func validPostingMode(mode models.PostingMode) bool {
	return mode == models.PublicPostingMode || mode == models.AnonymousPostingMode
}

func postResponse(post *models.Post, persona *models.Persona) PostResponse {
	res := PostResponse{
		ID:           post.ID.String(),
		Author:       PostAuthor{Mode: post.PostingMode},
		Body:         post.Body,
		PostType:     post.PostType,
		MediaURL:     post.MediaURL,
		MediaType:    post.MediaType,
		Location:     post.Location,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		RepostCount:  post.RepostCount,
		IsRepost:     post.IsRepost,
		Hashtags:     []string{},
		CreatedAt:    post.CreatedAt,
	}
	if post.EventID != nil {
		eventID := post.EventID.String()
		res.EventID = &eventID
	}
	if post.RepostOf != nil {
		repostOf := post.RepostOf.String()
		res.RepostOf = &repostOf
	}
	if post.PostingMode == models.PublicPostingMode && persona != nil {
		res.Author.Persona = &PostPersonaAuthor{
			ID:          persona.ID.String(),
			Handle:      persona.Handle,
			DisplayName: persona.DisplayName,
			AvatarURL:   persona.AvatarURL,
		}
	}
	return res
}

func postResponses(posts []*models.Post, personas map[string]*models.Persona) []PostResponse {
	responses := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		var persona *models.Persona
		if post.PersonaID != nil {
			persona = personas[post.PersonaID.String()]
		}
		response := postResponse(post, persona)
		responses = append(responses, response)
	}
	return responses
}

func postIDs(posts []*models.Post) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, post.ID)
	}
	return ids
}
