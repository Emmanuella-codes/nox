package pipes

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *AdminPipe) getMessageEvidencePipe(ctx context.Context, actorID, reportID, messageID uuid.UUID, reason string) *shared.PipeRes[models.AdminContentRecord] {
	var conversationID, senderUserID, senderPersonaID uuid.UUID
	var body, messageType string
	var attachments json.RawMessage
	var createdAt time.Time
	var editedAt, deletedAt *time.Time
	err := p.db.QueryRow(ctx, `SELECT conversation_id, sender_user_id, sender_persona_id, body, message_type, attachments, message_created_at, message_edited_at, message_deleted_at
		FROM message_report_evidence WHERE report_id = $1 AND expires_at > now()`, reportID).
		Scan(&conversationID, &senderUserID, &senderPersonaID, &body, &messageType, &attachments, &createdAt, &editedAt, &deletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return shared.PipeError[models.AdminContentRecord](messages.Moderation_Entity_Not_Found)
	}
	if err != nil {
		logInternalError(err, "message_evidence.read")
		return shared.PipeError[models.AdminContentRecord](messages.Internal_Error)
	}
	data := map[string]any{"conversation_id": conversationID, "sender_user_id": senderUserID, "sender_persona_id": senderPersonaID, "body": body, "message_type": messageType, "attachments": json.RawMessage(attachments), "created_at": createdAt, "edited_at": editedAt, "deleted_at": deletedAt}
	item := &models.AdminContentRecord{EntityType: models.ModerationEntityType(models.ReportTargetMessage), EntityID: messageID, Data: data, ModerationStatus: models.ModerationStatusActive}
	p.audit(ctx, actorID, "admin.message_evidence.read", map[string]any{"report_id": reportID.String(), "message_id": messageID.String(), "reason": strings.TrimSpace(reason)})
	return shared.PipeSuccess(messages.Admin_Message_Evidence_Loaded, item)
}

func (p *AdminPipe) resolveMessageReport(ctx context.Context, actorID, reportID uuid.UUID, reason, resolution string) (*models.Report, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var targetID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT target_id FROM reports WHERE id = $1 AND target_type = 'message' AND status IN ('open', 'reviewing') FOR UPDATE`, reportID).Scan(&targetID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, adminrepo.ErrReportNotFound
		}
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE messages SET body = '', media_asset_id = NULL, deleted_at = COALESCE(deleted_at, now()) WHERE id = $1`, targetID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM message_attachments WHERE message_id = $1`, targetID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE conversations c SET last_message_id = (
		SELECT m.id FROM messages m WHERE m.conversation_id = c.id AND m.deleted_at IS NULL
		ORDER BY m.created_at DESC, m.id DESC LIMIT 1), updated_at = now()
		WHERE c.id = (SELECT conversation_id FROM messages WHERE id = $1)`, targetID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO moderation_actions (entity_type, entity_id, previous_status, status, reason, admin_user_id, report_id) VALUES ('message', $1, NULL, 'removed', $2, $3, $4)`, targetID, reason, actorID, reportID); err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, `UPDATE reports SET status = 'resolved', assigned_admin_id = $2, resolution_note = $3, resolved_at = now(), updated_at = now() WHERE id = $1 RETURNING id, reporter_user_id, reporter_persona_id, target_type, target_id, reason, description, status, assigned_admin_id, resolution_note, created_at, updated_at, resolved_at`, reportID, actorID, resolution)
	report, err := scanAdminReport(row)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return report, nil
}

func scanAdminReport(row interface{ Scan(dest ...any) error }) (*models.Report, error) {
	var report models.Report
	err := row.Scan(&report.ID, &report.ReporterUserID, &report.ReporterPersonaID, &report.TargetType, &report.TargetID, &report.Reason, &report.Description, &report.Status, &report.AssignedAdminID, &report.ResolutionNote, &report.CreatedAt, &report.UpdatedAt, &report.ResolvedAt)
	return &report, err
}
