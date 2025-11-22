package repository

import (
	"clean-arch/app/model/postgre"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func GetAllPekerjaanWithPagination(db *sql.DB, params model.PaginationParams) ([]model.PekerjaanAlumni, int, error) {
	// Build WHERE clause for search
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	if params.Search != "" {
		whereClause += " AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1 OR status_pekerjaan ILIKE $1)"
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	// Validate and set sort column
	validSortColumns := map[string]bool{
		"id": true, "nama_perusahaan": true, "posisi_jabatan": true, "bidang_industri": true,
		"lokasi_kerja": true, "status_pekerjaan": true, "tanggal_mulai_kerja": true, "created_at": true,
	}
	if !validSortColumns[params.SortBy] {
		params.SortBy = "created_at"
	}

	// Validate sort order
	if strings.ToLower(params.Order) != "desc" {
		params.Order = "asc"
	}

	// Get total count for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM pekerjaan_alumni %s", whereClause)
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build main query with pagination
	offset := (params.Page - 1) * params.Limit
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
		       deskripsi_pekerjaan, deleted_at, deleted_by, created_at, updated_at 
		FROM pekerjaan_alumni %s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`,
		whereClause, params.SortBy, params.Order, argIndex, argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		var tanggalMulai time.Time
		var tanggalSelesai *time.Time

		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
			&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
			&pekerjaan.GajiRange, &tanggalMulai, &tanggalSelesai,
			&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan,
			&pekerjaan.DeletedAt, &pekerjaan.DeletedBy,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		pekerjaan.TanggalMulaiKerja = model.Date{Time: tanggalMulai}
		if tanggalSelesai != nil {
			pekerjaan.TanggalSelesaiKerja = &model.Date{Time: *tanggalSelesai}
		}

		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, total, nil
}

func GetAllPekerjaan(db *sql.DB) ([]model.PekerjaanAlumni, error) {
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
	          gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
	          deskripsi_pekerjaan, deleted_at, deleted_by, created_at, updated_at 
	          FROM pekerjaan_alumni WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		var tanggalMulai time.Time
		var tanggalSelesai *time.Time

		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
			&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
			&pekerjaan.GajiRange, &tanggalMulai, &tanggalSelesai,
			&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan,
			&pekerjaan.DeletedAt, &pekerjaan.DeletedBy,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		pekerjaan.TanggalMulaiKerja = model.Date{Time: tanggalMulai}
		if tanggalSelesai != nil {
			pekerjaan.TanggalSelesaiKerja = &model.Date{Time: *tanggalSelesai}
		}

		pekerjaanList = append(pekerjaanList, pekerjaan)
	}
	return pekerjaanList, nil
}

func GetPekerjaanByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	pekerjaan := new(model.PekerjaanAlumni)
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
	          gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
	          deskripsi_pekerjaan, deleted_at, deleted_by, created_at, updated_at 
	          FROM pekerjaan_alumni WHERE id = $1 AND deleted_at IS NULL`

	var tanggalMulai time.Time
	var tanggalSelesai *time.Time

	err := db.QueryRow(query, id).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange, &tanggalMulai, &tanggalSelesai,
		&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan,
		&pekerjaan.DeletedAt, &pekerjaan.DeletedBy,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	pekerjaan.TanggalMulaiKerja = model.Date{Time: tanggalMulai}
	if tanggalSelesai != nil {
		pekerjaan.TanggalSelesaiKerja = &model.Date{Time: *tanggalSelesai}
	}

	return pekerjaan, nil
}

func GetPekerjaanByAlumniID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
	          gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
	          deskripsi_pekerjaan, deleted_at, deleted_by, created_at, updated_at 
	          FROM pekerjaan_alumni WHERE alumni_id = $1 AND deleted_at IS NULL ORDER BY tanggal_mulai_kerja DESC`

	rows, err := db.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		var tanggalMulai time.Time
		var tanggalSelesai *time.Time

		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
			&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
			&pekerjaan.GajiRange, &tanggalMulai, &tanggalSelesai,
			&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan,
			&pekerjaan.DeletedAt, &pekerjaan.DeletedBy,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		pekerjaan.TanggalMulaiKerja = model.Date{Time: tanggalMulai}
		if tanggalSelesai != nil {
			pekerjaan.TanggalSelesaiKerja = &model.Date{Time: *tanggalSelesai}
		}

		pekerjaanList = append(pekerjaanList, pekerjaan)
	}
	return pekerjaanList, nil
}

func CreatePekerjaan(db *sql.DB, req model.CreatePekerjaanRequest) (*model.PekerjaanAlumni, error) {
	now := time.Now()
	var id int
	query := `INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
	          lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
	          deskripsi_pekerjaan, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil {
		tanggalSelesai = &req.TanggalSelesaiKerja.Time
	}

	err := db.QueryRow(query, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan,
		req.BidangIndustri, req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja.Time,
		tanggalSelesai, req.StatusPekerjaan, req.DeskripsiPekerjaan, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return GetPekerjaanByID(db, id)
}

func UpdatePekerjaan(db *sql.DB, id int, req model.UpdatePekerjaanRequest) (*model.PekerjaanAlumni, error) {
	now := time.Now()
	query := `UPDATE pekerjaan_alumni SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3,
	          lokasi_kerja = $4, gaji_range = $5, tanggal_mulai_kerja = $6, tanggal_selesai_kerja = $7,
	          status_pekerjaan = $8, deskripsi_pekerjaan = $9, updated_at = $10 WHERE id = $11`

	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil {
		tanggalSelesai = &req.TanggalSelesaiKerja.Time
	}

	result, err := db.Exec(query, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja.Time, tanggalSelesai,
		req.StatusPekerjaan, req.DeskripsiPekerjaan, now, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return GetPekerjaanByID(db, id)
}

func DeletePekerjaan(db *sql.DB, id int) error {
	query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
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

func SoftDeletePekerjaan(db *sql.DB, id int, deletedBy int) error {
	now := time.Now()
	query := `UPDATE pekerjaan_alumni SET deleted_at = $1, deleted_by = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL`
	result, err := db.Exec(query, now, deletedBy, now, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func GetAlumniIDByPekerjaanID(db *sql.DB, pekerjaanID int) (int, error) {
	var alumniID int
	query := `SELECT alumni_id FROM pekerjaan_alumni WHERE id = $1 AND deleted_at IS NULL`
	err := db.QueryRow(query, pekerjaanID).Scan(&alumniID)
	if err != nil {
		return 0, err
	}
	return alumniID, nil
}

func GetUserIDByAlumniID(db *sql.DB, alumniID int) (*int, error) {
	var userID *int
	query := `SELECT user_id FROM alumni WHERE id = $1`
	err := db.QueryRow(query, alumniID).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return userID, nil
}

func RestorePekerjaan(db *sql.DB, id int) error {
	now := time.Now()
	query := `UPDATE pekerjaan_alumni SET deleted_at = NULL, deleted_by = NULL, updated_at = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
	result, err := db.Exec(query, now, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func HardDeletePekerjaan(db *sql.DB, id int) error {
	query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
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

func SoftDeletePekerjaanByAlumniID(db *sql.DB, alumniID int, deletedBy int) error {
	now := time.Now()
	query := `UPDATE pekerjaan_alumni SET deleted_at = $1, deleted_by = $2, updated_at = $3 WHERE alumni_id = $4 AND deleted_at IS NULL`

	var deletedByValue interface{}
	if deletedBy > 0 {
		deletedByValue = deletedBy
	} else {
		deletedByValue = nil
	}

	result, err := db.Exec(query, now, deletedByValue, now, alumniID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func HardDeletePekerjaanByAlumniID(db *sql.DB, alumniID int) error {
	query := `DELETE FROM pekerjaan_alumni WHERE alumni_id = $1 AND deleted_at IS NOT NULL`
	result, err := db.Exec(query, alumniID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
