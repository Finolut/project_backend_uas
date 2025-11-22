package repository

import (
	"clean-arch/app/model/postgre"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func GetAllAlumniWithPagination(db *sql.DB, params model.PaginationParams) ([]model.Alumni, int, error) {
	// Build WHERE clause for search
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	if params.Search != "" {
		whereClause += fmt.Sprintf(" AND (nama ILIKE $%d OR nim ILIKE $%d OR jurusan ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex, argIndex, argIndex)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	// Validate and set sort column
	validSortColumns := map[string]bool{
		"id": true, "nim": true, "nama": true, "jurusan": true,
		"angkatan": true, "tahun_lulus": true, "email": true, "created_at": true,
	}
	if !validSortColumns[params.SortBy] {
		params.SortBy = "created_at"
	}

	// Validate sort order
	if strings.ToLower(params.Order) != "desc" {
		params.Order = "asc"
	}

	// Get total count for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM alumni %s", whereClause)
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build main query with pagination
	offset := (params.Page - 1) * params.Limit
	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
		FROM alumni %s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`,
		whereClause, params.SortBy, params.Order, argIndex, argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
			&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
			&alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		alumniList = append(alumniList, alumni)
	}

	return alumniList, total, nil
}

func GetAllAlumni(db *sql.DB) ([]model.Alumni, error) {
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
	          FROM alumni WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
			&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
			&alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}
	return alumniList, nil
}

func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
	          FROM alumni WHERE id = $1 AND deleted_at IS NULL`

	err := db.QueryRow(query, id).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func CreateAlumni(db *sql.DB, req model.CreateAlumniRequest) (*model.Alumni, error) {
	now := time.Now()
	var id int
	query := `INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err := db.QueryRow(query, req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return GetAlumniByID(db, id)
}

func UpdateAlumni(db *sql.DB, id int, req model.UpdateAlumniRequest) (*model.Alumni, error) {
	now := time.Now()
	query := `UPDATE alumni SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, 
	          email = $5, no_telepon = $6, alamat = $7, updated_at = $8 
			  WHERE id = $9 AND deleted_at IS NULL`

	result, err := db.Exec(query, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, now, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return GetAlumniByID(db, id)
}

func DeleteAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM alumni WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func CheckAlumniByNim(db *sql.DB, nim string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
	          FROM alumni WHERE nim = $1 AND deleted_at IS NULL`

	err := db.QueryRow(query, nim).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func GetAlumniStatistics(db *sql.DB) (*model.AlumniStatistics, error) {
	stats := &model.AlumniStatistics{
		AlumniByJurusan:    make(map[string]int),
		AlumniByAngkatan:   make(map[string]int),
		AlumniByTahunLulus: make(map[string]int),
	}

	// Get total alumni count
	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM alumni WHERE deleted_at IS NULL").Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	stats.TotalAlumni = totalCount

	// Get alumni count by jurusan
	rows, err := db.Query("SELECT jurusan, COUNT(*) FROM alumni WHERE deleted_at IS NULL GROUP BY jurusan ORDER BY jurusan")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var jurusan string
		var count int
		if err := rows.Scan(&jurusan, &count); err != nil {
			return nil, err
		}
		stats.AlumniByJurusan[jurusan] = count
	}

	// Get alumni count by angkatan
	rows, err = db.Query("SELECT angkatan, COUNT(*) FROM alumni WHERE deleted_at IS NULL GROUP BY angkatan ORDER BY angkatan")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var angkatan int
		var count int
		if err := rows.Scan(&angkatan, &count); err != nil {
			return nil, err
		}
		stats.AlumniByAngkatan[string(rune(angkatan+'0'))] = count
	}

	// Get alumni count by tahun lulus
	rows, err = db.Query("SELECT tahun_lulus, COUNT(*) FROM alumni WHERE deleted_at IS NULL GROUP BY tahun_lulus ORDER BY tahun_lulus")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tahunLulus int
		var count int
		if err := rows.Scan(&tahunLulus, &count); err != nil {
			return nil, err
		}
		stats.AlumniByTahunLulus[string(rune(tahunLulus+'0'))] = count
	}

	return stats, nil
}

func GetTrashedAlumni(db *sql.DB) ([]model.Alumni, error) {
	rows, err := db.Query(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at, deleted_at, deleted_by
		FROM alumni
		WHERE deleted_at IS NOT NULL
		ORDER BY deleted_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email,
			&a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt, &a.DeletedBy,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func SoftDeleteAlumni(db *sql.DB, id int, userID *int) error {
	_, err := db.Exec(`
		UPDATE alumni 
		SET deleted_at = $1, deleted_by = $2, updated_at = $1
		WHERE id = $3 AND deleted_at IS NULL`,
		time.Now(), userID, id)
	return err
}

func RestoreAlumni(db *sql.DB, id int) error {
	result, err := db.Exec(`
		UPDATE alumni 
		SET deleted_at = NULL, deleted_by = NULL, updated_at = $1
		WHERE id = $2 AND deleted_at IS NOT NULL`,
		time.Now(), id)
	if err != nil {
		return err
	}
	ra, _ := result.RowsAffected()
	if ra == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func HardDeleteAlumni(db *sql.DB, id int) error {
	// hanya boleh hard delete jika sudah di-trash
	result, err := db.Exec(`DELETE FROM alumni WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return err
	}
	ra, _ := result.RowsAffected()
	if ra == 0 {
		return sql.ErrNoRows
	}
	return nil
}
