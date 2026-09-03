package persona

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

// newPgRepository builds the Postgres-backed profile repository.
func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

// CreatePersona creates the single public profile attached to a user account.
func (r *pgRepository) CreatePersona(ctx context.Context, userID uuid.UUID, dto dtos.CreatePersonaDTO) (*models.Persona, error) {
	existing, err := r.FindPersonasByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, ErrProfileAlreadyExists
	}
	row := r.db.QueryRow(ctx, `
		INSERT INTO personas (
			user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, category, genre_tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, category, genre_tags,
		          follower_count, following_count, post_count, created_at, updated_at
	`, userID, dto.Handle, dto.DisplayName, dto.Bio, dto.AvatarURL, dto.CoverURL, dto.PersonaType, dto.Category, dto.GenreTags)

	persona, err := scanPersona(row)
	if err != nil {
		return nil, mapPersonaError(err)
	}

	return persona, nil
}

// FindPersonaByID fetches one public profile by id.
func (r *pgRepository) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, category, genre_tags,
		       follower_count, following_count, post_count, created_at, updated_at
		FROM personas
		WHERE id = $1
	`, personaID)

	persona, err := scanPersona(row)
	if err != nil {
		return nil, mapPersonaError(err)
	}

	return persona, nil
}

// FindPersonasByUserID fetches profiles belonging to one user account.
func (r *pgRepository) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, category, genre_tags,
		       follower_count, following_count, post_count, created_at, updated_at
		FROM personas
		WHERE user_id = $1
		ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var personas []*models.Persona
	for rows.Next() {
		persona, err := scanPersona(rows)
		if err != nil {
			return nil, err
		}
		personas = append(personas, persona)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return personas, nil
}

// FindPersonaByHandle fetches one public profile by handle.
func (r *pgRepository) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, category, genre_tags,
		       follower_count, following_count, post_count, created_at, updated_at
		FROM personas
		WHERE handle = $1
	`, handle)

	persona, err := scanPersona(row)
	if err != nil {
		return nil, mapPersonaError(err)
	}

	return persona, nil
}

// UpdatePersona updates mutable public profile fields.
func (r *pgRepository) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto dtos.UpdatePersonaDTO) (*models.Persona, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE personas
		SET display_name = $2,
		    bio = $3,
		    avatar_url = $4,
		    cover_url = $5,
		    category = $6,
		    genre_tags = $7,
		    updated_at = now()
		WHERE id = $1
		RETURNING id, user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, category, genre_tags,
		          follower_count, following_count, post_count, created_at, updated_at
	`, personaID, dto.DisplayName, dto.Bio, dto.AvatarURL, dto.CoverURL, dto.Category, dto.GenreTags)

	persona, err := scanPersona(row)
	if err != nil {
		return nil, mapPersonaError(err)
	}

	return persona, nil
}

type personaScanner interface {
	Scan(dest ...any) error
}

// scanPersona maps one public profile row.
func scanPersona(scanner personaScanner) (*models.Persona, error) {
	var persona models.Persona
	err := scanner.Scan(
		&persona.ID,
		&persona.UserID,
		&persona.Handle,
		&persona.DisplayName,
		&persona.Bio,
		&persona.AvatarURL,
		&persona.CoverURL,
		&persona.PersonaType,
		&persona.Category,
		&persona.GenreTags,
		&persona.FollowerCount,
		&persona.FollowingCount,
		&persona.PostCount,
		&persona.CreatedAt,
		&persona.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &persona, nil
}

// mapPersonaError maps repository-specific profile errors.
func mapPersonaError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPersonaNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "personas_handle_key":
			return ErrHandleAlreadyTaken
		case "personas_one_profile_per_user_idx":
			return ErrProfileAlreadyExists
		case "personas_one_ghost_per_user_idx":
			return ErrGhostPersonaAlreadyUsed
		}
	}

	return err
}
