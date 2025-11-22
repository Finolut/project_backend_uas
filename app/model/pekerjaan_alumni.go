package model

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Parse date in YYYY-MM-DD format
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte("\"" + d.Time.Format("2006-01-02") + "\""), nil
}

type PekerjaanAlumni struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AlumniID            primitive.ObjectID `json:"alumni_id" bson:"alumni_id"`
	NamaPerusahaan      string             `json:"nama_perusahaan" bson:"nama_perusahaan"`
	PosisiJabatan       string             `json:"posisi_jabatan" bson:"posisi_jabatan"`
	BidangIndustri      string             `json:"bidang_industri" bson:"bidang_industri"`
	LokasiKerja         string             `json:"lokasi_kerja" bson:"lokasi_kerja"`
	GajiRange           *string            `json:"gaji_range" bson:"gaji_range,omitempty"`
	TanggalMulaiKerja   Date               `json:"tanggal_mulai_kerja" bson:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *Date              `json:"tanggal_selesai_kerja" bson:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `json:"status_pekerjaan" bson:"status_pekerjaan"`
	DeskripsiPekerjaan  *string            `json:"deskripsi_pekerjaan" bson:"deskripsi_pekerjaan,omitempty"`
	DeletedAt           *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	DeletedBy           *string            `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreatePekerjaanRequest struct {
	AlumniID            primitive.ObjectID `json:"alumni_id" validate:"required"`
	NamaPerusahaan      string             `json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string             `json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string             `json:"bidang_industri" validate:"required"`
	LokasiKerja         string             `json:"lokasi_kerja" validate:"required"`
	GajiRange           *string            `json:"gaji_range"`
	TanggalMulaiKerja   Date               `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *Date              `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string             `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string            `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanRequest struct {
	NamaPerusahaan      string  `json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string  `json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string  `json:"bidang_industri" validate:"required"`
	LokasiKerja         string  `json:"lokasi_kerja" validate:"required"`
	GajiRange           *string `json:"gaji_range"`
	TanggalMulaiKerja   Date    `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *Date   `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string  `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan"`
}

type SoftDeletePekerjaanRequest struct {
	Reason string `json:"reason,omitempty" validate:"max=255"`
}
