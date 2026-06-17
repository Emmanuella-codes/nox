package pipes

import (
	"time"

	"github.com/emmanuella-codes/nox/models"
)

type SetResponse struct {
	ID              string             `json:"id"`
	PersonaID       string             `json:"persona_id"`
	MediaAssetID    string             `json:"media_asset_id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	GenreTags       []string           `json:"genre_tags"`
	DurationSeconds int                `json:"duration_seconds"`
	LikeCount       int                `json:"like_count"`
	CommentCount    int                `json:"comment_count"`
	PlayCount       int                `json:"play_count"`
	IsLiked         bool               `json:"is_liked"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Persona         *models.Persona    `json:"persona,omitempty"`
	MediaAsset      *models.MediaAsset `json:"media_asset,omitempty"`
}

type SetListResponse struct {
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
	HasMore    bool          `json:"has_more"`
	NextOffset *int          `json:"next_offset,omitempty"`
	Sets       []SetResponse `json:"sets"`
}

type SetCommentResponse struct {
	ID        string          `json:"id"`
	PersonaID string          `json:"persona_id"`
	SetID     string          `json:"set_id"`
	Body      string          `json:"body"`
	ParentID  string          `json:"parent_id,omitempty"`
	LikeCount int             `json:"like_count"`
	CreatedAt time.Time       `json:"created_at"`
	Persona   *models.Persona `json:"persona,omitempty"`
}

type SetCommentListResponse struct {
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
	HasMore    bool                 `json:"has_more"`
	NextOffset *int                 `json:"next_offset,omitempty"`
	Comments   []SetCommentResponse `json:"comments"`
}

func setResponse(set *models.Set, liked bool) SetResponse {
	return SetResponse{
		ID:              set.ID.String(),
		PersonaID:       set.PersonaID.String(),
		MediaAssetID:    set.MediaAssetID.String(),
		Title:           set.Title,
		Description:     set.Description,
		GenreTags:       set.GenreTags,
		DurationSeconds: set.DurationSeconds,
		LikeCount:       set.LikeCount,
		CommentCount:    set.CommentCount,
		PlayCount:       set.PlayCount,
		IsLiked:         liked,
		CreatedAt:       set.CreatedAt,
		UpdatedAt:       set.UpdatedAt,
		Persona:         set.Persona,
		MediaAsset:      set.MediaAsset,
	}
}

func setCommentResponse(comment *models.SetComment) SetCommentResponse {
	response := SetCommentResponse{
		ID:        comment.ID.String(),
		PersonaID: comment.PersonaID.String(),
		SetID:     comment.SetID.String(),
		Body:      comment.Body,
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt,
		Persona:   comment.Persona,
	}
	if comment.ParentID.String() != "00000000-0000-0000-0000-000000000000" {
		response.ParentID = comment.ParentID.String()
	}
	return response
}
