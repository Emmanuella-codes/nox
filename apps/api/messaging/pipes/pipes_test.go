package pipes

import (
	"context"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/media/dtos"
	messaging_dtos "github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/models"
	persona_dtos "github.com/emmanuella-codes/nox/persona/dtos"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TestCreateDirectConversationPipeAllowsSelfDM verifies that self-DMs are allowed without follow checks.
func TestCreateDirectConversationPipeAllowsSelfDM(t *testing.T) {
	userID, personaID, conversationID := uuid.New(), uuid.New(), uuid.New()
	persona := testPersona(userID, personaID, "self")
	messagingRepo := &messagingTestRepo{
		directConversation: &models.Conversation{ID: conversationID, ConversationType: models.DirectConversationType, CreatedBy: personaID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		membersByConversation: map[uuid.UUID][]*models.ConversationMember{
			conversationID: {testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember)},
		},
	}
	pipe := NewMessagingPipe(messagingRepo, &messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: persona}}, &messagingTestMediaRepo{}, &messagingTestFollowRepo{})

	res := pipe.CreateDirectConversationPipe(context.Background(), userID, messaging_dtos.CreateDirectConversationDTO{
		SenderPersonaID:    personaID,
		RecipientPersonaID: personaID,
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if messagingRepo.createDirectCalls != 0 {
		t.Fatalf("expected no direct create call when conversation exists, got %d", messagingRepo.createDirectCalls)
	}
}

// TestCreateDirectConversationPipeRequiresMutualFollow verifies that two different users must follow each other.
func TestCreateDirectConversationPipeRequiresMutualFollow(t *testing.T) {
	userID, senderID, recipientID := uuid.New(), uuid.New(), uuid.New()
	pipe := NewMessagingPipe(
		&messagingTestRepo{directFindErr: messaging_repo.ErrConversationNotFound},
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{
			senderID:    testPersona(userID, senderID, "sender"),
			recipientID: testPersona(uuid.New(), recipientID, "recipient"),
		}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{following: map[string]bool{followKey(senderID, recipientID): true}},
	)

	res := pipe.CreateDirectConversationPipe(context.Background(), userID, messaging_dtos.CreateDirectConversationDTO{
		SenderPersonaID:    senderID,
		RecipientPersonaID: recipientID,
	})
	if res.Message != messages.Forbidden {
		t.Fatalf("expected %q, got %q", messages.Forbidden, res.Message)
	}
}

// TestCreateGroupConversationPipeRequiresMutualFollow verifies group invites require mutual follows with the creator.
func TestCreateGroupConversationPipeRequiresMutualFollow(t *testing.T) {
	userID, creatorID, memberID := uuid.New(), uuid.New(), uuid.New()
	pipe := NewMessagingPipe(
		&messagingTestRepo{},
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{
			creatorID: testPersona(userID, creatorID, "creator"),
			memberID:  testPersona(uuid.New(), memberID, "member"),
		}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{following: map[string]bool{followKey(creatorID, memberID): true}},
	)

	res := pipe.CreateGroupConversationPipe(context.Background(), userID, messaging_dtos.CreateGroupConversationDTO{
		CreatorPersonaID: creatorID,
		Title:            "Weekend plan",
		MemberPersonaIDs: []uuid.UUID{memberID},
	})
	if res.Message != messages.Forbidden {
		t.Fatalf("expected %q, got %q", messages.Forbidden, res.Message)
	}
}

// TestSendMessagePipeSupportsAttachmentOnlyAudio verifies attachment-only audio messages are accepted.
func TestSendMessagePipeSupportsAttachmentOnlyAudio(t *testing.T) {
	userID, personaID, conversationID, messageID, assetID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	persona := testPersona(userID, personaID, "sender")
	messagingRepo := &messagingTestRepo{
		memberByConversationPersona: map[string]*models.ConversationMember{
			memberKey(conversationID, personaID): testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember),
		},
		createMessageResult: &models.Message{
			ID:              messageID,
			ConversationID:  conversationID,
			SenderUserID:    userID,
			SenderPersonaID: personaID,
			Body:            "",
			MessageType:     models.AudioMessageType,
			CreatedAt:       time.Now(),
		},
		attachmentsByMessage: map[uuid.UUID][]*models.MediaAsset{
			messageID: {testMediaAsset(assetID, userID, personaID, models.AudioMediaKind)},
		},
	}
	pipe := NewMessagingPipe(
		messagingRepo,
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: persona}},
		&messagingTestMediaRepo{assets: map[uuid.UUID]*models.MediaAsset{assetID: testMediaAsset(assetID, userID, personaID, models.AudioMediaKind)}},
		&messagingTestFollowRepo{},
	)

	res := pipe.SendMessagePipe(context.Background(), userID, conversationID, messaging_dtos.SendMessageDTO{
		SenderPersonaID: personaID,
		MediaAssetIDs:   []uuid.UUID{assetID},
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if messagingRepo.createdMessageDTO.MessageType != models.AudioMessageType {
		t.Fatalf("expected audio message type, got %q", messagingRepo.createdMessageDTO.MessageType)
	}
	if len(messagingRepo.createdMessageDTO.MediaAssetIDs) != 1 || messagingRepo.createdMessageDTO.MediaAssetIDs[0] != assetID {
		t.Fatalf("expected one audio attachment, got %#v", messagingRepo.createdMessageDTO.MediaAssetIDs)
	}
	if res.Data == nil || len(res.Data.Attachments) != 1 || res.Data.Attachments[0].ID != assetID {
		t.Fatal("expected response attachments to be hydrated")
	}
}

// TestSendMessagePipeRejectsTooManyAttachments verifies the five-attachment limit.
func TestSendMessagePipeRejectsTooManyAttachments(t *testing.T) {
	pipe := NewMessagingPipe(&messagingTestRepo{}, &messagingTestPersonaRepo{}, &messagingTestMediaRepo{}, &messagingTestFollowRepo{})
	attachments := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()}

	res := pipe.SendMessagePipe(context.Background(), uuid.New(), uuid.New(), messaging_dtos.SendMessageDTO{
		SenderPersonaID: uuid.New(),
		MediaAssetIDs:   attachments,
	})
	if res.Message != messages.Invalid_Payload {
		t.Fatalf("expected %q, got %q", messages.Invalid_Payload, res.Message)
	}
}

