package user

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateUser(ctx context.Context, fullname string, email string, passwordHash string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (fullname, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, fullname, email, password, created_at, updated_at
	`, fullname, email, passwordHash).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}

func (r *pgRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, fullname, email, password, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Fullname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
