package pipes

import (
	"context"
	"github.com/emmanuella-codes/nox/admission/dtos"
	"github.com/emmanuella-codes/nox/admission/messages"
	"github.com/emmanuella-codes/nox/models"
	repo "github.com/emmanuella-codes/nox/repositories/admission"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"strings"
)

type AdmissionPipe struct{ repo repo.AdmissionRepository }

func NewAdmissionPipe(r repo.AdmissionRepository) *AdmissionPipe { return &AdmissionPipe{repo: r} }

func (p *AdmissionPipe) Create(ctx context.Context, actor, event uuid.UUID, d dtos.CreateAdmissionDTO) *shared.PipeRes[models.Admission] {
	sourceID, userID := uuid.Nil, uuid.Nil
	if d.SourceID != "" {
		id, err := uuid.Parse(d.SourceID)
		if err != nil {
			return shared.PipeError[models.Admission](messages.Invalid)
		}
		sourceID = id
	}
	if d.UserID != "" {
		id, err := uuid.Parse(d.UserID)
		if err != nil {
			return shared.PipeError[models.Admission](messages.Invalid)
		}
		userID = id
	}
	item, err := p.repo.Create(ctx, actor, event, d.SourceType, sourceID, userID, d.Name, d.Email, d.Phone, "")
	if err != nil {
		return shared.PipeError[models.Admission](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Created, item)
}
func (p *AdmissionPipe) QR(ctx context.Context, user, event uuid.UUID, credential string) *shared.PipeRes[models.CheckInResult] {
	if strings.TrimSpace(credential) == "" {
		return shared.PipeError[models.CheckInResult](messages.Invalid)
	}
	item, err := p.repo.CheckIn(ctx, user, event, credential, "qr")
	if err != nil {
		return checkInError[models.CheckInResult](err)
	}
	return shared.PipeSuccess(messages.CheckedIn, item)
}
func (p *AdmissionPipe) Manual(ctx context.Context, user, event uuid.UUID, d dtos.ManualCheckInDTO) *shared.PipeRes[models.CheckInResult] {
	id, err := uuid.Parse(d.AdmissionID)
	if err != nil {
		return shared.PipeError[models.CheckInResult](messages.Invalid)
	}
	item, err := p.repo.ManualCheckIn(ctx, user, event, id)
	if err != nil {
		return checkInError[models.CheckInResult](err)
	}
	return shared.PipeSuccess(messages.CheckedIn, item)
}
func (p *AdmissionPipe) Override(ctx context.Context, user, admission uuid.UUID, d dtos.OverrideDTO) *shared.PipeRes[models.CheckInResult] {
	if strings.TrimSpace(d.Reason) == "" {
		return shared.PipeError[models.CheckInResult](messages.Invalid)
	}
	item, err := p.repo.OverrideCheckIn(ctx, user, admission, d.Reason)
	if err != nil {
		return checkInError[models.CheckInResult](err)
	}
	return shared.PipeSuccess(messages.CheckedIn, item)
}
func (p *AdmissionPipe) Search(ctx context.Context, user, event uuid.UUID, term string) *shared.PipeRes[[]*models.AdmissionSearchResult] {
	if strings.TrimSpace(term) == "" {
		return shared.PipeError[[]*models.AdmissionSearchResult](messages.Invalid)
	}
	items, err := p.repo.Search(ctx, user, event, term)
	if err != nil {
		return shared.PipeError[[]*models.AdmissionSearchResult](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Listed, &items)
}
func (p *AdmissionPipe) Reconcile(ctx context.Context, user, event uuid.UUID) *shared.PipeRes[map[string]int] {
	item, err := p.repo.Reconcile(ctx, user, event)
	if err != nil {
		return shared.PipeError[map[string]int](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Reconciled, &item)
}
func checkInError[T any](err error) *shared.PipeRes[T] {
	if invalid, ok := err.(*repo.InvalidError); ok {
		return &shared.PipeRes[T]{Success: false, Message: messages.InvalidAdmission, HookData: map[string]string{"reason": invalid.Reason}}
	}
	if err == repo.ErrAlreadyCheckedIn {
		return shared.PipeError[T](messages.AlreadyCheckedIn)
	}
	return shared.PipeError[T](messages.Forbidden)
}