// TestEditMessagePipeRejectsExpiredMessage verifies messages cannot be edited after one hour.
func TestEditMessagePipeRejectsExpiredMessage(t *testing.T) {
	userID, messageID := uuid.New(), uuid.New()
	pipe := NewMessagingPipe(
		&messagingTestRepo{
			messages: map[uuid.UUID]*models.Message{
				messageID: {ID: messageID, SenderUserID: userID, CreatedAt: time.Now().Add(-messageMutationWindow - time.Minute)},
			},
		},
		&messagingTestPersonaRepo{},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
	)

	res := pipe.EditMessagePipe(context.Background(), userID, messageID, messaging_dtos.EditMessageDTO{Body: "updated"})
	if res.Message != messages.Forbidden {
		t.Fatalf("expected %q, got %q", messages.Forbidden, res.Message)
	}
}

// TestDeleteMessagePipeDeletesRecentMessage verifies recent sender-owned messages can be deleted for everyone.
func TestDeleteMessagePipeDeletesRecentMessage(t *testing.T) {
	userID, personaID, messageID := uuid.New(), uuid.New(), uuid.New()
	now := time.Now()
	pipe := NewMessagingPipe(
		&messagingTestRepo{
			messages: map[uuid.UUID]*models.Message{
				messageID: {ID: messageID, SenderUserID: userID, SenderPersonaID: personaID, ConversationID: uuid.New(), CreatedAt: now},
			},
			deletedMessageResult: &models.Message{
				ID:              messageID,
				SenderUserID:    userID,
				SenderPersonaID: personaID,
				ConversationID:  uuid.New(),
				CreatedAt:       now,
				DeletedAt:       timePtr(now),
			},
		},
		&messagingTestPersonaRepo{},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
	)

	res := pipe.DeleteMessagePipe(context.Background(), userID, messageID)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data == nil || !res.Data.Deleted {
		t.Fatal("expected deleted response")
	}
}

