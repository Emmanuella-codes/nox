package pipes

import (
	"context"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/repositories/persona"
	ticketrepo "github.com/emmanuella-codes/nox/repositories/ticket"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/ticket/dtos"
	"github.com/emmanuella-codes/nox/ticket/messages"
	"github.com/google/uuid"
	"strings"
	"time"
)

type TicketPipe struct {
	repo        ticketrepo.TicketRepository
	personaRepo persona.PersonaRepository
}

func NewTicketPipe(repo ticketrepo.TicketRepository, personas persona.PersonaRepository) *TicketPipe {
	return &TicketPipe{repo: repo, personaRepo: personas}
}
func (p *TicketPipe) CreateWave(ctx context.Context, userID, eventID uuid.UUID, dto dtos.CreateWaveDTO) *shared.PipeRes[models.TicketWave] {
	if strings.TrimSpace(dto.Name) == "" {
		return shared.PipeError[models.TicketWave](messages.Invalid)
	}
	wave, err := p.repo.CreateWave(ctx, eventID, userID, dto)
	if err != nil {
		return shared.PipeError[models.TicketWave](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Created, wave)
}
func (p *TicketPipe) CreateTier(ctx context.Context, userID, waveID uuid.UUID, dto dtos.CreateTierDTO) *shared.PipeRes[models.TicketTier] {
	if strings.TrimSpace(dto.Name) == "" || dto.Inventory <= 0 || dto.AttendeesPerPurchase <= 0 || dto.MaxPurchasesPerUser <= 0 || dto.PriceMinor < 0 {
		return shared.PipeError[models.TicketTier](messages.Invalid)
	}
	tier, err := p.repo.CreateTier(ctx, userID, waveID, dto)
	if err != nil {
		return shared.PipeError[models.TicketTier](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Created, tier)
}
func (p *TicketPipe) Products(ctx context.Context, eventID uuid.UUID) *shared.PipeRes[[]*models.TicketWave] {
	items, err := p.repo.FindProducts(ctx, eventID)
	if err != nil {
		return shared.PipeError[[]*models.TicketWave](messages.Internal)
	}
	return shared.PipeSuccess(messages.Listed, &items)
}
func (p *TicketPipe) ReleaseWave(ctx context.Context, userID, waveID uuid.UUID) *shared.PipeRes[models.TicketWave] {
	item, err := p.repo.ReleaseWave(ctx, userID, waveID)
	if err != nil {
		return shared.PipeError[models.TicketWave](messages.Forbidden)
	}
	return shared.PipeSuccess(messages.Created, item)
}
func (p *TicketPipe) Reserve(ctx context.Context, userID, tierID uuid.UUID, dto dtos.ReserveDTO) *shared.PipeRes[models.TicketReservation] {
	if dto.Quantity <= 0 || strings.TrimSpace(dto.IdempotencyKey) == "" {
		return shared.PipeError[models.TicketReservation](messages.Invalid)
	}
	item, err := p.repo.Reserve(ctx, tierID, userID, dto, 10*time.Minute)
	if err != nil {
		if err == ticketrepo.ErrUnavailable {
			return shared.PipeError[models.TicketReservation](messages.Unavailable)
		}
		return shared.PipeError[models.TicketReservation](messages.Internal)
	}
	return shared.PipeSuccess(messages.Reserved, item)
}
