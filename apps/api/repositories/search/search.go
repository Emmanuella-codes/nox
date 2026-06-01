package search

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	// "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostResult struct {
	Post    *models.Post
	Persona *models.Persona
}

type Results struct {
	Personas []*models.Persona
	Posts    []*PostResult
	Events   []*models.Event
	Hashtags []*models.Hashtag
	HasMore  bool
}

type Options struct {
	Limit  int
	Offset int
}

type SearchRepository interface {
	Search(ctx context.Context, query string, options Options) (*Results, error)
}

func NewSearchRepository(db *pgxpool.Pool) SearchRepository {
	return newPgRepository(db)
}

func NormalizeOptions(options Options) Options {
	options.Limit = normalizeLimit(options.Limit)
	if options.Offset < 0 {
		options.Offset = 0
	}
	return options
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 30 {
		return 30
	}
	return limit
}

func textMatchParam(query string) string {
	return "%" + query + "%"
}

func tagMatchParam(query string) string {
	return "%" + query + "%"
}

func prefixMatchParam(query string) string {
	return query + "%"
}

func normalizeHashtagQuery(query string) string {
	query = strings.TrimLeft(strings.TrimSpace(query), "#")
	query = strings.Trim(query, "-_")
	return strings.ToLower(query)
}

// func nullableUUID(valid bool, value uuid.UUID) *uuid.UUID {
// 	if !valid {
// 		return nil
// 	}
// 	v := value
// 	return &v
// }
