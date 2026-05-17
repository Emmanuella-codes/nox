package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventDTO struct {
	Title       string    `json:"title" validate:"required"`
	Venue       string    `json:"venue" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	EventDate   time.Time `json:"event_date" validate:"required"`
	Description string    `json:"description" validate:"required"`
	CoverURL    string    `json:"cover_url"`
	TicketURL   string    `json:"ticket_url"`
	Price       int       `json:"price_ngn"`
	GenreTags   []string  `json:"genre_tags"`
	OrganizerID uuid.UUID `json:"organizer_id" validate:"required"`
}
