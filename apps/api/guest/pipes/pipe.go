package pipes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/emmanuella-codes/nox/guest/dtos"
	"github.com/emmanuella-codes/nox/guest/messages"
	"github.com/emmanuella-codes/nox/models"
	guest_repo "github.com/emmanuella-codes/nox/repositories/guest"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/mail"
	"github.com/google/uuid"
	"strings"
)

type GuestPipe struct {
	repo guest_repo.GuestRepository
	mail mail.Provider
}

func NewGuestPipe(repo guest_repo.GuestRepository, providers ...mail.Provider) *GuestPipe {
	p := &GuestPipe{repo: repo}
	if len(providers) > 0 {
		p.mail = providers[0]
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
	}
	return shared.PipeSuccess(messages.Created, &items)
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
func hashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
