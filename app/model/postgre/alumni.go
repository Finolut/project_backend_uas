package model

import "time"

type Alumni struct {
	ID         int        `json:"id" db:"id"`
	NIM        string     `json:"nim" db:"nim"`
	Nama       string     `json:"nama" db:"nama"`
	Jurusan    string     `json:"jurusan" db:"jurusan"`
	Angkatan   int        `json:"angkatan" db:"angkatan"`
	TahunLulus int        `json:"tahun_lulus" db:"tahun_lulus"`
	Email      string     `json:"email" db:"email"`
	Password   string     `json:"-" db:"password_hash"`
	Role       string     `json:"role" db:"role"`
	NoTelepon  *string    `json:"no_telepon" db:"no_telepon"`
	Alamat     *string    `json:"alamat" db:"alamat"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy  *int       `json:"deleted_by,omitempty" db:"deleted_by"`
}

type AlumniLoginRequest struct {
	NIM      string `json:"nim" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AlumniLoginResponse struct {
	Alumni Alumni `json:"alumni"`
	Token  string `json:"token"`
}

type AlumniWithJobs struct {
	Alumni
	PekerjaanList []PekerjaanAlumni `json:"pekerjaan_list"`
}

type CreateAlumniRequest struct {
	UserID     *int    `json:"user_id"`
	NIM        string  `json:"nim" validate:"required"`
	Nama       string  `json:"nama" validate:"required"`
	Jurusan    string  `json:"jurusan" validate:"required"`
	Angkatan   int     `json:"angkatan" validate:"required"`
	TahunLulus int     `json:"tahun_lulus" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required,min=6"`
	NoTelepon  *string `json:"no_telepon"`
	Alamat     *string `json:"alamat"`
}

type UpdateAlumniRequest struct {
	UserID     *int    `json:"user_id"`
	Nama       string  `json:"nama" validate:"required"`
	Jurusan    string  `json:"jurusan" validate:"required"`
	Angkatan   int     `json:"angkatan" validate:"required"`
	TahunLulus int     `json:"tahun_lulus" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	NoTelepon  *string `json:"no_telepon"`
	Alamat     *string `json:"alamat"`
}

type AlumniStatistics struct {
	TotalAlumni        int            `json:"total_alumni"`
	AlumniByJurusan    map[string]int `json:"alumni_by_jurusan"`
	AlumniByAngkatan   map[string]int `json:"alumni_by_angkatan"`
	AlumniByTahunLulus map[string]int `json:"alumni_by_tahun_lulus"`
}
