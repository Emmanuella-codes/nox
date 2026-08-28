package post

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// EnsureAnonymousThreadIdentity creates or reuses one anonymous identity per thread owner.
func (r *pgRepository) EnsureAnonymousThreadIdentity(ctx context.Context, threadID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, anonymousHandle string, anonymousAvatarKey string) (*models.AnonymousThreadIdentity, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO anonymous_thread_identities (thread_id, user_id, persona_id, anonymous_handle, anonymous_avatar_key)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (thread_id, user_id) DO UPDATE
		SET anonymous_handle = anonymous_thread_identities.anonymous_handle,
		    anonymous_avatar_key = anonymous_thread_identities.anonymous_avatar_key
		RETURNING id, thread_id, user_id, persona_id, anonymous_handle, anonymous_avatar_key, created_at
	`, threadID, userID, personaID, anonymousHandle, anonymousAvatarKey)
	return scanAnonymousThreadIdentity(row)
}

// FindAnonymousThreadIdentities fetches thread identities for the supplied users.
func (r *pgRepository) FindAnonymousThreadIdentities(ctx context.Context, threadID uuid.UUID, userIDs []uuid.UUID) (map[uuid.UUID]*models.AnonymousThreadIdentity, error) {
	identities := make(map[uuid.UUID]*models.AnonymousThreadIdentity)
	if len(userIDs) == 0 {
		return identities, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, thread_id, user_id, persona_id, anonymous_handle, anonymous_avatar_key, created_at
		FROM anonymous_thread_identities
		WHERE thread_id = $1 AND user_id = ANY($2)
	`, threadID, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		identity, err := scanAnonymousThreadIdentity(rows)
		if err != nil {
			return nil, err
		}
		identities[identity.UserID] = identity
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return identities, nil
}
