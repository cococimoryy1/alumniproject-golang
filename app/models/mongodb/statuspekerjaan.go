package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AlumniPekerjaan -> gabungan data alumni dan pekerjaan
type AlumniPekerjaan struct {
	ID                      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nama                    string             `json:"nama" bson:"nama"`
	Jurusan                 string             `json:"jurusan" bson:"jurusan"`
	Angkatan                int                `json:"angkatan" bson:"angkatan"`
	BidangIndustri          string             `json:"bidang_industri" bson:"bidang_industri"`
	NamaPerusahaan          string             `json:"nama_perusahaan" bson:"nama_perusahaan"`
	PosisiJabatan           string             `json:"posisi_jabatan" bson:"posisi_jabatan"`
	TanggalMulaiKerja       time.Time          `json:"tanggal_mulai_kerja" bson:"tanggal_mulai_kerja"`
	GajiRange               string             `json:"gaji_range" bson:"gaji_range"`
	StatusPekerjaan         string             `json:"status_pekerjaan" bson:"status_pekerjaan"`
	TotalBekerjaLebih1Tahun int                `json:"total_bekerja_lebih_1_tahun" bson:"total_bekerja_lebih_1_tahun"`
}
