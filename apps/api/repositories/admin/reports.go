package admin

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var reportTargets = map[models.ReportTargetType]struct {
	table string
	owner string
}{
	models.ReportTargetPersona:   {table: "personas", owner: "user_id"},
	models.ReportTargetPost:      {table: "posts", owner: "author_user_id"},
	models.ReportTargetComment:   {table: "comments", owner: "author_user_id"},
	models.ReportTargetStory:     {table: "stories", owner: "owner_user_id"},
	models.ReportTargetStoryItem: {table: "story_items", owner: "contributor_user_id"},
	models.ReportTargetSet:       {table: "sets", owner: "author_user_id"},
	models.ReportTargetEvent:     {table: "events", owner: "organizer_id"},
}

func (r *pgRepository) CreateReport(ctx context.Context, params CreateReportParams) (*models.Report, error) {
	if params.TargetType == models.ReportTargetMessage {
		return r.createMessageReport(ctx, params)
	}
	target, ok := reportTargets[params.TargetType]
	if !ok {
		return nil, ErrModerationEntityNotFound
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var ownerID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT `+target.owner+` FROM `+target.table+` WHERE id = $1`, params.TargetID).Scan(&ownerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrModerationEntityNotFound
		}
		return nil, err
	}
	if ownerID == params.ReporterUserID {
		return nil, ErrSelfReport
	}
	row := tx.QueryRow(ctx, `
		INSERT INTO reports (reporter_user_id, reporter_persona_id, target_type, target_id, reason, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, reporter_user_id, reporter_persona_id, target_type, target_id, reason, description,
		          status, assigned_admin_id, resolution_note, created_at, updated_at, resolved_at
	`, params.ReporterUserID, params.ReporterPersonaID, params.TargetType, params.TargetID, params.Reason, params.Description)
	report, err := scanReport(row)
	if err != nil {
		if isAdminUniqueViolation(err) {
			return nil, ErrReportAlreadyExists
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return report, nil
}

func (r *pgRepository) createMessageReport(ctx context.Context, params CreateReportParams) (*models.Report, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var senderID, conversationID, senderPersonaID uuid.UUID
	var body string
	var messageType models.MessageType
	var createdAt time.Time
	var editedAt, deletedAt *time.Time
	err = tx.QueryRow(ctx, `SELECT m.sender_user_id, m.conversation_id, m.sender_persona_id, m.body, m.message_type, m.created_at, m.edited_at, m.deleted_at
		FROM messages m JOIN conversation_members cm ON cm.conversation_id = m.conversation_id
		WHERE m.id = $1 AND cm.user_id = $2 AND cm.left_at IS NULL`, params.TargetID, params.ReporterUserID).
		Scan(&senderID, &conversationID, &senderPersonaID, &body, &messageType, &createdAt, &editedAt, &deletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrModerationEntityNotFound
	}
	if err != nil {
		return nil, err
	}
	if senderID == params.ReporterUserID {
		return nil, ErrSelfReport
	}
	row := tx.QueryRow(ctx, `INSERT INTO reports (reporter_user_id, reporter_persona_id, target_type, target_id, reason, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, reporter_user_id, reporter_persona_id, target_type, target_id, reason, description, status, assigned_admin_id, resolution_note, created_at, updated_at, resolved_at`,
		params.ReporterUserID, params.ReporterPersonaID, params.TargetType, params.TargetID, params.Reason, params.Description)
	report, err := scanReport(row)
	if err != nil {
		if isAdminUniqueViolation(err) {
			return nil, ErrReportAlreadyExists
		}
		return nil, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO message_report_evidence
		(report_id, message_id, conversation_id, sender_user_id, sender_persona_id, body, message_type, attachments, message_created_at, message_edited_at, message_deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE((SELECT jsonb_agg(jsonb_build_object('media_asset_id', media_asset_id, 'position', position) ORDER BY position) FROM message_attachments WHERE message_id = $2), '[]'::jsonb), $8, $9, $10)`,
		report.ID, params.TargetID, conversationID, senderID, senderPersonaID, body, messageType, createdAt, editedAt, deletedAt)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return report, nil
}

func (r *pgRepository) ListReports(ctx context.Context, params ListReportsParams) ([]models.Report, error) {
	limit, offset := reportPage(params.Limit, params.Offset)
	rows, err := r.db.Query(ctx, `
		SELECT id, reporter_user_id, reporter_persona_id, target_type, target_id, reason, description,
		       status, assigned_admin_id, resolution_note, created_at, updated_at, resolved_at
		FROM reports
		WHERE ($1::text IS NULL OR status = $1)
		  AND ($2::text IS NULL OR target_type = $2)
		  AND ($3::uuid IS NULL OR assigned_admin_id = $3)
		ORDER BY CASE WHEN status IN ('open', 'reviewing') THEN 0 ELSE 1 END, created_at ASC
		LIMIT $4 OFFSET $5
	`, params.Status, params.TargetType, params.AssignedID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reports := make([]models.Report, 0)
	for rows.Next() {
		report, err := scanReport(rows)
		if err != nil {
			return nil, err
		}
		reports = append(reports, *report)
	}
	return reports, rows.Err()
}

func (r *pgRepository) FindReportByID(ctx context.Context, reportID uuid.UUID) (*models.Report, error) {
	report, err := scanReport(r.db.QueryRow(ctx, reportSelect+` WHERE id = $1`, reportID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrReportNotFound
	}
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (r *pgRepository) UpdateReport(ctx context.Context, params UpdateReportParams) (*models.Report, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE reports
		SET status = $2, assigned_admin_id = COALESCE($3, assigned_admin_id),
		    resolution_note = CASE WHEN $4 <> '' THEN $4 ELSE resolution_note END,
		    resolved_at = CASE WHEN $2 IN ('resolved', 'dismissed') THEN COALESCE(resolved_at, now()) ELSE NULL END,
		    updated_at = now()
		WHERE id = $1
		RETURNING id, reporter_user_id, reporter_persona_id, target_type, target_id, reason, description,
		          status, assigned_admin_id, resolution_note, created_at, updated_at, resolved_at
	`, params.ReportID, params.Status, params.AssignedAdminID, params.ResolutionNote)
	report, err := scanReport(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrReportNotFound
	}
	return report, err
}

func (r *pgRepository) ResolveReport(ctx context.Context, params ResolveReportParams) (*models.Report, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var targetType models.ReportTargetType
	var targetID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT target_type, target_id FROM reports WHERE id = $1 AND status IN ('open', 'reviewing') FOR UPDATE`, params.ReportID).Scan(&targetType, &targetID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}
	target, ok := reportTargets[targetType]
	if !ok {
		return nil, ErrModerationEntityNotFound
	}
	var previous models.ModerationStatus
	if err := tx.QueryRow(ctx, `SELECT moderation_status FROM `+target.table+` WHERE id = $1 FOR UPDATE`, targetID).Scan(&previous); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrModerationEntityNotFound
		}
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE `+target.table+` SET moderation_status = $2, moderation_reason = $3,
		moderated_at = now(), moderated_by_user_id = $4 WHERE id = $1`, targetID, params.Status, params.Reason, params.AdminID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO moderation_actions
		(entity_type, entity_id, previous_status, status, reason, admin_user_id, report_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, targetType, targetID, previous, params.Status, params.Reason, params.AdminID, params.ReportID); err != nil {
		return nil, err
	}
	report, err := scanReport(tx.QueryRow(ctx, `
		UPDATE reports SET status = 'resolved', assigned_admin_id = $2, resolution_note = $3,
			resolved_at = now(), updated_at = now() WHERE id = $1
		RETURNING id, reporter_user_id, reporter_persona_id, target_type, target_id, reason, description,
			status, assigned_admin_id, resolution_note, created_at, updated_at, resolved_at
	`, params.ReportID, params.AdminID, params.Resolution))
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return report, nil
}

const reportSelect = `SELECT id, reporter_user_id, reporter_persona_id, target_type, target_id, reason, description,
	status, assigned_admin_id, resolution_note, created_at, updated_at, resolved_at FROM reports`

func scanReport(row interface{ Scan(dest ...any) error }) (*models.Report, error) {
	report := &models.Report{}
	err := row.Scan(&report.ID, &report.ReporterUserID, &report.ReporterPersonaID, &report.TargetType,
		&report.TargetID, &report.Reason, &report.Description, &report.Status, &report.AssignedAdminID,
		&report.ResolutionNote, &report.CreatedAt, &report.UpdatedAt, &report.ResolvedAt)
	return report, err
}

func reportPage(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
