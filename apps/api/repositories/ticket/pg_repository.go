package ticket

import (
	"context"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/ticket/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type pgRepository struct{ db *pgxpool.Pool }

func (r *pgRepository) CreateWave(ctx context.Context, eventID, userID uuid.UUID, dto dtos.CreateWaveDTO) (*models.TicketWave, error) {
	start, end, err := parseTimes(dto.SaleStartAt, dto.SaleEndAt)
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, `INSERT INTO ticket_waves(event_id,name,display_order,sale_start_at,sale_end_at,release_mode,status) SELECT $1,$3,$4,$5,$6,$7,CASE WHEN $5 IS NULL THEN 'not_released' ELSE 'scheduled' END FROM events WHERE id=$1 AND organizer_id=$2 AND status='published' RETURNING id,event_id,name,display_order,sale_start_at,sale_end_at,release_mode,status`, eventID, userID, dto.Name, dto.DisplayOrder, start, end, defaultRelease(dto.ReleaseMode))
	return scanWave(row)
}

func (r *pgRepository) CreateTier(ctx context.Context, userID, waveID uuid.UUID, dto dtos.CreateTierDTO) (*models.TicketTier, error) {
	if dto.Currency == "" {
		dto.Currency = "NGN"
	}
	row := r.db.QueryRow(ctx, `INSERT INTO ticket_tiers(wave_id,name,description,price_minor,currency,inventory,remaining_inventory,attendees_per_purchase,max_purchases_per_user) SELECT w.id,$3,$4,$5,$6,$7,$7,$8,$9 FROM ticket_waves w JOIN events e ON e.id=w.event_id WHERE w.id=$1 AND e.organizer_id=$2 AND w.status IN ('not_released','scheduled') RETURNING id,wave_id,name,description,price_minor,currency,inventory,remaining_inventory,attendees_per_purchase,max_purchases_per_user,status`, waveID, userID, dto.Name, dto.Description, dto.PriceMinor, dto.Currency, dto.Inventory, dto.AttendeesPerPurchase, dto.MaxPurchasesPerUser)
	return scanTier(row)
}

func (r *pgRepository) FindProducts(ctx context.Context, eventID uuid.UUID) ([]*models.TicketWave, error) {
	rows, err := r.db.Query(ctx, `SELECT w.id,w.event_id,w.name,w.display_order,w.sale_start_at,w.sale_end_at,w.release_mode,w.status,t.id,t.wave_id,t.name,t.description,t.price_minor,t.currency,t.inventory,t.remaining_inventory,t.attendees_per_purchase,t.max_purchases_per_user,t.status FROM ticket_waves w JOIN events e ON e.id=w.event_id LEFT JOIN ticket_tiers t ON t.wave_id=w.id AND t.status NOT IN ('cancelled','closed') WHERE w.event_id=$1 AND e.status IN ('published','pre_event','live','closing') AND w.status NOT IN ('cancelled','closed') ORDER BY w.display_order,w.id,t.created_at`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var waves []*models.TicketWave
	for rows.Next() {
		var w models.TicketWave
		var t models.TicketTier
		var tid *uuid.UUID
		if err := rows.Scan(&w.ID, &w.EventID, &w.Name, &w.DisplayOrder, &w.SaleStartAt, &w.SaleEndAt, &w.ReleaseMode, &w.Status, &tid, &t.WaveID, &t.Name, &t.Description, &t.PriceMinor, &t.Currency, &t.Inventory, &t.RemainingInventory, &t.AttendeesPerPurchase, &t.MaxPurchasesPerUser, &t.Status); err != nil {
			return nil, err
		}
		if len(waves) == 0 || waves[len(waves)-1].ID != w.ID {
			w.Tiers = []*models.TicketTier{}
			waves = append(waves, &w)
		}
		if tid != nil {
			t.ID = *tid
			waves[len(waves)-1].Tiers = append(waves[len(waves)-1].Tiers, &t)
		}
	}
	return waves, rows.Err()
}

func (r *pgRepository) ReleaseWave(ctx context.Context, userID, waveID uuid.UUID) (*models.TicketWave, error) {
	row := r.db.QueryRow(ctx, `UPDATE ticket_waves w SET status='on_sale', updated_at=now() FROM events e WHERE w.id=$1 AND e.id=w.event_id AND e.organizer_id=$2 AND e.status='published' AND w.status IN ('not_released','scheduled') RETURNING w.id,w.event_id,w.name,w.display_order,w.sale_start_at,w.sale_end_at,w.release_mode,w.status`, waveID, userID)
	return scanWave(row)
}

func (r *pgRepository) AdvanceSales(ctx context.Context, now time.Time, limit int) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	closed, err := tx.Exec(ctx, `UPDATE ticket_waves SET status='closed',updated_at=now() WHERE status IN ('scheduled','on_sale') AND sale_end_at IS NOT NULL AND sale_end_at <= $1`, now)
	if err != nil {
		return 0, err
	}
	opened, err := tx.Exec(ctx, `UPDATE ticket_waves SET status='on_sale',updated_at=now() WHERE status IN ('scheduled','not_released') AND (sale_start_at IS NOT NULL AND sale_start_at <= $1 OR (release_mode='after_previous_sold_out' AND NOT EXISTS (SELECT 1 FROM ticket_waves previous WHERE previous.event_id=ticket_waves.event_id AND previous.display_order < ticket_waves.display_order AND previous.status NOT IN ('sold_out','closed','cancelled')))) AND (sale_end_at IS NULL OR sale_end_at > $1) AND EXISTS (SELECT 1 FROM events e WHERE e.id=ticket_waves.event_id AND e.status IN ('published','pre_event','live','closing'))`, now)
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, `UPDATE ticket_tiers SET status=CASE WHEN remaining_inventory=0 THEN 'sold_out' ELSE 'on_sale' END,updated_at=now() WHERE status IN ('not_released','on_sale') AND EXISTS (SELECT 1 FROM ticket_waves w WHERE w.id=ticket_tiers.wave_id AND w.status='on_sale')`); err != nil {
		return 0, err
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return closed.RowsAffected() + opened.RowsAffected(), nil
}

