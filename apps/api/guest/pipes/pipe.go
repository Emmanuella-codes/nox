package pipes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/emmanuella-codes/nox/guest/dtos"
	"github.com/emmanuella-codes/nox/guest/messages"
	"github.com/emmanuella-codes/nox/models"
	guest_repo "github.com/emmanuella-codes/nox/repositories/guest"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/mail"
	"github.com/google/uuid"
	"strings"
)

type GuestPipe struct {
	repo          guest_repo.GuestRepository
	mail          mail.Provider
	personas      persona_repo.PersonaRepository
	notifications notification_repo.NotificationRepository
	publisher     interface {
		PublishCreatedNotification(context.Context, *models.Notification)
	}
}

func NewGuestPipe(repo guest_repo.GuestRepository, deps ...any) *GuestPipe {
	p := &GuestPipe{repo: repo}
	for _, dep := range deps {
		switch value := dep.(type) {
		case mail.Provider:
			p.mail = value
		case persona_repo.PersonaRepository:
			p.personas = value
		case notification_repo.NotificationRepository:
			p.notifications = value
		case interface {
			PublishCreatedNotification(context.Context, *models.Notification)
		}:
			p.publisher = value
		}
	}
	return p
}

func (p *GuestPipe) Configure(ctx context.Context, userID, eventID uuid.UUID, dto dtos.ConfigureDTO) *shared.PipeRes[models.GuestQuota] {
	if dto.Capacity <= 0 || p.repo.Configure(ctx, userID, eventID, dto) != nil {
		return shared.PipeError[models.GuestQuota](messages.Invalid)
	}
	return shared.PipeSuccess(messages.Created, &models.GuestQuota{EventID: eventID, Capacity: dto.Capacity})
}
func (p *GuestPipe) Quota(ctx context.Context, userID, eventID uuid.UUID, dto dtos.QuotaDTO) *shared.PipeRes[models.GuestQuota] {
	if dto.Capacity <= 0 || strings.TrimSpace(dto.Label) == "" {
		return shared.PipeError[models.GuestQuota](messages.Invalid)
	}
	item, err := p.repo.CreateQuota(ctx, userID, eventID, dto)
	if err != nil {
		return shared.PipeError[models.GuestQuota](messages.Conflict)
	}
	return shared.PipeSuccess(messages.Created, item)
}
func (p *GuestPipe) Add(ctx context.Context, userID, eventID, quotaID uuid.UUID, dto dtos.GuestDTO) *shared.PipeRes[[]*models.GuestEntry] {
	if strings.TrimSpace(dto.Name) == "" || (strings.TrimSpace(dto.Email) == "" && strings.TrimSpace(dto.Phone) == "") {
		return shared.PipeError[[]*models.GuestEntry](messages.Invalid)
	}
	if dto.Quantity <= 0 {
		dto.Quantity = 1
	}
	tokens, hashes := make([]string, dto.Quantity), make([]string, dto.Quantity)
	for i := range tokens {
		tokens[i] = uuid.NewString() + "-" + uuid.NewString()
		hashes[i] = hashToken(tokens[i])
	}
	items, err := p.repo.AddGuest(ctx, userID, eventID, quotaID, dto, hashes)
	if err != nil {
		return shared.PipeError[[]*models.GuestEntry](messages.Conflict)
	}
	for i, item := range items {
		item.InvitationURL = "/guest-invitations/claim?token=" + tokens[i]
		if p.mail != nil && strings.TrimSpace(item.Email) != "" {
			_ = p.mail.Send(ctx, mail.Message{ToEmail: item.Email, ToName: item.Name, Subject: "Your Nox guest pass", TextContent: "Claim your guest pass: " + item.InvitationURL, HTMLContent: "<p>Claim your guest pass: <a href=\"" + item.InvitationURL + "\">open invitation</a></p>"})
		}
		p.notifyMatchedUser(ctx, item)
	}
	return shared.PipeSuccess(messages.Created, &items)
}

func (p *GuestPipe) notifyMatchedUser(ctx context.Context, entry *models.GuestEntry) {
	if entry.MatchedUserID == nil || p.personas == nil || p.notifications == nil {
		return
	}
	personas, err := p.personas.FindPersonasByUserID(ctx, *entry.MatchedUserID)
	if err != nil {
		return
	}
	for _, persona := range personas {
		if persona.PersonaType != models.VisiblePersonaType {
			continue
		}
		created, err := p.notifications.CreateNotifications(ctx, []notification_repo.CreateNotificationInput{{RecipientUserID: *entry.MatchedUserID, RecipientPersonaID: persona.ID, EventID: &entry.EventID, NotificationType: models.GuestInvitationNotificationType}})
		if err == nil && len(created) > 0 && p.publisher != nil {
			p.publisher.PublishCreatedNotification(ctx, created[0])
		}
		return
	}
}
func (p *GuestPipe) Claim(ctx context.Context, token string, userID *uuid.UUID) *shared.PipeRes[models.GuestEntry] {
	if strings.TrimSpace(token) == "" {
		return shared.PipeError[models.GuestEntry](messages.Invalid)
	}
	item, err := p.repo.Claim(ctx, token, userID)
	if err != nil {
		return shared.PipeError[models.GuestEntry](messages.NotFound)
	}
	return shared.PipeSuccess(messages.Claimed, item)
}
func (p *GuestPipe) Decline(ctx context.Context, token string) *shared.PipeRes[any] {
	if strings.TrimSpace(token) == "" {
		return shared.PipeError[any](messages.Invalid)
	}
	if err := p.repo.Decline(ctx, token); err != nil {
		return shared.PipeError[any](messages.NotFound)
	}
	var done any = true
	return shared.PipeSuccess(messages.Declined, &done)
}
func (p *GuestPipe) Revoke(ctx context.Context, userID, entryID uuid.UUID) *shared.PipeRes[any] {
	if err := p.repo.Revoke(ctx, userID, entryID); err != nil {
		return shared.PipeError[any](messages.Forbidden)
	}
	var done any = true
	return shared.PipeSuccess(messages.Revoked, &done)
}
func (p *GuestPipe) Attach(ctx context.Context, userID, entryID uuid.UUID) *shared.PipeRes[any] {
	if err := p.repo.Attach(ctx, userID, entryID); err != nil {
		return shared.PipeError[any](messages.Conflict)
	}
	var done any = true
	return shared.PipeSuccess(messages.Attached, &done)
}
func (p *GuestPipe) CheckIn(ctx context.Context, staffID, eventID uuid.UUID, credential string) *shared.PipeRes[models.GuestEntry] {
	if strings.TrimSpace(credential) == "" {
		return shared.PipeError[models.GuestEntry](messages.Invalid)
	}
	item, err := p.repo.CheckIn(ctx, eventID, staffID, credential)
	if err != nil {
		return shared.PipeError[models.GuestEntry](messages.Conflict)
	}
	return shared.PipeSuccess(messages.CheckedIn, item)
}
func (p *GuestPipe) Report(ctx context.Context, userID, eventID uuid.UUID) *shared.PipeRes[models.GuestListReport] {
	item, err := p.repo.Report(ctx, userID, eventID)
	if err != nil {
		return shared.PipeError[models.GuestListReport](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Reported, item)
}
func hashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
