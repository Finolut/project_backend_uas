package model

import (
	"strings"
	"time"
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
	ID                  int        `json:"id" db:"id"`
	AlumniID            int        `json:"alumni_id" db:"alumni_id"`
	NamaPerusahaan      string     `json:"nama_perusahaan" db:"nama_perusahaan"`
	PosisiJabatan       string     `json:"posisi_jabatan" db:"posisi_jabatan"`
	BidangIndustri      string     `json:"bidang_industri" db:"bidang_industri"`
	LokasiKerja         string     `json:"lokasi_kerja" db:"lokasi_kerja"`
	GajiRange           *string    `json:"gaji_range" db:"gaji_range"`
	TanggalMulaiKerja   Date       `json:"tanggal_mulai_kerja" db:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *Date      `json:"tanggal_selesai_kerja" db:"tanggal_selesai_kerja"`
	StatusPekerjaan     string     `json:"status_pekerjaan" db:"status_pekerjaan"`
	DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan" db:"deskripsi_pekerjaan"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy           *int       `json:"deleted_by,omitempty" db:"deleted_by"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

type CreatePekerjaanRequest struct {
	AlumniID            int     `json:"alumni_id" validate:"required"`
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
