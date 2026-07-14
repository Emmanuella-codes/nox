export interface Event {
  id: string;
  title: string;
  venue: string;
  location: string;
  event_date: string;
  description: string;
  cover_url: string;
  ticket_url: string;
  price_ngn: number;
  genre_tags: string[];
  organizer_id: string;
  created_at: string;
}

export interface CreateEventRequest {
  title: string;
  venue: string;
  location: string;
  event_date: string;
  description: string;
  cover_url?: string;
  ticket_url?: string;
  price_ngn: number;
  genre_tags: string[];
  organizer_id: string;
}
