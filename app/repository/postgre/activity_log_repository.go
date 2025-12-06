package postgre

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	pgmodel "clean-arch-copy/app/model/postgre"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *pgmodel.ActivityLog) error
	ListByEntity(ctx context.Context, entityType string, entityID string, limit, offset int) ([]*pgmodel.ActivityLog, error)
}

type activityLogRepo struct {
	db *sql.DB
}

func NewActivityLogRepository(db *sql.DB) ActivityLogRepository {
	return &activityLogRepo{db: db}
}

func toJSONb(m map[string]interface{}) ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (r *activityLogRepo) Create(ctx context.Context, l *pgmodel.ActivityLog) error {
	now := time.Now()
	if l.CreatedAt.IsZero() {
		l.CreatedAt = now
	}
	prevBytes, _ := toJSONb(l.Previous)
	currBytes, _ := toJSONb(l.Current)
	metaBytes, _ := toJSONb(l.Metadata)

	q := `INSERT INTO activity_logs (id, entity_type, entity_id, event_type, actor_id, actor_role, previous, current, metadata, created_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, q,
		l.ID, l.EntityType, l.EntityID, l.EventType, l.ActorID, l.ActorRole,
		prevBytes, currBytes, metaBytes, l.CreatedAt,
	)
	return err
}

func (r *activityLogRepo) ListByEntity(ctx context.Context, entityType string, entityID string, limit, offset int) ([]*pgmodel.ActivityLog, error) {
	q := `SELECT id, entity_type, entity_id, event_type, actor_id, actor_role, previous, current, metadata, created_at
	      FROM activity_logs
	      WHERE entity_type=$1 AND entity_id=$2
	      ORDER BY created_at DESC
	      LIMIT $3 OFFSET $4`
	rows, err := r.db.QueryContext(ctx, q, entityType, entityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*pgmodel.ActivityLog
	for rows.Next() {
		var l pgmodel.ActivityLog
		var prevBytes, currBytes, metaBytes sql.NullString
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.EventType, &l.ActorID, &l.ActorRole, &prevBytes, &currBytes, &metaBytes, &l.CreatedAt); err != nil {
			return nil, err
		}
		// unmarshal jsonb fields
		if prevBytes.Valid {
			_ = json.Unmarshal([]byte(prevBytes.String), &l.Previous)
		}
		if currBytes.Valid {
			_ = json.Unmarshal([]byte(currBytes.String), &l.Current)
		}
		if metaBytes.Valid {
			_ = json.Unmarshal([]byte(metaBytes.String), &l.Metadata)
		}
		out = append(out, &l)
	}
	return out, nil
}
