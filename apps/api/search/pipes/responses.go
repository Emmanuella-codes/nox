package pipes

import (
	"time"

	"github.com/emmanuella-codes/nox/models"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
)

type SearchResponse struct {
	Query      string                  `json:"query"`
	Limit      int                     `json:"limit"`
	Offset     int                     `json:"offset"`
	HasMore    bool                    `json:"has_more"`
	NextOffset *int                    `json:"next_offset,omitempty"`
	Personas   []SearchPersonaResponse `json:"personas"`
	Posts      []SearchPostResponse    `json:"posts"`
	Events     []SearchEventResponse   `json:"events"`
	Hashtags   []SearchHashtagResponse `json:"hashtags"`
}

type SearchPersonaResponse struct {
	ID             string    `json:"id"`
	Handle         string    `json:"handle"`
	DisplayName    string    `json:"display_name"`
	Bio            string    `json:"bio"`
	AvatarURL      string    `json:"avatar_url"`
	CoverURL       string    `json:"cover_url"`
	PersonaType    string    `json:"persona_type"`
	Category       string    `json:"category"`
	GenreTags      []string  `json:"genre_tags"`
	FollowerCount  int       `json:"follower_count"`
	FollowingCount int       `json:"following_count"`
	IsFollowing    bool      `json:"is_following"`
	PostCount      int       `json:"post_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SearchPostResponse struct {
	ID           string           `json:"id"`
	Author       SearchPostAuthor `json:"author"`
	EventID      *string          `json:"event_id,omitempty"`
	Body         string           `json:"body"`
	PostType     models.PostType  `json:"post_type"`
	MediaURL     string           `json:"media_url,omitempty"`
	MediaType    models.MediaType `json:"media_type,omitempty"`
	Location     string           `json:"location,omitempty"`
	LikeCount    int              `json:"like_count"`
	CommentCount int              `json:"comment_count"`
	RepostCount  int              `json:"repost_count"`
	IsLiked      bool             `json:"is_liked"`
	IsRepost     bool             `json:"is_repost"`
	RepostOf     *string          `json:"repost_of,omitempty"`
	Hashtags     []string         `json:"hashtags"`
	CreatedAt    time.Time        `json:"created_at"`
}

type SearchPostAuthor struct {
	Mode      models.PostingMode   `json:"mode"`
	Persona   *SearchPostPersona   `json:"persona,omitempty"`
	Anonymous *SearchPostAnonymous `json:"anonymous,omitempty"`
}

type SearchPostPersona struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type SearchPostAnonymous struct {
	Handle string `json:"handle"`
}

type SearchEventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Venue       string    `json:"venue"`
	Location    string    `json:"location"`
	EventDate   time.Time `json:"event_date"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	TicketURL   string    `json:"ticket_url"`
	Price       int       `json:"price_ngn"`
	GenreTags   []string  `json:"genre_tags"`
	OrganizerID string    `json:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type SearchHashtagResponse struct {
	ID        string    `json:"id"`
	Tag       string    `json:"tag"`
	PostCount int       `json:"post_count"`
	CreatedAt time.Time `json:"created_at"`
}

func personaResponses(personas []*models.Persona) []SearchPersonaResponse {
	responses := make([]SearchPersonaResponse, 0, len(personas))
	for _, persona := range personas {
		responses = append(responses, SearchPersonaResponse{
			ID:             persona.ID.String(),
			Handle:         persona.Handle,
			DisplayName:    persona.DisplayName,
			Bio:            persona.Bio,
			AvatarURL:      persona.AvatarURL,
			CoverURL:       persona.CoverURL,
			PersonaType:    string(persona.PersonaType),
			Category:       string(persona.Category),
			GenreTags:      persona.GenreTags,
			FollowerCount:  persona.FollowerCount,
			FollowingCount: persona.FollowingCount,
			IsFollowing:    false,
			PostCount:      persona.PostCount,
			CreatedAt:      persona.CreatedAt,
			UpdatedAt:      persona.UpdatedAt,
		})
	}
	return responses
}

func postResponses(posts []*searchrepo.PostResult) []SearchPostResponse {
	responses := make([]SearchPostResponse, 0, len(posts))
	for _, result := range posts {
		post := result.Post
		response := SearchPostResponse{
			ID:           post.ID.String(),
			Author:       postAuthor(result),
			Body:         post.Body,
			PostType:     post.PostType,
			MediaURL:     post.MediaURL,
			MediaType:    post.MediaType,
			Location:     post.Location,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			RepostCount:  post.RepostCount,
			IsLiked:      false,
			IsRepost:     post.IsRepost,
			Hashtags:     []string{},
			CreatedAt:    post.CreatedAt,
		}
		if post.EventID != nil {
			eventID := post.EventID.String()
			response.EventID = &eventID
		}
		if post.RepostOf != nil {
			repostOf := post.RepostOf.String()
			response.RepostOf = &repostOf
		}
		responses = append(responses, response)
	}
	return responses
}

func postAuthor(result *searchrepo.PostResult) SearchPostAuthor {
	author := SearchPostAuthor{Mode: result.Post.PostingMode}
	if result.Post.PostingMode == models.AnonymousPostingMode {
		author.Anonymous = &SearchPostAnonymous{Handle: "anonymous"}
		return author
	}
	if result.Post.PostingMode == models.PublicPostingMode && result.Persona != nil {
		author.Persona = &SearchPostPersona{
			ID:          result.Persona.ID.String(),
			Handle:      result.Persona.Handle,
			DisplayName: result.Persona.DisplayName,
			AvatarURL:   result.Persona.AvatarURL,
		}
	}
	return author
}

func eventResponses(events []*models.Event) []SearchEventResponse {
	responses := make([]SearchEventResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, SearchEventResponse{
			ID:          event.ID.String(),
			Title:       event.Title,
			Venue:       event.Venue,
			Location:    event.Location,
			EventDate:   event.EventDate,
			Description: event.Description,
			CoverURL:    event.CoverURL,
			TicketURL:   event.TicketURL,
			Price:       event.Price,
			GenreTags:   event.GenreTags,
			OrganizerID: event.OrganizerID.String(),
			CreatedAt:   event.CreatedAt,
		})
	}
	return responses
}

func hashtagResponses(hashtags []*models.Hashtag) []SearchHashtagResponse {
	responses := make([]SearchHashtagResponse, 0, len(hashtags))
	for _, hashtag := range hashtags {
		responses = append(responses, SearchHashtagResponse{
			ID:        hashtag.ID.String(),
			Tag:       hashtag.Tag,
			PostCount: hashtag.PostCount,
			CreatedAt: hashtag.CreatedAt,
		})
	}
	return responses
}
