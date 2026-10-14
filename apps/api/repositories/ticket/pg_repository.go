package ticket

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/ticket/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateOrder(ctx context.Context, reservationID, userID uuid.UUID) (*models.TicketOrder, error) {
	reference := "nox_" + uuid.NewString()
	row := r.db.QueryRow(ctx, `INSERT INTO ticket_orders(user_id,event_id,reservation_id,subtotal_minor,provider_fee_minor,platform_fee_minor,total_minor,currency,provider_reference,expires_at) SELECT r.user_id,w.event_id,r.id,t.price_minor*r.quantity,0,0,t.price_minor*r.quantity,t.currency,$3,r.expires_at FROM ticket_reservations r JOIN ticket_tiers t ON t.id=r.tier_id JOIN ticket_waves w ON w.id=t.wave_id WHERE r.id=$1 AND r.user_id=$2 AND r.status IN ('active','converted') ON CONFLICT (reservation_id) DO UPDATE SET updated_at=now() RETURNING id,user_id,event_id,reservation_id,subtotal_minor,provider_fee_minor,platform_fee_minor,total_minor,currency,status,provider_reference,expires_at`, reservationID, userID, reference)
	return scanOrder(row)
}

func (r *pgRepository) FindOrder(ctx context.Context, orderID, userID uuid.UUID) (*models.TicketOrder, error) {
	row := r.db.QueryRow(ctx, `SELECT id,user_id,event_id,reservation_id,subtotal_minor,provider_fee_minor,platform_fee_minor,total_minor,currency,status,provider_reference,expires_at FROM ticket_orders WHERE id=$1 AND user_id=$2`, orderID, userID)
	return scanOrder(row)
}

func (r *pgRepository) SettlePayment(ctx context.Context, reference, providerEventID string, amount int64, currency string, payload json.RawMessage) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var orderID, reservationID, tierID uuid.UUID
	var userID, eventID uuid.UUID
	var quantity, attendees int
	var subtotal int64
	var orderCurrency string
	var status, reservationStatus string
	if err := tx.QueryRow(ctx, `SELECT o.id,o.user_id,o.event_id,o.reservation_id,o.subtotal_minor,o.currency,o.status,r.status,r.quantity,r.tier_id,t.attendees_per_purchase FROM ticket_orders o JOIN ticket_reservations r ON r.id=o.reservation_id JOIN ticket_tiers t ON t.id=r.tier_id WHERE o.provider_reference=$1 FOR UPDATE`, reference).Scan(&orderID, &userID, &eventID, &reservationID, &subtotal, &orderCurrency, &status, &reservationStatus, &quantity, &tierID, &attendees); err != nil {
		return ErrNotFound
	}
	if status == "paid" {
		return nil
	}
	if amount < subtotal || !strings.EqualFold(currency, orderCurrency) {
		return ErrUnavailable
	}
	if reservationStatus == "expired" {
		result, err := tx.Exec(ctx, `UPDATE ticket_tiers SET remaining_inventory=remaining_inventory-$2,updated_at=now() WHERE id=$1 AND remaining_inventory >= $2`, tierID, quantity)
		if err != nil {
			return err
		}
		if result.RowsAffected() != 1 {
			return ErrUnavailable
		}
		if _, err := tx.Exec(ctx, `UPDATE ticket_reservations SET status='active',updated_at=now() WHERE id=$1`, reservationID); err != nil {
			return err
		}
	} else if reservationStatus != "active" {
		return ErrUnavailable
	}
	if _, err := tx.Exec(ctx, `INSERT INTO ticket_payments(order_id,provider_event_id,provider_reference,status,amount_minor,currency,payload,processed_at) VALUES($1,$2,$3,'success',$4,$5,$6,now()) ON CONFLICT(provider_event_id) DO NOTHING`, orderID, providerEventID, reference, amount, currency, payload); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE ticket_orders SET status='paid',total_minor=$2,provider_fee_minor=GREATEST($2-subtotal_minor,0),updated_at=now() WHERE id=$1`, orderID, amount); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE ticket_reservations SET status='converted',updated_at=now() WHERE id=$1`, reservationID); err != nil {
		return err
	}
	for i := 0; i < quantity; i++ {
		var ticketID uuid.UUID
		ticketCredential := "nox_" + uuid.NewString()
		if err := tx.QueryRow(ctx, `INSERT INTO tickets(order_id,event_id,tier_id,credential) VALUES($1,$2,$3,$4) RETURNING id`, orderID, eventID, tierID, ticketCredential).Scan(&ticketID); err != nil {
			return err
		}
		for j := 0; j < attendees; j++ {
			var entitlementID, attendeeID uuid.UUID
			if err := tx.QueryRow(ctx, `INSERT INTO attendee_entitlements(ticket_id,event_id,user_id) VALUES($1,$2,$3) RETURNING id`, ticketID, eventID, userID).Scan(&entitlementID); err != nil {
				return err
			}
			if err := tx.QueryRow(ctx, `INSERT INTO event_attendees(event_id,user_id,display_name) SELECT $1,$2,COALESCE(fullname,'Attendee') FROM users WHERE id=$2 RETURNING id`, eventID, userID).Scan(&attendeeID); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `INSERT INTO event_admissions(attendee_id,event_id,source_type,source_id) VALUES($1,$2,'paid_ticket',$3)`, attendeeID, eventID, entitlementID); err != nil {
				return err
			}
			credential := ticketCredential
			if j > 0 {
				credential = "nox_" + uuid.NewString()
			}
			hash := sha256.Sum256([]byte(credential))
			if _, err := tx.Exec(ctx, `INSERT INTO admission_credentials(admission_id,credential_hash) SELECT id,$2 FROM event_admissions WHERE source_type='paid_ticket' AND source_id=$1`, entitlementID, hex.EncodeToString(hash[:])); err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

func scanOrder(s interface{ Scan(...any) error }) (*models.TicketOrder, error) {
	var o models.TicketOrder
	err := s.Scan(&o.ID, &o.UserID, &o.EventID, &o.ReservationID, &o.SubtotalMinor, &o.ProviderFeeMinor, &o.PlatformFeeMinor, &o.TotalMinor, &o.Currency, &o.Status, &o.ProviderReference, &o.ExpiresAt)
	return &o, err
}

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
		if _, err = tx.Exec(ctx, `UPDATE ticket_orders SET status='expired',updated_at=now() WHERE reservation_id=$1 AND status='pending'`, id); err != nil {
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
