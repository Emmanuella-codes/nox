package pipes

import (
	"context"
	"github.com/emmanuella-codes/nox/models"
	staffrepo "github.com/emmanuella-codes/nox/repositories/staff"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/staff/dtos"
	"github.com/emmanuella-codes/nox/staff/messages"
	"github.com/google/uuid"
	"strings"
)

type StaffPipe struct{ repo staffrepo.StaffRepository }

func NewStaffPipe(repo staffrepo.StaffRepository) *StaffPipe { return &StaffPipe{repo: repo} }
func (p *StaffPipe) Invite(ctx context.Context, actor, eventID uuid.UUID, d dtos.InviteDTO) *shared.PipeRes[models.EventStaffMembership] {
	u, e := uuid.Parse(d.UserID)
	if e != nil {
		return shared.PipeError[models.EventStaffMembership](messages.Invalid)
	}
	item, e := p.repo.Invite(ctx, actor, eventID, u)
	if e != nil {
		return shared.PipeError[models.EventStaffMembership](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Created, item)
}
func (p *StaffPipe) Roles(ctx context.Context, actor, membershipID uuid.UUID, d dtos.RolesDTO) *shared.PipeRes[any] {
	if len(d.Roles) == 0 {
		return shared.PipeError[any](messages.Invalid)
	}
	if e := p.repo.AssignRoles(ctx, actor, membershipID, d.Roles); e != nil {
		return shared.PipeError[any](messages.Forbidden)
	}
	var ok any = true
	return shared.PipeSuccess(messages.Created, &ok)
}
func (p *StaffPipe) Accept(ctx context.Context, userID, membershipID uuid.UUID) *shared.PipeRes[models.EventStaffMembership] {
	item, e := p.repo.Accept(ctx, userID, membershipID)
	if e != nil {
		return shared.PipeError[models.EventStaffMembership](messages.NotFound)
	}
	return shared.PipeSuccess(messages.Accepted, item)
}
func (p *StaffPipe) Revoke(ctx context.Context, actor, membershipID uuid.UUID) *shared.PipeRes[any] {
	if e := p.repo.Revoke(ctx, actor, membershipID); e != nil {
		return shared.PipeError[any](messages.Forbidden)
	}
	var ok any = true
	return shared.PipeSuccess(messages.Revoked, &ok)
}
func (p *StaffPipe) Shift(ctx context.Context, actor, membershipID uuid.UUID, d dtos.ShiftDTO) *shared.PipeRes[models.EventStaffShift] {
	if strings.TrimSpace(d.StartsAt) == "" || strings.TrimSpace(d.EndsAt) == "" {
		return shared.PipeError[models.EventStaffShift](messages.Invalid)
	}
	item, e := p.repo.CreateShift(ctx, actor, membershipID, d)
	if e != nil {
		return shared.PipeError[models.EventStaffShift](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Created, item)
}
func (p *StaffPipe) List(ctx context.Context, actor, eventID uuid.UUID) *shared.PipeRes[[]*models.EventStaffMembership] {
	items, e := p.repo.List(ctx, actor, eventID)
	if e != nil {
		return shared.PipeError[[]*models.EventStaffMembership](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Listed, &items)
}
