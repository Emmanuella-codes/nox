package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Venue            string    `json:"venue"`
	Location         string    `json:"location"`
	EventDate        time.Time `json:"event_date"`
	EventEndDate     time.Time `json:"event_end_date"`
	Status           string    `json:"status"`
	TicketSaleStatus string    `json:"ticket_sale_status"`
	EventModeStatus  string    `json:"event_mode_status"`
	Description      string    `json:"description"`
	CoverURL         string    `json:"cover_url"`
	TicketURL        string    `json:"ticket_url"`
	Price            int       `json:"price_ngn"`
	GenreTags        []string  `json:"genre_tags"`
	OrganizerID      uuid.UUID `json:"organizer_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
