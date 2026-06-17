package set

import (
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type setScanner interface {
	Scan(dest ...any) error
}

type setRows interface {
	setScanner
	Next() bool
	Err() error
}

func scanSets(rows setRows) ([]*models.Set, error) {
	var sets []*models.Set
	for rows.Next() {
		set, err := scanSet(rows)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	return sets, rows.Err()
}

func scanSet(scanner setScanner) (*models.Set, error) {
	var set models.Set
	err := scanner.Scan(
		&set.ID,
		&set.AuthorUserID,
		&set.PersonaID,
		&set.MediaAssetID,
		&set.Title,
		&set.Description,
		&set.GenreTags,
		&set.DurationSeconds,
		&set.LikeCount,
		&set.CommentCount,
		&set.PlayCount,
		&set.CreatedAt,
		&set.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &set, nil
}

func scanSetComments(rows setRows) ([]*models.SetComment, error) {
	var comments []*models.SetComment
	for rows.Next() {
		comment, err := scanSetComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func scanSetComment(scanner setScanner) (*models.SetComment, error) {
	var comment models.SetComment
	err := scanner.Scan(
		&comment.ID,
		&comment.PersonaID,
		&comment.SetID,
		&comment.Body,
		&comment.ParentID,
		&comment.LikeCount,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func mapSetError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrSetNotFound
	}
	return err
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}