// TestMarkReadPipeRejectsDeletedMessage verifies deleted messages cannot be used as read cursors.
func TestMarkReadPipeRejectsDeletedMessage(t *testing.T) {
	userID, personaID, conversationID, messageID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	pipe := NewMessagingPipe(
		&messagingTestRepo{
			memberByConversationPersona: map[string]*models.ConversationMember{
				memberKey(conversationID, personaID): testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember),
			},
			messages: map[uuid.UUID]*models.Message{
				messageID: {ID: messageID, ConversationID: conversationID, DeletedAt: timePtr(time.Now())},
			},
		},
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: testPersona(userID, personaID, "reader")}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
	)

	res := pipe.MarkReadPipe(context.Background(), userID, conversationID, messaging_dtos.MarkReadDTO{PersonaID: personaID, MessageID: messageID})
	if res.Message != messages.Message_Not_Found {
		t.Fatalf("expected %q, got %q", messages.Message_Not_Found, res.Message)
	}
}

// TestLeaveConversationPipeLeavesGroup verifies members can leave group conversations.
func TestLeaveConversationPipeLeavesGroup(t *testing.T) {
	userID, personaID, conversationID := uuid.New(), uuid.New(), uuid.New()
	pipe := NewMessagingPipe(
		&messagingTestRepo{
			conversations: map[uuid.UUID]*models.Conversation{
				conversationID: {ID: conversationID, ConversationType: models.GroupConversationType},
			},
			memberByConversationPersona: map[string]*models.ConversationMember{
				memberKey(conversationID, personaID): testMember(conversationID, userID, personaID, models.ConversationMemberRoleMember),
			},
		},
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{personaID: testPersona(userID, personaID, "member")}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
	)

	res := pipe.LeaveConversationPipe(context.Background(), userID, conversationID, messaging_dtos.LeaveConversationDTO{PersonaID: personaID})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
}

