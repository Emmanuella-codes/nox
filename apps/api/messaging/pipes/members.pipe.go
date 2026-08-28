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

// AddMembersPipe adds new members to a group conversation after admin and follow checks.
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
	if message := p.requireMutualFollows(ctx, admin.PersonaID, personas); message != "" {
		return shared.PipeError[[]MemberResponse](message)
	}
	added, err := p.messagingRepo.AddConversationMembers(ctx, conversationID, personas)
	if err != nil {
		return pipeInternalError[[]MemberResponse](err, "messaging.add_members")
	}
	personasByID, err := p.memberPersonas(ctx, added)
	if err != nil {
		return pipeInternalError[[]MemberResponse](err, "messaging.add_member_personas")
	}
	responses := memberResponses(added, personasByID)
	for _, member := range added {
		p.publishMemberEvent(ctx, "conversation.member_added", member)
	}
	return shared.PipeSuccess(messages.Members_Added, &responses)
}

// RemoveMemberPipe removes another member from a group conversation.
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
	p.publishConversationEvent(ctx, conversationID, "conversation.member_removed", map[string]string{
		"conversation_id": conversationID.String(),
		"persona_id":      targetPersonaID.String(),
	})
	return shared.PipeSuccess[any](messages.Member_Removed, nil)
}

// LeaveConversationPipe removes the current member from a group conversation.
func (p *MessagingPipe) LeaveConversationPipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, dto dtos.LeaveConversationDTO) *shared.PipeRes[any] {
	member, message := p.requireMember(ctx, userID, conversationID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[any](message)
	}
	conversation, err := p.messagingRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrConversationNotFound) {
			return shared.PipeError[any](messages.Conversation_Not_Found)
		}
		return pipeInternalError[any](err, "messaging.leave_conversation")
	}
	if conversation.ConversationType != models.GroupConversationType {
		return shared.PipeError[any](messages.Forbidden)
	}
	if err := p.messagingRepo.RemoveConversationMember(ctx, conversationID, member.PersonaID); err != nil {
		if errors.Is(err, messaging_repo.ErrMembershipNotFound) {
			return shared.PipeError[any](messages.Persona_Not_Found)
		}
		return pipeInternalError[any](err, "messaging.leave_member")
	}
	p.publishConversationEvent(ctx, conversationID, "conversation.member_left", map[string]string{
		"conversation_id": conversationID.String(),
		"persona_id":      member.PersonaID.String(),
	})
	return shared.PipeSuccess[any](messages.Member_Left, nil)
}

// UpdateMemberRolePipe promotes or demotes a group member through an admin action.
func (p *MessagingPipe) UpdateMemberRolePipe(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, targetPersonaID uuid.UUID, dto dtos.UpdateMemberRoleDTO) *shared.PipeRes[MemberResponse] {
	admin, message := p.requireMember(ctx, userID, conversationID, dto.AdminPersonaID)
	if message != "" {
		return shared.PipeError[MemberResponse](message)
	}
	conversation, err := p.messagingRepo.FindConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrConversationNotFound) {
			return shared.PipeError[MemberResponse](messages.Conversation_Not_Found)
		}
		return pipeInternalError[MemberResponse](err, "messaging.update_role_conversation")
	}
	if conversation.ConversationType != models.GroupConversationType || admin.Role != models.ConversationMemberRoleAdmin || !validMemberRole(dto.Role) {
		return shared.PipeError[MemberResponse](messages.Forbidden)
	}
	updated, err := p.messagingRepo.UpdateConversationMemberRole(ctx, conversationID, targetPersonaID, dto.Role)
	if err != nil {
		if errors.Is(err, messaging_repo.ErrMembershipNotFound) {
			return shared.PipeError[MemberResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[MemberResponse](err, "messaging.update_member_role")
	}
	personas, err := p.memberPersonas(ctx, []*models.ConversationMember{updated})
	if err != nil {
		return pipeInternalError[MemberResponse](err, "messaging.update_member_role_persona")
	}
	response := memberResponses([]*models.ConversationMember{updated}, personas)[0]
	p.publishMemberEvent(ctx, "conversation.member_role_updated", updated)
	return shared.PipeSuccess(messages.Member_Role_Updated, &response)
}

// validMemberRole validates the supported group member role values.
func validMemberRole(role models.ConversationMemberRole) bool {
	return role == models.ConversationMemberRoleAdmin || role == models.ConversationMemberRoleMember
}
