package pipes

import (
	"context"
	"errors"
	"strings"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// CreateDirectConversationPipe creates or reuses a direct conversation between two public profiles.
func (p *MessagingPipe) CreateDirectConversationPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateDirectConversationDTO) *shared.PipeRes[ConversationResponse] {
	sender, message := p.profilePersona(ctx, userID, dto.SenderPersonaID, true)
	if message != "" {
		return shared.PipeError[ConversationResponse](message)
	}
	recipient, message := p.profilePersona(ctx, userID, dto.RecipientPersonaID, false)
	if message != "" {
		return shared.PipeError[ConversationResponse](message)
	}
	if sender.ID != recipient.ID {
		if message := p.requireMutualFollow(ctx, sender.ID, recipient.ID); message != "" {
			return shared.PipeError[ConversationResponse](message)
		}
	}

	conversation, err := p.messagingRepo.FindDirectConversationBetweenPersonas(ctx, sender.ID, recipient.ID)
	if err != nil && !errors.Is(err, messaging_repo.ErrConversationNotFound) {
		return pipeInternalError[ConversationResponse](err, "messaging.find_direct")
	}
	if conversation == nil {
		conversation, err = p.messagingRepo.CreateDirectConversation(ctx, sender, recipient)
		if err != nil {
			conversation, err = p.messagingRepo.FindDirectConversationBetweenPersonas(ctx, sender.ID, recipient.ID)
			if err != nil {
				return pipeInternalError[ConversationResponse](err, "messaging.create_direct")
			}
		}
	}

	members, err := p.messagingRepo.FindConversationMembers(ctx, conversation.ID)
	if err != nil {
		return pipeInternalError[ConversationResponse](err, "messaging.direct_members")
	}
	personas, err := p.memberPersonas(ctx, members)
	if err != nil {
		return pipeInternalError[ConversationResponse](err, "messaging.direct_member_personas")
	}
	response := p.conversationResponse(ctx, conversation, members, personas, nil, 0)
	p.publishConversationEvent(ctx, conversation.ID, "conversation.created", response)
	return shared.PipeSuccess(messages.Conversation_Created, &response)
}

// CreateGroupConversationPipe creates a group conversation for one owner profile and invited profiles.
func (p *MessagingPipe) CreateGroupConversationPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateGroupConversationDTO) *shared.PipeRes[ConversationResponse] {
	dto.Title = strings.TrimSpace(dto.Title)
	if dto.Title == "" {
		return shared.PipeError[ConversationResponse](messages.Invalid_Payload)
	}
	creator, message := p.profilePersona(ctx, userID, dto.CreatorPersonaID, true)
	if message != "" {
		return shared.PipeError[ConversationResponse](message)
	}
	memberPersonas, message := p.visiblePersonas(ctx, userID, dto.MemberPersonaIDs)
	if message != "" {
		return shared.PipeError[ConversationResponse](message)
	}
	if message := p.requireMutualFollows(ctx, creator.ID, memberPersonas); message != "" {
		return shared.PipeError[ConversationResponse](message)
	}
	hasOtherMember := false
	for _, member := range memberPersonas {
		if member.ID != creator.ID {
			hasOtherMember = true
			break
		}
	}
	if !hasOtherMember {
		return shared.PipeError[ConversationResponse](messages.Invalid_Payload)
	}

	conversation, err := p.messagingRepo.CreateGroupConversation(ctx, creator, memberPersonas, dto)
	if err != nil {
		return pipeInternalError[ConversationResponse](err, "messaging.create_group")
	}
	members, err := p.messagingRepo.FindConversationMembers(ctx, conversation.ID)
	if err != nil {
		return pipeInternalError[ConversationResponse](err, "messaging.group_members")
	}
	personas, err := p.memberPersonas(ctx, members)
	if err != nil {
		return pipeInternalError[ConversationResponse](err, "messaging.group_member_personas")
	}
	response := p.conversationResponse(ctx, conversation, members, personas, nil, 0)
	p.publishConversationEvent(ctx, conversation.ID, "conversation.created", response)
	return shared.PipeSuccess(messages.Conversation_Created, &response)
}

// visiblePersonas fetches and deduplicates profile participants for a group conversation.
func (p *MessagingPipe) visiblePersonas(ctx context.Context, userID uuid.UUID, personaIDs []uuid.UUID) ([]*models.Persona, shared.PipeMessage) {
	seen := map[uuid.UUID]bool{}
	personas := make([]*models.Persona, 0, len(personaIDs))
	for _, personaID := range personaIDs {
		if personaID == uuid.Nil || seen[personaID] {
			continue
		}
		persona, message := p.profilePersona(ctx, userID, personaID, false)
		if message != "" {
			return nil, message
		}
		seen[personaID] = true
		personas = append(personas, persona)
	}
	return personas, ""
}
