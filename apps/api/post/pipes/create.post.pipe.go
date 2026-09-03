package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/dtos"
	"github.com/emmanuella-codes/nox/post/messages"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/anonymousidentity"
	"github.com/google/uuid"
)

// CreatePostPipe validates the author profile and persists a new post.
func (p *PostPipe) CreatePostPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreatePostDTO) *shared.PipeRes[PostResponse] {
	dto.Body = strings.TrimSpace(dto.Body)
	dto.MediaURL = strings.TrimSpace(dto.MediaURL)
	dto.Location = strings.TrimSpace(dto.Location)
	if !validPostingMode(dto.PostingMode) {
		return shared.PipeError[PostResponse](messages.Invalid_Posting_Mode)
	}
	if dto.PersonaID == nil || *dto.PersonaID == uuid.Nil {
		return shared.PipeError[PostResponse](messages.Persona_Required)
	}
	profile, err := p.personaRepo.FindPersonaByID(ctx, *dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[PostResponse](err, "post.find_profile")
	}
	if !profile.IsOwnedBy(userID) {
		return shared.PipeError[PostResponse](messages.Forbidden)
	}
	if err := p.validatePostMedia(ctx, userID, profile.ID, dto.MediaAssetIDs); err != nil {
		if err == errInvalidPostMedia {
			return shared.PipeError[PostResponse](messages.Invalid_Payload)
		}
		return pipeInternalError[PostResponse](err, "post.validate_media")
	}
	tags := hashtag_repo.ExtractTags(dto.Body)
	post, err := p.postRepo.CreatePostWithHashtagsAndMedia(ctx, userID, dto, tags, dto.MediaAssetIDs)
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.create")
	}
	var identity *models.AnonymousThreadIdentity
	if post.PostingMode == models.AnonymousPostingMode && post.PersonaID != nil {
		identity, err = p.postRepo.EnsureAnonymousThreadIdentity(
			ctx,
			post.ID,
			userID,
			*post.PersonaID,
			anonymousidentity.GenerateHandle(),
			anonymousidentity.GenerateAvatarKey(),
		)
		if err != nil {
			return pipeInternalError[PostResponse](err, "post.anonymous_identity")
		}
	}
	response := postResponse(post, publicAuthorProfile(post, profile), identity)
	response.Hashtags = tags
	if len(dto.MediaAssetIDs) > 0 {
		mediaByPost, err := p.postRepo.FindMediaAssetsByPostIDs(ctx, []uuid.UUID{post.ID})
		if err != nil {
			return pipeInternalError[PostResponse](err, "post.media_response")
		}
		response.Media = mediaByPost[post.ID]
	}
	return shared.PipeSuccess(messages.Post_Created, &response)
}

// publicAuthorProfile returns the public profile only for public posting mode.
func publicAuthorProfile(post *models.Post, profile *models.Persona) *models.Persona {
	if post.PostingMode != models.PublicPostingMode {
		return nil
	}
	return profile
}
