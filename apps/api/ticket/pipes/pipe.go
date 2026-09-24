package pipes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/repositories/persona"
	ticketrepo "github.com/emmanuella-codes/nox/repositories/ticket"
	userrepo "github.com/emmanuella-codes/nox/repositories/user"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/payments"
	"github.com/emmanuella-codes/nox/ticket/dtos"
	"github.com/emmanuella-codes/nox/ticket/messages"
	"github.com/google/uuid"
	"strings"
	"time"
)

type TicketPipe struct {
	repo        ticketrepo.TicketRepository
	personaRepo persona.PersonaRepository
	userRepo    userrepo.UserRepository
	payment     payments.Provider
}

func NewTicketPipe(repo ticketrepo.TicketRepository, personas persona.PersonaRepository, deps ...any) *TicketPipe {
	p := &TicketPipe{repo: repo, personaRepo: personas}
	for _, dep := range deps {
		switch value := dep.(type) {
		case userrepo.UserRepository:
			p.userRepo = value
		case payments.Provider:
			p.payment = value
		}
	}
	return p
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
	order, err := p.repo.CreateOrder(ctx, item.ID, userID)
	if err != nil {
		return shared.PipeError[models.TicketReservation](messages.Internal)
	}
	item.OrderID = order.ID
	return shared.PipeSuccess(messages.Reserved, item)
}

func (p *TicketPipe) Order(ctx context.Context, userID, orderID uuid.UUID) *shared.PipeRes[models.TicketOrder] {
	order, err := p.repo.FindOrder(ctx, orderID, userID)
	if err != nil {
		return shared.PipeError[models.TicketOrder](messages.NotFound)
	}
	return shared.PipeSuccess(messages.OrderReady, order)
}

func (p *TicketPipe) InitializePayment(ctx context.Context, userID, orderID uuid.UUID) *shared.PipeRes[models.TicketOrder] {
	if p.payment == nil || p.userRepo == nil {
		return shared.PipeError[models.TicketOrder](messages.Internal)
	}
	order, err := p.repo.FindOrder(ctx, orderID, userID)
	if err != nil || order.Status != "pending" {
		return shared.PipeError[models.TicketOrder](messages.NotFound)
	}
	user, err := p.userRepo.FindUserByID(ctx, userID.String())
	if err != nil {
		return shared.PipeError[models.TicketOrder](messages.Internal)
	}
	result, err := p.payment.Initialize(ctx, payments.InitRequest{Email: user.Email, AmountMinor: order.SubtotalMinor, Currency: order.Currency, Reference: order.ProviderReference})
	if err != nil {
		return shared.PipeError[models.TicketOrder](messages.Internal)
	}
	order.AuthorizationURL = result.AuthorizationURL
	return shared.PipeSuccess(messages.PaymentReady, order)
}

func (p *TicketPipe) HandleWebhook(ctx context.Context, body []byte) error {
	var input struct {
		Event string `json:"event"`
		Data  struct {
			ID        string `json:"id"`
			Reference string `json:"reference"`
			Amount    int64  `json:"amount"`
			Currency  string `json:"currency"`
			Status    string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &input); err != nil || input.Event != "charge.success" || input.Data.Status != "success" || input.Data.Reference == "" {
		return ticketrepo.ErrUnavailable
	}
	eventID := input.Data.ID
	if eventID == "" {
		hash := sha256.Sum256(body)
		eventID = hex.EncodeToString(hash[:])
	}
	return p.repo.SettlePayment(ctx, input.Data.Reference, eventID, input.Data.Amount, input.Data.Currency, body)
}

func (p *TicketPipe) VerifyWebhook(body []byte, signature string) bool {
	return p.payment != nil && p.payment.VerifyWebhook(body, signature)
}
