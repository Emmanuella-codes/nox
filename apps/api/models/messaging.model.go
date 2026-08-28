package models

import (
	"time"

	"github.com/google/uuid"
)

type ConversationType string
type ConversationMemberRole string
type MessageType string

const (
	DirectConversationType ConversationType = "direct"
	GroupConversationType  ConversationType = "group"
)

const (
	ConversationMemberRoleMember ConversationMemberRole = "member"
	ConversationMemberRoleAdmin  ConversationMemberRole = "admin"
)

const (
	TextMessageType   MessageType = "text"
	ImageMessageType  MessageType = "image"
	VideoMessageType  MessageType = "video"
	AudioMessageType  MessageType = "audio"
	SystemMessageType MessageType = "system"
)

type Conversation struct {
	ID               uuid.UUID        `json:"id"`
	ConversationType ConversationType `json:"conversation_type"`
	Title            string           `json:"title"`
	CreatedBy        uuid.UUID        `json:"created_by"`
	LastMessageID    *uuid.UUID       `json:"last_message_id,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type ConversationMember struct {
	ConversationID    uuid.UUID              `json:"conversation_id"`
	UserID            uuid.UUID              `json:"-"`
	PersonaID         uuid.UUID              `json:"persona_id"`
	Role              ConversationMemberRole `json:"role"`
	LastReadMessageID *uuid.UUID             `json:"last_read_message_id,omitempty"`
	JoinedAt          time.Time              `json:"joined_at"`
	LeftAt            *time.Time             `json:"left_at,omitempty"`
}

type Message struct {
	ID              uuid.UUID   `json:"id"`
	ConversationID  uuid.UUID   `json:"conversation_id"`
	SenderUserID    uuid.UUID   `json:"-"`
	SenderPersonaID uuid.UUID   `json:"sender_persona_id"`
	Body            string      `json:"body"`
	MessageType     MessageType `json:"message_type"`
	MediaAssetID    *uuid.UUID  `json:"media_asset_id,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	EditedAt        *time.Time  `json:"edited_at,omitempty"`
	DeletedAt       *time.Time  `json:"deleted_at,omitempty"`
}

type MessageAttachment struct {
	MessageID    uuid.UUID `json:"message_id"`
	MediaAssetID uuid.UUID `json:"media_asset_id"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
}
