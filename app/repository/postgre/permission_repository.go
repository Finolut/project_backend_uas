package postgre

import (
	"context"
	"database/sql"

	pgmodel "app/model/postgre"
)

// PermissionRepository defines data access for permissions.
type PermissionRepository interface {
	Create(ctx context.Context, p *pgmodel.Permission) error
	GetByID(ctx context.Context, id string) (*pgmodel.Permission, error)
	GetByName(ctx context.Context, name string) (*pgmodel.Permission, error)
	ListAll(ctx context.Context) ([]*pgmodel.Permission, error)
}

// Implementation
type permissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, p *pgmodel.Permission) error {
	q := `INSERT INTO permissions (id, name, resource, action, description) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Resource, p.Action, p.Description)
	return err
}

func (r *permissionRepository) GetByID(ctx context.Context, id string) (*pgmodel.Permission, error) {
	var out pgmodel.Permission
	q := `SELECT id, name, resource, action, description FROM permissions WHERE id=$1`
	row := r.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(&out.ID, &out.Name, &out.Resource, &out.Action, &out.Description); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *permissionRepository) GetByName(ctx context.Context, name string) (*pgmodel.Permission, error) {
	var out pgmodel.Permission
	q := `SELECT id, name, resource, action, description FROM permissions WHERE name=$1`
	row := r.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&out.ID, &out.Name, &out.Resource, &out.Action, &out.Description); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *permissionRepository) ListAll(ctx context.Context) ([]*pgmodel.Permission, error) {
	q := `SELECT id, name, resource, action, description FROM permissions ORDER BY resource, action`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*pgmodel.Permission
	for rows.Next() {
		var p pgmodel.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, nil
}
