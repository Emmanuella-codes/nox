package search

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
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
}

type Repository interface {
	Search(ctx context.Context, query string, limit int) (*Results, error)
}

func NewSearchRepository(db *pgxpool.Pool) Repository {
	return newPgRepository(db)
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

func nullableUUID(valid bool, value uuid.UUID) *uuid.UUID {
	if !valid {
		return nil
	}
	v := value
	return &v
}
