package hashtag

import (
	"context"
	"regexp"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HashtagRepository interface {
	SyncPostHashtags(ctx context.Context, postID uuid.UUID, tags []string) error
	DeletePostHashtags(ctx context.Context, postID uuid.UUID) error
	FindTagsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]string, error)
	FindTrending(ctx context.Context, limit int) ([]*models.Hashtag, error)
	FindByTag(ctx context.Context, tag string) (*models.Hashtag, error)
	FindPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*models.Post, error)
	Search(ctx context.Context, query string, limit int, offset int) ([]*models.Hashtag, error)
}

func NewHashtagRepository(db *pgxpool.Pool) HashtagRepository {
	return newPgRepository(db)
}

var hashtagPattern = regexp.MustCompile(`#([A-Za-z0-9_][A-Za-z0-9_-]*)`)

func ExtractTags(body string) []string {
	matches := hashtagPattern.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool, len(matches))
	tags := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match[1]) > 50 {
			continue
		}
		tag := NormalizeTag(match[1])
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		tags = append(tags, tag)
	}
	return tags
}

func NormalizeTag(tag string) string {
	tag = strings.TrimLeft(strings.TrimSpace(tag), "#")
	tag = strings.Trim(tag, "-_")
	return strings.ToLower(tag)
}