// TestUpdateMemberRolePipePromotesAdmin verifies an admin can promote another member.
func TestUpdateMemberRolePipePromotesAdmin(t *testing.T) {
	userID, adminID, memberID, conversationID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	updatedMember := testMember(conversationID, uuid.New(), memberID, models.ConversationMemberRoleAdmin)
	pipe := NewMessagingPipe(
		&messagingTestRepo{
			conversations: map[uuid.UUID]*models.Conversation{
				conversationID: {ID: conversationID, ConversationType: models.GroupConversationType},
			},
			memberByConversationPersona: map[string]*models.ConversationMember{
				memberKey(conversationID, adminID): testMember(conversationID, userID, adminID, models.ConversationMemberRoleAdmin),
			},
			updatedMember: updatedMember,
		},
		&messagingTestPersonaRepo{personas: map[uuid.UUID]*models.Persona{
			adminID:  testPersona(userID, adminID, "admin"),
			memberID: testPersona(uuid.New(), memberID, "member"),
		}},
		&messagingTestMediaRepo{},
		&messagingTestFollowRepo{},
	)

	res := pipe.UpdateMemberRolePipe(context.Background(), userID, conversationID, memberID, messaging_dtos.UpdateMemberRoleDTO{
		AdminPersonaID: adminID,
		Role:           models.ConversationMemberRoleAdmin,
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data == nil || res.Data.Role != string(models.ConversationMemberRoleAdmin) {
		t.Fatal("expected updated admin role")
	}
}

type messagingTestRepo struct {
	conversations               map[uuid.UUID]*models.Conversation
	directConversation          *models.Conversation
	directFindErr               error
	memberByConversationPersona map[string]*models.ConversationMember
	membersByConversation       map[uuid.UUID][]*models.ConversationMember
	createDirectCalls           int
	createGroupCalls            int
	createdMessageDTO           messaging_dtos.SendMessageDTO
	createMessageResult         *models.Message
	createMessageCreated        bool
	messages                    map[uuid.UUID]*models.Message
	deletedMessageResult        *models.Message
	attachmentsByMessage        map[uuid.UUID][]*models.MediaAsset
	updatedMember               *models.ConversationMember
}

// CreateDirectConversation records direct conversation creation calls for assertions.
func (r *messagingTestRepo) CreateDirectConversation(ctx context.Context, creator *models.Persona, recipient *models.Persona) (*models.Conversation, error) {
	r.createDirectCalls++
	if r.directConversation != nil {
		return r.directConversation, nil
	}
	conversation := &models.Conversation{ID: uuid.New(), ConversationType: models.DirectConversationType, CreatedBy: creator.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	r.directConversation = conversation
	return conversation, nil
}

// FindDirectConversationBetweenPersonas returns the configured direct conversation fixture.
func (r *messagingTestRepo) FindDirectConversationBetweenPersonas(ctx context.Context, personaAID uuid.UUID, personaBID uuid.UUID) (*models.Conversation, error) {
	if r.directFindErr != nil {
		return nil, r.directFindErr
	}
	if r.directConversation == nil {
		return nil, messaging_repo.ErrConversationNotFound
	}
	return r.directConversation, nil
}

// CreateGroupConversation records group conversation creation calls for assertions.
func (r *messagingTestRepo) CreateGroupConversation(ctx context.Context, creator *models.Persona, members []*models.Persona, dto messaging_dtos.CreateGroupConversationDTO) (*models.Conversation, error) {
	r.createGroupCalls++
	conversation := &models.Conversation{ID: uuid.New(), ConversationType: models.GroupConversationType, Title: dto.Title, CreatedBy: creator.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if r.conversations == nil {
		r.conversations = map[uuid.UUID]*models.Conversation{}
	}
	r.conversations[conversation.ID] = conversation
	r.membersByConversation = map[uuid.UUID][]*models.ConversationMember{
		conversation.ID: {testMember(conversation.ID, creator.UserID, creator.ID, models.ConversationMemberRoleAdmin)},
	}
	return conversation, nil
}

// FindConversationByID returns the configured conversation fixture.
func (r *messagingTestRepo) FindConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	if conversation, ok := r.conversations[conversationID]; ok {
		return conversation, nil
	}
	return nil, messaging_repo.ErrConversationNotFound
}

// DeleteConversation is unused in these tests.
func (r *messagingTestRepo) DeleteConversation(ctx context.Context, conversationID uuid.UUID) error {
	return nil
}

// ConversationBelongsToInactiveCrew returns false for all test conversations.
func (r *messagingTestRepo) ConversationBelongsToInactiveCrew(ctx context.Context, conversationID uuid.UUID) (bool, error) {
	return false, nil
}

// FindPersonaConversations is unused in these tests.
func (r *messagingTestRepo) FindPersonaConversations(ctx context.Context, userID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*messaging_repo.ConversationListItem, error) {
	return nil, nil
}

// FindConversationMembers returns configured member fixtures for one conversation.
func (r *messagingTestRepo) FindConversationMembers(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationMember, error) {
	return r.membersByConversation[conversationID], nil
}

// FindConversationMemberUserIDs is unused in these tests.
func (r *messagingTestRepo) FindConversationMemberUserIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

// FindRelatedConversationUserIDs is unused in these tests.
func (r *messagingTestRepo) FindRelatedConversationUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

// FindMember returns the configured member fixture for one conversation profile pair.
func (r *messagingTestRepo) FindMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) (*models.ConversationMember, error) {
	member, ok := r.memberByConversationPersona[memberKey(conversationID, personaID)]
	if !ok {
		return nil, messaging_repo.ErrMembershipNotFound
	}
	return member, nil
}

// AddConversationMembers is unused in these tests.
func (r *messagingTestRepo) AddConversationMembers(ctx context.Context, conversationID uuid.UUID, members []*models.Persona) ([]*models.ConversationMember, error) {
	return nil, nil
}

// UpdateConversationMemberRole returns the configured updated member fixture.
func (r *messagingTestRepo) UpdateConversationMemberRole(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, role models.ConversationMemberRole) (*models.ConversationMember, error) {
	if r.updatedMember == nil {
		return nil, messaging_repo.ErrMembershipNotFound
	}
	return r.updatedMember, nil
}

// RemoveConversationMember marks one membership removal call as successful.
func (r *messagingTestRepo) RemoveConversationMember(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID) error {
	if _, ok := r.memberByConversationPersona[memberKey(conversationID, personaID)]; !ok {
		return messaging_repo.ErrMembershipNotFound
	}
	return nil
}

// CreateMessage records the outbound message payload and returns the configured message result.
func (r *messagingTestRepo) CreateMessage(ctx context.Context, conversationID uuid.UUID, senderUserID uuid.UUID, dto messaging_dtos.SendMessageDTO) (*models.Message, bool, error) {
	r.createdMessageDTO = dto
	if r.createMessageResult != nil {
		return r.createMessageResult, r.createMessageCreated, nil
	}
	return &models.Message{ID: uuid.New(), ConversationID: conversationID, SenderUserID: senderUserID, SenderPersonaID: dto.SenderPersonaID, Body: dto.Body, MessageType: dto.MessageType, CreatedAt: time.Now()}, true, nil
}

// UpdateMessageBody returns the updated message fixture for edit flows.
func (r *messagingTestRepo) UpdateMessageBody(ctx context.Context, messageID uuid.UUID, body string) (*models.Message, error) {
	messageModel, ok := r.messages[messageID]
	if !ok {
		return nil, messaging_repo.ErrMessageNotFound
	}
	messageModel.Body = body
	messageModel.EditedAt = timePtr(time.Now())
	return messageModel, nil
}

// FindMessagesByConversationID is unused in these tests.
func (r *messagingTestRepo) FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID, limit int, offset int) ([]*models.Message, error) {
	return nil, nil
}

// FindMessageByID returns the configured message fixture for one id.
func (r *messagingTestRepo) FindMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	messageModel, ok := r.messages[messageID]
	if !ok {
		return nil, messaging_repo.ErrMessageNotFound
	}
	return messageModel, nil
}