func (r *pgRepository) Reserve(ctx context.Context, tierID, userID uuid.UUID, dto dtos.ReserveDTO, hold time.Duration) (*models.TicketReservation, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var existing models.TicketReservation
	err = tx.QueryRow(ctx, `SELECT id,tier_id,user_id,quantity,attendee_quantity,status,expires_at FROM ticket_reservations WHERE user_id=$1 AND idempotency_key=$2`, userID, dto.IdempotencyKey).Scan(&existing.ID, &existing.TierID, &existing.UserID, &existing.Quantity, &existing.AttendeeQuantity, &existing.Status, &existing.ExpiresAt)
	if err == nil {
		return &existing, nil
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}
	var perUser, attendees int
	if err := tx.QueryRow(ctx, `SELECT max_purchases_per_user,attendees_per_purchase FROM ticket_tiers WHERE id=$1 AND status='on_sale' FOR UPDATE`, tierID).Scan(&perUser, &attendees); err != nil {
		return nil, ErrUnavailable
	}
	var purchased int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(SUM(quantity),0) FROM ticket_reservations WHERE tier_id=$1 AND user_id=$2 AND status IN ('active','converted')`, tierID, userID).Scan(&purchased); err != nil {
		return nil, err
	}
	if purchased+dto.Quantity > perUser {
		return nil, ErrUnavailable
	}
	expires := time.Now().UTC().Add(hold)
	row := tx.QueryRow(ctx, `UPDATE ticket_tiers SET remaining_inventory=remaining_inventory-$2,updated_at=now() WHERE id=$1 AND remaining_inventory >= $2 RETURNING $1`, tierID, dto.Quantity)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return nil, ErrUnavailable
	}
	reservation := tx.QueryRow(ctx, `INSERT INTO ticket_reservations(tier_id,user_id,quantity,attendee_quantity,idempotency_key,expires_at) VALUES($1,$2,$3,$3*$4,$5,$6) RETURNING id,tier_id,user_id,quantity,attendee_quantity,status,expires_at`, tierID, userID, dto.Quantity, attendees, dto.IdempotencyKey, expires)
	result := &models.TicketReservation{}
	if err := reservation.Scan(&result.ID, &result.TierID, &result.UserID, &result.Quantity, &result.AttendeeQuantity, &result.Status, &result.ExpiresAt); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *pgRepository) ReleaseExpired(ctx context.Context, now time.Time, limit int) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `SELECT id,tier_id,quantity FROM ticket_reservations WHERE status='active' AND expires_at <= $1 ORDER BY expires_at FOR UPDATE SKIP LOCKED LIMIT $2`, now, limit)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var count int64
	for rows.Next() {
		var id, tier uuid.UUID
		var qty int
		if err := rows.Scan(&id, &tier, &qty); err != nil {
			return 0, err
		}
		if _, err = tx.Exec(ctx, `UPDATE ticket_reservations SET status='expired',updated_at=now() WHERE id=$1`, id); err != nil {
			return 0, err
		}
		if _, err = tx.Exec(ctx, `UPDATE ticket_tiers SET remaining_inventory=LEAST(inventory,remaining_inventory+$2),updated_at=now() WHERE id=$1`, tier, qty); err != nil {
			return 0, err
		}
		count++
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return count, nil
}

func scanWave(s interface{ Scan(...any) error }) (*models.TicketWave, error) {
	var w models.TicketWave
	err := s.Scan(&w.ID, &w.EventID, &w.Name, &w.DisplayOrder, &w.SaleStartAt, &w.SaleEndAt, &w.ReleaseMode, &w.Status)
	return &w, err
}
func scanTier(s interface{ Scan(...any) error }) (*models.TicketTier, error) {
	var t models.TicketTier
	err := s.Scan(&t.ID, &t.WaveID, &t.Name, &t.Description, &t.PriceMinor, &t.Currency, &t.Inventory, &t.RemainingInventory, &t.AttendeesPerPurchase, &t.MaxPurchasesPerUser, &t.Status)
	return &t, err
}
func parseTimes(start, end string) (*time.Time, *time.Time, error) {
	var a, b *time.Time
	for value, target := range map[string]**time.Time{start: &a, end: &b} {
		if value == "" {
			continue
		}
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, nil, err
		}
		*target = &parsed
	}
	return a, b, nil
}
func defaultRelease(value string) string {
	if value == "after_previous_sold_out" || value == "manual" {
		return value
	}
	return "scheduled"
}
