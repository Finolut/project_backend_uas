package postgre

import (
	"context"
	"database/sql"
	"time"

	pgmodel "app/model/postgre"
)

// ----------------------
// INTERFACE
// ----------------------
type UserRepository interface {
	Create(ctx context.Context, u *pgmodel.User) error
	GetByID(ctx context.Context, id string) (*pgmodel.User, error)
	GetByUsername(ctx context.Context, username string) (*pgmodel.User, error)
	Update(ctx context.Context, u *pgmodel.User) error
	Delete(ctx context.Context, id string) error
}

// ----------------------
// IMPLEMENTATION
// ----------------------
type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *pgmodel.User) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	query := `
		INSERT INTO users (
			id, username, email, password_hash, full_name, role_id,
			is_active, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Username, u.Email, u.PasswordHash, u.FullName,
		u.RoleID, u.IsActive, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*pgmodel.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name,
		       role_id, is_active, created_at, updated_at
		FROM users WHERE id=$1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	var u pgmodel.User

	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
		&u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*pgmodel.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name,
		       role_id, is_active, created_at, updated_at
		FROM users WHERE username=$1
	`

	row := r.db.QueryRowContext(ctx, query, username)
	var u pgmodel.User

	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
		&u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *pgmodel.User) error {
	u.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET username=$1, email=$2, password_hash=$3, full_name=$4,
		    role_id=$5, is_active=$6, updated_at=$7
		WHERE id=$8
	`
	_, err := r.db.ExecContext(ctx, query,
		u.Username, u.Email, u.PasswordHash, u.FullName,
		u.RoleID, u.IsActive, u.UpdatedAt, u.ID,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
