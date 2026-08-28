package story

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	storydtos "github.com/emmanuella-codes/nox/story/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateStoryContributionRequest stores one pending contribution request for a story.
func (r *pgRepository) CreateStoryContributionRequest(ctx context.Context, storyID uuid.UUID, contributorUserID uuid.UUID, dto storydtos.CreateStoryContributionRequestDTO) (*models.StoryContributionRequest, error) {
	if inUse, err := r.mediaAssetAlreadyAttached(ctx, dto.MediaAssetID); err != nil {
		return nil, err
	} else if inUse {
		return nil, ErrStoryMediaInUse
	}
	row := r.db.QueryRow(ctx, `
		INSERT INTO story_contribution_requests (
			story_id, media_asset_id, contributor_user_id, contributor_persona_id, status
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		          status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
	`, storyID, dto.MediaAssetID, contributorUserID, dto.ContributorPersonaID, models.PendingStoryContributionRequestStatus)
	request, err := scanStoryContributionRequest(row)
	if err != nil {
		if isUniqueViolation(err, "story_contribution_requests_pending_media_asset_idx") {
			return nil, ErrStoryContributionRequestPending
		}
		return nil, err
	}
	return request, nil
}

// FindStoryContributionRequestByID loads one contribution request for a story.
func (r *pgRepository) FindStoryContributionRequestByID(ctx context.Context, storyID uuid.UUID, requestID uuid.UUID) (*models.StoryContributionRequest, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
		FROM story_contribution_requests
		WHERE story_id = $1 AND id = $2
	`, storyID, requestID)
	return scanStoryContributionRequest(row)
}

// FindStoryContributionRequests lists contribution requests for one story with an optional status filter.
func (r *pgRepository) FindStoryContributionRequests(ctx context.Context, storyID uuid.UUID, status *models.StoryContributionRequestStatus, limit int, offset int) ([]*models.StoryContributionRequest, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
		FROM story_contribution_requests
		WHERE story_id = $1 AND ($2::text IS NULL OR status = $2)
		ORDER BY created_at DESC, id DESC
		LIMIT $3 OFFSET $4
	`, storyID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStoryContributionRequests(rows)
}

// AcceptStoryContributionRequest accepts one pending request and creates the linked story item in one transaction.
func (r *pgRepository) AcceptStoryContributionRequest(ctx context.Context, story *models.Story, requestID uuid.UUID, reviewerPersonaID uuid.UUID, durationSeconds int) (*models.StoryContributionRequest, *models.StoryItem, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	request, err := contributionRequestForReview(ctx, tx, story.ID, requestID)
	if err != nil {
		return nil, nil, err
	}
	var currentTotal int
	if err := tx.QueryRow(ctx, `SELECT total_duration_seconds FROM stories WHERE id = $1 FOR UPDATE`, story.ID).Scan(&currentTotal); err != nil {
		return nil, nil, mapStoryError(err)
	}
	if currentTotal+durationSeconds > maxStoryDurationSeconds {
		return nil, nil, ErrStoryDurationLimitExceeded
	}
	var position int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(MAX(position), 0) + 1 FROM story_items WHERE story_id = $1`, story.ID).Scan(&position); err != nil {
		return nil, nil, err
	}
	item, err := insertStoryItemInTx(ctx, tx, story.ID, request.MediaAssetID, request.ContributorUserID, request.ContributorPersonaID, models.PublicPostingMode, "", durationSeconds, position)
	if err != nil {
		return nil, nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE stories SET total_duration_seconds = total_duration_seconds + $2, updated_at = now() WHERE id = $1`, story.ID, durationSeconds); err != nil {
		return nil, nil, err
	}
	row := tx.QueryRow(ctx, `
		UPDATE story_contribution_requests
		SET status = $3,
		    reviewed_by_persona_id = $4,
		    story_item_id = $5,
		    reviewed_at = now()
		WHERE story_id = $1 AND id = $2
		RETURNING id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		          status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
	`, story.ID, requestID, models.AcceptedStoryContributionRequestStatus, reviewerPersonaID, item.ID)
	updatedRequest, err := scanStoryContributionRequest(row)
	if err != nil {
		return nil, nil, err
	}
	return updatedRequest, item, tx.Commit(ctx)
}

// RejectStoryContributionRequest rejects one pending request in one transaction.
func (r *pgRepository) RejectStoryContributionRequest(ctx context.Context, storyID uuid.UUID, requestID uuid.UUID, reviewerPersonaID uuid.UUID) (*models.StoryContributionRequest, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := contributionRequestForReview(ctx, tx, storyID, requestID); err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, `
		UPDATE story_contribution_requests
		SET status = $3,
		    reviewed_by_persona_id = $4,
		    reviewed_at = now()
		WHERE story_id = $1 AND id = $2
		RETURNING id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		          status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
	`, storyID, requestID, models.RejectedStoryContributionRequestStatus, reviewerPersonaID)
	request, err := scanStoryContributionRequest(row)
	if err != nil {
		return nil, err
	}
	return request, tx.Commit(ctx)
}

// contributionRequestForReview locks one request row and enforces that it is still pending.
func contributionRequestForReview(ctx context.Context, tx pgx.Tx, storyID uuid.UUID, requestID uuid.UUID) (*models.StoryContributionRequest, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, story_id, media_asset_id, contributor_user_id, contributor_persona_id,
		       status, reviewed_by_persona_id, story_item_id, created_at, reviewed_at
		FROM story_contribution_requests
		WHERE story_id = $1 AND id = $2
		FOR UPDATE
	`, storyID, requestID)
	request, err := scanStoryContributionRequest(row)
	if err != nil {
		return nil, err
	}
	if request.Status != models.PendingStoryContributionRequestStatus {
		return nil, ErrStoryContributionRequestReviewed
	}
	return request, nil
}

// mediaAssetAlreadyAttached checks whether one media asset is already used by a story item.
func (r *pgRepository) mediaAssetAlreadyAttached(ctx context.Context, mediaAssetID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM story_items WHERE media_asset_id = $1)`, mediaAssetID).Scan(&exists)
	return exists, err
}
