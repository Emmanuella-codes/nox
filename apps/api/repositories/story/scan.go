package story

import (
	"errors"

	"github.com/emmanuella-codes/nox/models"
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
