package postgre

import (
	"context"
	"database/sql"

	pgmodel "app/model/postgre"
)

// RolePermissionRepository manages assignment between roles and permissions.
type RolePermissionRepository interface {
	Assign(ctx context.Context, roleID string, permissionID string) error
	Remove(ctx context.Context, roleID string, permissionID string) error
	ListByRole(ctx context.Context, roleID string) ([]*pgmodel.Permission, error)
}

// Implementation
type rolePermissionRepository struct {
	db *sql.DB
}

func NewRolePermissionRepository(db *sql.DB) RolePermissionRepository {
	return &rolePermissionRepository{db: db}
}

func (r *rolePermissionRepository) Assign(ctx context.Context, roleID string, permissionID string) error {
	q := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, q, roleID, permissionID)
	return err
}

func (r *rolePermissionRepository) Remove(ctx context.Context, roleID string, permissionID string) error {
	q := `DELETE FROM role_permissions WHERE role_id=$1 AND permission_id=$2`
	_, err := r.db.ExecContext(ctx, q, roleID, permissionID)
	return err
}

func (r *rolePermissionRepository) ListByRole(ctx context.Context, roleID string) ([]*pgmodel.Permission, error) {
	q := `
	SELECT p.id, p.name, p.resource, p.action, p.description
	FROM permissions p
	JOIN role_permissions rp ON rp.permission_id = p.id
	WHERE rp.role_id = $1
	ORDER BY p.resource, p.action
	`
	rows, err := r.db.QueryContext(ctx, q, roleID)
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
