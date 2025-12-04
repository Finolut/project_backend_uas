package postgres

type Permission struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`         // e.g. achievement:create
	Resource    string `db:"resource" json:"resource"` // e.g. achievement
	Action      string `db:"action" json:"action"`     // e.g. create, read
	Description string `db:"description" json:"description"`
}
