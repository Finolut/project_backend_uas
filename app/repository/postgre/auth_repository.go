package repository

import (
	"clean-arch/app/model/postgre"
	"database/sql"
)

func GetUserByUsernameOrEmail(db *sql.DB, identifier string) (*model.User, string, error) {
	var user model.User
	var passwordHash string

	query := `SELECT id, username, email, password_hash, role, created_at 
	          FROM users WHERE username = $1 OR email = $1`

	err := db.QueryRow(query, identifier).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash,
		&user.Role, &user.CreatedAt,
	)

	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}

func GetAlumniByNIM(db *sql.DB, nim string) (*model.Alumni, error) {
	var alumni model.Alumni

	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
	          password_hash, role, no_telepon, alamat, created_at, updated_at 
	          FROM alumni WHERE nim = $1`

	err := db.QueryRow(query, nim).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.Password, &alumni.Role, &alumni.NoTelepon,
		&alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &alumni, nil
}

func GetAlumniWithJobs(db *sql.DB, alumniID int) (*model.AlumniWithJobs, error) {
	// Get alumni data
	alumni, err := GetAlumniByID(db, alumniID)
	if err != nil {
		return nil, err
	}

	// Get job history
	jobs, err := GetPekerjaanByAlumniID(db, alumniID)
	if err != nil {
		return nil, err
	}

	return &model.AlumniWithJobs{
		Alumni:        *alumni,
		PekerjaanList: jobs,
	}, nil
}

func CreateAlumniWithAuth(db *sql.DB, req model.CreateAlumniRequest, hashedPassword string) (*model.Alumni, error) {
	var alumni model.Alumni

	query := `INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, 
	          password_hash, role, no_telepon, alamat) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	          RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
	          role, no_telepon, alamat, created_at, updated_at`

	role := "user" // Default role for alumni
	if req.UserID != nil {
		role = "admin" // If user_id is provided, make them admin
	}

	err := db.QueryRow(query,
		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, hashedPassword, role, req.NoTelepon, req.Alamat,
	).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.Role, &alumni.NoTelepon, &alumni.Alamat,
		&alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &alumni, nil
}