// FindMessageAttachmentsByMessageIDs returns configured attachment fixtures for message responses.
func (r *messagingTestRepo) FindMessageAttachmentsByMessageIDs(ctx context.Context, messageIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error) {
	return r.attachmentsByMessage, nil
}

// MarkConversationRead returns the member fixture for read-tracking assertions.
func (r *messagingTestRepo) MarkConversationRead(ctx context.Context, conversationID uuid.UUID, personaID uuid.UUID, messageID uuid.UUID) (*models.ConversationMember, error) {
	return r.memberByConversationPersona[memberKey(conversationID, personaID)], nil
}

// SoftDeleteMessage returns the configured deleted message fixture.
func (r *messagingTestRepo) SoftDeleteMessage(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	if r.deletedMessageResult == nil {
		return nil, messaging_repo.ErrMessageNotFound
	}
	return r.deletedMessageResult, nil
}

type messagingTestPersonaRepo struct {
	personas map[uuid.UUID]*models.Persona
}

// CreatePersona is unused in these tests.
func (r *messagingTestPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto persona_dtos.CreatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

// FindPersonaByID returns the configured profile fixture for one id.
func (r *messagingTestPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona, ok := r.personas[personaID]
	if !ok {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

// FindPersonasByUserID is unused in these tests.
func (r *messagingTestPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

// FindPersonaByHandle is unused in these tests.
func (r *messagingTestPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

// UpdatePersona is unused in these tests.
func (r *messagingTestPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto persona_dtos.UpdatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

type messagingTestMediaRepo struct {
	assets map[uuid.UUID]*models.MediaAsset
}

// CreateMediaAsset is unused in these tests.
func (r *messagingTestMediaRepo) CreateMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto dtos.CreateMediaAssetDTO) (*models.MediaAsset, error) {
	return nil, nil
}

// CreatePostMediaAsset is unused in these tests.
func (r *messagingTestMediaRepo) CreatePostMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto dtos.ConfirmPostMediaUploadDTO) (*models.MediaAsset, error) {
	return nil, nil
}

// CreatePendingMediaAsset is unused in these tests.
func (r *messagingTestMediaRepo) CreatePendingMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateSetVideoUploadDTO) (*models.MediaAsset, error) {
	return nil, nil
}

// CreatePendingStoryMediaAsset is unused in these tests.
func (r *messagingTestMediaRepo) CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateStoryVideoUploadDTO) (*models.MediaAsset, error) {
	return nil, nil
}

// FindMediaAssetByID returns the configured attachment fixture for one id.
func (r *messagingTestMediaRepo) FindMediaAssetByID(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	asset, ok := r.assets[mediaAssetID]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	return asset, nil
}

// MarkMediaAssetReady is unused in these tests.
func (r *messagingTestMediaRepo) MarkMediaAssetReady(ctx context.Context, mediaAssetID uuid.UUID, dto dtos.CompleteMediaProcessingDTO) (*models.MediaAsset, error) {
	return nil, nil
}

// MarkMediaAssetFailed is unused in these tests.
func (r *messagingTestMediaRepo) MarkMediaAssetFailed(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	return nil, nil
}

// DeleteOrphanedMediaAssets is unused in these tests.
func (r *messagingTestMediaRepo) DeleteOrphanedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	return 0, nil
}

type messagingTestFollowRepo struct {
	following map[string]bool
}

// Follow is unused in these tests.
func (r *messagingTestFollowRepo) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return nil
}

