package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Venue       string    `json:"venue"`
	Location    string    `json:"location"`
	EventDate   time.Time `json:"event_date"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	TicketURL   string    `json:"ticket_url"`
	Price       int       `json:"price_ngn"`
	GenreTags   []string  `json:"genre_tags"`
	OrganizerID uuid.UUID `json:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
}
