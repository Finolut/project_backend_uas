package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NIM        string             `json:"nim" bson:"nim"`
	Nama       string             `json:"nama" bson:"nama"`
	Jurusan    string             `json:"jurusan" bson:"jurusan"`
	Angkatan   int                `json:"angkatan" bson:"angkatan"`
	TahunLulus int                `json:"tahun_lulus" bson:"tahun_lulus"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"-" bson:"password"`
	Role       string             `json:"role" bson:"role"`
	NoTelepon  *string            `json:"no_telepon" bson:"no_telepon,omitempty"`
	Alamat     *string            `json:"alamat" bson:"alamat,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	DeletedBy  *string            `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
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
