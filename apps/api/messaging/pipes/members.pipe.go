package pipes

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *MessagingPipe) AddMembersPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, dto dtos.AddMembersDTO) *shared.PipeRes[[]MemberResponse] {
	admin, message := p.requireMember(ctx, userID, conversationID, dto.AdminPersonaID)
	if message != "" {
		return shared.PipeError[[]MemberResponse](message)
	}
	conversation, err := p.messagingRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrConversationNotFound) {
			return shared.PipeError[[]MemberResponse](messages.Conversation_Not_Found)
		}
		return pipeInternalError[[]MemberResponse](err, "messaging.add_members_conversation")
	}
	if conversation.ConversationType != models.GroupConversationType || admin.Role != models.ConversationMemberRoleAdmin {
		return shared.PipeError[[]MemberResponse](messages.Forbidden)
	}
	personas, pipeMessage := p.visiblePersonas(ctx, userID, dto.MemberPersonaIDs)
	if pipeMessage != "" {
		return shared.PipeError[[]MemberResponse](pipeMessage)
	}
	added, err := p.messagingRepo.AddConversationMembers(ctx, conversationID, personas)
	if err != nil {
		return pipeInternalError[[]MemberResponse](err, "messaging.add_members")
	}
	responses := memberResponses(added)
	return shared.PipeSuccess(messages.Members_Added, &responses)
}

func (p *MessagingPipe) RemoveMemberPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, targetPersonaID uuid.UUID, dto dtos.RemoveMemberDTO) *shared.PipeRes[any] {
	admin, message := p.requireMember(ctx, userID, conversationID, dto.AdminPersonaID)
	if message != "" {
		return shared.PipeError[any](message)
	}
	conversation, err := p.messagingRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrConversationNotFound) {
			return shared.PipeError[any](messages.Conversation_Not_Found)
		}
		return pipeInternalError[any](err, "messaging.remove_member_conversation")
	}
	if conversation.ConversationType != models.GroupConversationType || admin.Role != models.ConversationMemberRoleAdmin || admin.PersonaID == targetPersonaID {
		return shared.PipeError[any](messages.Forbidden)
	}
	if err := p.messagingRepo.RemoveConversationMember(ctx, conversationID, targetPersonaID); err != nil {
		if errors.Is(err, messaging_repo.ErrMembershipNotFound) {
			return shared.PipeError[any](messages.Persona_Not_Found)
		}
		return pipeInternalError[any](err, "messaging.remove_member")
	}
	return shared.PipeSuccess[any](messages.Member_Removed, nil)
}
