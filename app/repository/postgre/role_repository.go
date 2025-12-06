package postgre

import (
	"context"
	"database/sql"
	"time"

	pgmodel "clean-arch-copy/app/model/postgre"
)

// -----------------------------
// INTERFACE
// -----------------------------
type RoleRepository interface {
	Create(ctx context.Context, r *pgmodel.Role) error
	GetByID(ctx context.Context, id string) (*pgmodel.Role, error)
	GetByName(ctx context.Context, name string) (*pgmodel.Role, error)
	ListAll(ctx context.Context) ([]*pgmodel.Role, error)
}

// -----------------------------
// IMPLEMENTATION
// -----------------------------
type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *pgmodel.Role) error {
	role.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO roles (id, name, description, created_at)
		 VALUES ($1,$2,$3,$4)`,
		role.ID, role.Name, role.Description, role.CreatedAt,
	)

	return err
}


func (r *roleRepository) GetByID(ctx context.Context, id string) (*pgmodel.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles WHERE id=$1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var out pgmodel.Role
	err := row.Scan(&out.ID, &out.Name, &out.Description, &out.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*pgmodel.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles WHERE name=$1
	`
	row := r.db.QueryRowContext(ctx, query, name)

	var out pgmodel.Role
	err := row.Scan(&out.ID, &out.Name, &out.Description, &out.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *roleRepository) ListAll(ctx context.Context) ([]*pgmodel.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*pgmodel.Role
	for rows.Next() {
		var r pgmodel.Role
		err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, &r)
	}
	return out, nil
}
