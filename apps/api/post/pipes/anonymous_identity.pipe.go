package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared/anonymousidentity"
	"github.com/google/uuid"
)

// ensures a per-thread anonymous identity for one post owner.
func (p *PostPipe) anonymousPostIdentity(ctx context.Context, post *models.Post) (*models.AnonymousThreadIdentity, error) {
	if post.PostingMode != models.AnonymousPostingMode || post.PersonaID == nil {
		return nil, nil
	}
	return p.postRepo.EnsureAnonymousThreadIdentity(
		ctx,
		post.ID,
		post.AuthorUserID,
		*post.PersonaID,
		anonymousidentity.GenerateHandle(),
		anonymousidentity.GenerateAvatarKey(),
	)
}

// anonymousPostIdentities ensures or fetches anonymous identities for a list of posts.
func (p *PostPipe) anonymousPostIdentities(ctx context.Context, posts []*models.Post) (map[uuid.UUID]*models.AnonymousThreadIdentity, error) {
	identities := make(map[uuid.UUID]*models.AnonymousThreadIdentity)
	for _, post := range posts {
		identity, err := p.anonymousPostIdentity(ctx, post)
		if err != nil {
			return nil, err
		}
		if identity != nil {
			identities[post.ID] = identity
		}
	}
	return identities, nil
}

// anonymousHandleValue extracts the anonymous handle or returns the default label.
func anonymousHandleValue(identity *models.AnonymousThreadIdentity) string {
	if identity == nil || identity.AnonymousHandle == "" {
		return "anonymous"
	}
	return identity.AnonymousHandle
}

// anonymousAvatarValue extracts the anonymous avatar key or returns the default key.
func anonymousAvatarValue(identity *models.AnonymousThreadIdentity) string {
	if identity == nil || identity.AnonymousAvatarKey == "" {
		return "mask_01"
	}
	return identity.AnonymousAvatarKey
}
