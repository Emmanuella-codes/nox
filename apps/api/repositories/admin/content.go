package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type contentTable struct {
	table  string
	owner  string
	parent string
}

var adminContentTables = map[models.ModerationEntityType]contentTable{
	models.ModerationEntityPersona:   {table: "personas", owner: "user_id"},
	models.ModerationEntityPost:      {table: "posts", owner: "author_user_id"},
	models.ModerationEntityComment:   {table: "comments", owner: "author_user_id"},
	models.ModerationEntityStory:     {table: "stories", owner: "owner_user_id"},
	models.ModerationEntityStoryItem: {table: "story_items", owner: "contributor_user_id", parent: "story_id"},
	models.ModerationEntitySet:       {table: "sets", owner: "author_user_id"},
	models.ModerationEntityEvent:     {table: "events", owner: "organizer_id"},
}

func (r *pgRepository) ListAdminContent(ctx context.Context, params ListAdminContentParams) (*models.AdminContentPage, error) {
	target, ok := adminContentTables[params.EntityType]
	if !ok {
		return nil, ErrModerationEntityNotFound
	}
	limit, offset := reportPage(params.Limit, params.Offset)
	where := []string{"TRUE"}
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if params.Status != nil {
		where = append(where, "t.moderation_status = "+addArg(params.Status))
	}
	if params.OwnerID != nil {
		where = append(where, "t."+target.owner+" = "+addArg(params.OwnerID))
	}
	if params.ParentID != nil && target.parent != "" {
		where = append(where, "t."+target.parent+" = "+addArg(params.ParentID))
	}
	if params.CreatedFrom != nil {
		where = append(where, "t.created_at >= "+addArg(params.CreatedFrom))
	}
	if params.CreatedTo != nil {
		where = append(where, "t.created_at <= "+addArg(params.CreatedTo))
	}
	args = append(args, limit+1, offset)
	query := `SELECT t.id, t.moderation_status, t.moderation_reason, t.moderated_at, t.moderated_by_user_id,
		to_jsonb(t) - ARRAY['password', 'user_id', 'author_user_id', 'owner_user_id', 'owner_persona_id', 'contributor_user_id', 'contributor_persona_id', 'persona_id', 'moderation_status', 'moderation_reason', 'moderated_at', 'moderated_by_user_id']::text[]
		FROM ` + target.table + ` t WHERE ` + joinStrings(where, " AND ") + `
		ORDER BY t.created_at DESC, t.id DESC LIMIT $` + fmt.Sprint(len(args)-1) + ` OFFSET $` + fmt.Sprint(len(args))
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AdminContentRecord, 0, limit)
	for rows.Next() {
		item, err := scanAdminContentRecord(rows, params.EntityType)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	page := &models.AdminContentPage{EntityType: string(params.EntityType), Limit: limit, Offset: offset, Items: items}
	if len(items) > limit {
		page.HasMore = true
		page.Items = items[:limit]
	}
	return page, nil
}

func (r *pgRepository) FindAdminContent(ctx context.Context, entityType models.ModerationEntityType, entityID uuid.UUID) (*models.AdminContentRecord, error) {
	target, ok := adminContentTables[entityType]
	if !ok {
		return nil, ErrModerationEntityNotFound
	}
	row := r.db.QueryRow(ctx, `SELECT t.id, t.moderation_status, t.moderation_reason, t.moderated_at,
		t.moderated_by_user_id, to_jsonb(t) - ARRAY['password', 'user_id', 'author_user_id', 'owner_user_id', 'owner_persona_id', 'contributor_user_id', 'contributor_persona_id', 'persona_id', 'moderation_status', 'moderation_reason', 'moderated_at', 'moderated_by_user_id']::text[]
		FROM `+target.table+` t WHERE t.id = $1`, entityID)
	item, err := scanAdminContentRecord(row, entityType)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrModerationEntityNotFound
	}
	return item, err
}

func scanAdminContentRecord(row interface{ Scan(dest ...any) error }, entityType models.ModerationEntityType) (*models.AdminContentRecord, error) {
	item := &models.AdminContentRecord{EntityType: entityType}
	var raw json.RawMessage
	if err := row.Scan(&item.EntityID, &item.ModerationStatus, &item.ModerationReason, &item.ModeratedAt, &item.ModeratedBy, &raw); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &item.Data); err != nil {
		return nil, err
	}
	return item, nil
}

func joinStrings(values []string, separator string) string {
	return strings.Join(values, separator)
}
