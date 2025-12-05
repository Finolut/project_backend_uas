package postgres

import "time"

type ActivityLog struct {
	ID         string                 `db:"id" json:"id"`
	EntityType string                 `db:"entity_type" json:"entity_type"`
	EntityID   string                 `db:"entity_id" json:"entity_id"`
	EventType  string                 `db:"event_type" json:"event_type"`
	ActorID    *string                `db:"actor_id" json:"actor_id,omitempty"`
	ActorRole  *string                `db:"actor_role" json:"actor_role,omitempty"`
	Previous   map[string]interface{} `db:"previous" json:"previous,omitempty"` // use jsonb
	Current    map[string]interface{} `db:"current" json:"current,omitempty"`
	Metadata   map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	CreatedAt  time.Time              `db:"created_at" json:"created_at"`
}
