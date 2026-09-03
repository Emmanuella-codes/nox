package story

import (
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type storyScanner interface {
	Scan(dest ...any) error
}

type storyRows interface {
	storyScanner
	Next() bool
	Err() error
}

func scanStories(rows storyRows) ([]*models.Story, error) {
	var stories []*models.Story
	for rows.Next() {
		story, err := scanStory(rows)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}
	return stories, rows.Err()
}

func scanStory(scanner storyScanner) (*models.Story, error) {
	var story models.Story
	err := scanner.Scan(
		&story.ID,
		&story.EventID,
		&story.OwnerUserID,
		&story.OwnerPersonaID,
		&story.Title,
		&story.ContributionMode,
		&story.TotalDurationSeconds,
		&story.ExpiresAt,
		&story.CreatedAt,
		&story.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

func scanStoryItems(rows storyRows) ([]*models.StoryItem, error) {
	var items []*models.StoryItem
	for rows.Next() {
		item, err := scanStoryItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanStoryItem(scanner storyScanner) (*models.StoryItem, error) {
	var item models.StoryItem
	err := scanner.Scan(
		&item.ID,
		&item.StoryID,
		&item.MediaAssetID,
		&item.ContributorUserID,
		&item.ContributorPersonaID,
		&item.PostingMode,
		&item.AnonymousLabel,
		&item.DurationSeconds,
		&item.Position,
		&item.ExpiresAt,
		&item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func scanEventHighlightStories(rows storyRows) ([]*models.EventHighlightStory, error) {
	var highlights []*models.EventHighlightStory
	for rows.Next() {
		highlight, err := scanEventHighlightStory(rows)
		if err != nil {
			return nil, err
		}
		highlights = append(highlights, highlight)
	}
	return highlights, rows.Err()
}

func scanEventHighlightStory(scanner storyScanner) (*models.EventHighlightStory, error) {
	var highlight models.EventHighlightStory
	err := scanner.Scan(
		&highlight.ID,
		&highlight.EventID,
		&highlight.StoryID,
		&highlight.AddedByPersonaID,
		&highlight.Position,
		&highlight.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &highlight, nil
}

// scanStoryUUIDs scans one-column uuid row sets into a slice.
func scanStoryUUIDs(rows storyRows) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// scanStoryContributionRequests scans many contribution request rows into model values.
func scanStoryContributionRequests(rows storyRows) ([]*models.StoryContributionRequest, error) {
	var requests []*models.StoryContributionRequest
	for rows.Next() {
		request, err := scanStoryContributionRequest(rows)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

// scanStoryContributionRequest scans one contribution request row into the model shape.
func scanStoryContributionRequest(scanner storyScanner) (*models.StoryContributionRequest, error) {
	var request models.StoryContributionRequest
	err := scanner.Scan(
		&request.ID,
		&request.StoryID,
		&request.MediaAssetID,
		&request.ContributorUserID,
		&request.ContributorPersonaID,
		&request.Status,
		&request.ReviewedByPersonaID,
		&request.StoryItemID,
		&request.CreatedAt,
		&request.ReviewedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStoryContributionRequestNotFound
		}
		return nil, err
	}
	return &request, nil
}

func mapStoryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrStoryNotFound
	}
	return err
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}
