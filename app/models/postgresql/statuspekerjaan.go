package models

import "time"

// AlumniPekerjaan represents the combined data from alumni and pekerjaan_alumni tables for the specific query
type AlumniPekerjaan struct {
	ID                      int       `json:"id"`
	Nama                    string    `json:"nama"`
	Jurusan                 string    `json:"jurusan"`
	Angkatan                int       `json:"angkatan"`
	BidangIndustri          string    `json:"bidang_industri"`
	NamaPerusahaan          string    `json:"nama_perusahaan"`
	PosisiJabatan           string    `json:"posisi_jabatan"`
	TanggalMulaiKerja       time.Time `json:"tanggal_mulai_kerja"`
	GajiRange               string    `json:"gaji_range"`
	StatusPekerjaan         string    `json:"status_pekerjaan"`
	TotalBekerjaLebih1Tahun int       `json:"total_bekerja_lebih_1_tahun"`
}

