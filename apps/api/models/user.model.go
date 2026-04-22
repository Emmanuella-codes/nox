package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID    `json:"id"`
	Fullname        string       `json:"fullname"`
	Email           string       `json:"email"`
	Password        string       `json:"-"`
	EmailVerified   bool         `json:"email_verified"`
	EmailVerifiedAt sql.NullTime `json:"-"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}
