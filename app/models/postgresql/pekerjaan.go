package models

import "time"

type Pekerjaan struct {
    ID                  int64      `json:"id"`
    AlumniID            int        `json:"alumni_id"`
    NamaPerusahaan      string     `json:"nama_perusahaan"`
    PosisiJabatan       string     `json:"posisi_jabatan"`
    BidangIndustri      string     `json:"bidang_industri"`
    LokasiKerja         string     `json:"lokasi_kerja"`
    GajiRange           string     `json:"gaji_range"`
    TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja,omitempty"`
    StatusPekerjaan     string     `json:"status_pekerjaan"`
    DeskripsiPekerjaan  string     `json:"deskripsi_pekerjaan"`
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
    CreatedBy           int        `json:"created_by"`  // siapa yang input
    DeletedAt           *time.Time `json:"deleted_at,omitempty"` // soft delete
}

type CreatePekerjaanRequest struct {
    AlumniID            int    `json:"alumni_id"`
    NamaPerusahaan      string `json:"nama_perusahaan"`
    PosisiJabatan       string `json:"posisi_jabatan"`
    BidangIndustri      string `json:"bidang_industri"`
    LokasiKerja         string `json:"lokasi_kerja"`
    GajiRange           string `json:"gaji_range"`
    TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"` // YYYY-MM-DD
    TanggalSelesaiKerja string `json:"tanggal_selesai_kerja,omitempty"`
    StatusPekerjaan     string `json:"status_pekerjaan"`
    DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanRequest struct {
	NamaPerusahaan        string `json:"nama_perusahaan"`
	PosisiJabatan         string `json:"posisi_jabatan"`
	BidangIndustri        string `json:"bidang_industri"`
	LokasiKerja           string `json:"lokasi_kerja"`
	GajiRange             string `json:"gaji_range"`
	TanggalMulaiKerja     string `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja   string `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan       string `json:"status_pekerjaan"`
	DeskripsiPekerjaan    string `json:"deskripsi_pekerjaan"`
}

type PekerjaanWithAlumni struct {
	ID                  int       `json:"id"`
	AlumniID            int       `json:"alumni_id"`
	NamaPerusahaan      string    `json:"nama_perusahaan"`
	PosisiJabatan       string    `json:"posisi_jabatan"`
	BidangIndustri      string    `json:"bidang_industri"`
	LokasiKerja         string    `json:"lokasi_kerja"`
	GajiRange           string    `json:"gaji_range"`
	TanggalMulaiKerja   time.Time `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string    `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string    `json:"deskripsi_pekerjaan"`
	Alumni              struct {
		Nama    string `json:"nama"`
		Jurusan string `json:"jurusan"`
		Email   string `json:"email"`
	} `json:"alumni"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TrashPekerjaan struct {
	ID            int        `json:"id"`
	NamaPekerjaan string     `json:"nama_pekerjaan"`
	TempatKerja   string     `json:"tempat_kerja"`
	DeletedAt     *time.Time `json:"deleted_at"`
	
}

type GetTrashPekerjaan struct {
	ID              int        `json:"id"`
	AlumniID        int        `json:"alumni_id"`
	NamaPerusahaan  string     `json:"nama_perusahaan"`
	PosisiJabatan   string     `json:"posisi_jabatan"`
	BidangIndustri  string     `json:"bidang_industri"`
	LokasiKerja     string     `json:"lokasi_kerja"`
	StatusPekerjaan string     `json:"status_pekerjaan"`
	DeletedAt       *time.Time `json:"deleted_at"`
	CreatedBy       int        `json:"created_by"`
}