package pipes

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/messages"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

var errInvalidPostMedia = errors.New("invalid post media")

type PostPipe struct {
	postRepo    post_repo.PostRepository
	personaRepo persona_repo.PersonaRepository
	likeRepo    like_repo.LikeRepository
	hashtagRepo hashtag_repo.HashtagRepository
	mediaRepo   media_repo.MediaRepository
}

type PostResponse struct {
	ID           string               `json:"id"`
	Author       PostAuthor           `json:"author"`
	EventID      *string              `json:"event_id,omitempty"`
	Body         string               `json:"body"`
	PostType     models.PostType      `json:"post_type"`
	MediaURL     string               `json:"media_url,omitempty"`
	MediaType    models.MediaType     `json:"media_type,omitempty"`
	Location     string               `json:"location,omitempty"`
	LikeCount    int                  `json:"like_count"`
	CommentCount int                  `json:"comment_count"`
	RepostCount  int                  `json:"repost_count"`
	IsLiked      bool                 `json:"is_liked"`
	IsRepost     bool                 `json:"is_repost"`
	RepostOf     *string              `json:"repost_of,omitempty"`
	Hashtags     []string             `json:"hashtags"`
	Media        []*models.MediaAsset `json:"media"`
	CreatedAt    time.Time            `json:"created_at"`
}

type PostAuthor struct {
	Mode      models.PostingMode   `json:"mode"`
	Persona   *PostPersonaAuthor   `json:"persona,omitempty"`
	Anonymous *PostAnonymousAuthor `json:"anonymous,omitempty"`
}

type PostPersonaAuthor struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type PostAnonymousAuthor struct {
	Handle    string `json:"handle"`
	AvatarKey string `json:"avatar_key"`
}

// NewPostPipe builds the post orchestration layer from repositories.
func NewPostPipe(postRepo post_repo.PostRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *PostPipe {
	var likes like_repo.LikeRepository
	var hashtags hashtag_repo.HashtagRepository
	var media media_repo.MediaRepository
	for _, dep := range deps {
		switch typed := dep.(type) {
		case like_repo.LikeRepository:
			likes = typed
		case hashtag_repo.HashtagRepository:
			hashtags = typed
		case media_repo.MediaRepository:
			media = typed
		}
	}
	return &PostPipe{postRepo: postRepo, personaRepo: personaRepo, likeRepo: likes, hashtagRepo: hashtags, mediaRepo: media}
}

// pipeInternalError maps internal post errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "post", operation, messages.Internal_Error)
}

// validPostingMode validates the supported post publishing modes.
func validPostingMode(mode models.PostingMode) bool {
	return mode == models.PublicPostingMode || mode == models.AnonymousPostingMode
}

// postResponse maps one post model into the API response shape.
func postResponse(post *models.Post, persona *models.Persona, anonymousIdentity *models.AnonymousThreadIdentity) PostResponse {
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
		Media:        []*models.MediaAsset{},
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
	if post.PostingMode == models.AnonymousPostingMode {
		res.Author.Anonymous = &PostAnonymousAuthor{
			Handle:    anonymousHandleValue(anonymousIdentity),
			AvatarKey: anonymousAvatarValue(anonymousIdentity),
		}
	}
	return res
}

// postIDs returns the ids for a slice of posts.
func postIDs(posts []*models.Post) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, post.ID)
	}
	return ids
}

// validatePostMedia checks that attached media belongs to the current user profile and is ready.
func (p *PostPipe) validatePostMedia(ctx context.Context, userID uuid.UUID, ownerPersonaID uuid.UUID, mediaAssetIDs []uuid.UUID) error {
	if len(mediaAssetIDs) == 0 {
		return nil
	}
	if p.mediaRepo == nil || ownerPersonaID == uuid.Nil || len(mediaAssetIDs) > 4 {
		return errInvalidPostMedia
	}
	seen := map[uuid.UUID]bool{}
	for _, mediaAssetID := range mediaAssetIDs {
		if mediaAssetID == uuid.Nil || seen[mediaAssetID] {
			return errInvalidPostMedia
		}
		seen[mediaAssetID] = true
		asset, err := p.mediaRepo.FindMediaAssetByID(ctx, mediaAssetID)
		if err != nil {
			return err
		}
		if asset.OwnerUserID != userID || asset.OwnerPersonaID != ownerPersonaID || asset.ProcessingStatus != models.ReadyMediaStatus {
			return errInvalidPostMedia
		}
		if asset.MediaKind != models.ImageMediaKind && asset.MediaKind != models.VideoMediaKind {
			return errInvalidPostMedia
		}
	}
	return nil
}
