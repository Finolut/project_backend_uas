package postgre

import (
	"context"
	"database/sql"
	"time"
)

// TokenRepository mendefinisikan operasi untuk blacklist token.
// Interface ini harus cocok dengan yang dibutuhkan oleh AuthService.
type TokenRepository interface {
	AddToBlacklist(ctx context.Context, token string, expiresAt time.Time) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

// Implementation
type tokenRepository struct {
	db *sql.DB
}

// NewTokenRepository membuat instance baru dari tokenRepository
func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

// AddToBlacklist menyimpan token ke tabel blacklist
func (r *tokenRepository) AddToBlacklist(ctx context.Context, token string, expiresAt time.Time) error {
	// Kita simpan expires_at agar nanti bisa dibuat cron job untuk cleanup data sampah
	query := `INSERT INTO token_blacklist (token, expires_at) VALUES ($1, $2) ON CONFLICT (token) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, token, expiresAt)
	return err
}

// IsBlacklisted mengecek apakah token ada di database
func (r *tokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	// Kita gunakan EXISTS agar performa lebih cepat (tidak perlu scan data)
	query := `SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE token=$1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, token).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