// Unfollow is unused in these tests.
func (r *messagingTestFollowRepo) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return nil
}

// IsFollowing returns the configured follow relationship for one pair.
func (r *messagingTestFollowRepo) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return r.following[followKey(followerID, followingID)], nil
}

// FindFollowingIDs is unused in these tests.
func (r *messagingTestFollowRepo) FindFollowingIDs(ctx context.Context, followerID uuid.UUID, followingIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	return nil, nil
}

// FindFollowers is unused in these tests.
func (r *messagingTestFollowRepo) FindFollowers(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) ([]*models.Persona, error) {
	return nil, nil
}

// FindFollowing is unused in these tests.
func (r *messagingTestFollowRepo) FindFollowing(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) ([]*models.Persona, error) {
	return nil, nil
}

// testPersona builds one public profile fixture.
func testPersona(userID uuid.UUID, personaID uuid.UUID, handle string) *models.Persona {
	return &models.Persona{ID: personaID, UserID: userID, Handle: handle, DisplayName: handle, PersonaType: models.VisiblePersonaType}
}

// testMember builds one active conversation member fixture.
func testMember(conversationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, role models.ConversationMemberRole) *models.ConversationMember {
	return &models.ConversationMember{ConversationID: conversationID, UserID: userID, PersonaID: personaID, Role: role, JoinedAt: time.Now()}
}

// testMediaAsset builds one ready attachment fixture.
func testMediaAsset(assetID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, kind models.MediaKind) *models.MediaAsset {
	return &models.MediaAsset{ID: assetID, OwnerUserID: userID, OwnerPersonaID: personaID, MediaKind: kind, ProcessingStatus: models.ReadyMediaStatus}
}

// memberKey builds one stable member lookup key.
func memberKey(conversationID uuid.UUID, personaID uuid.UUID) string {
	return conversationID.String() + ":" + personaID.String()
}

// followKey builds one stable follow lookup key.
func followKey(followerID uuid.UUID, followingID uuid.UUID) string {
	return followerID.String() + ":" + followingID.String()
}

// timePtr returns a pointer to the supplied timestamp.
func timePtr(value time.Time) *time.Time {
	return &value
}
