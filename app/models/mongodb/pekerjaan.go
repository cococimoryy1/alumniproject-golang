package models

import (
    "time"
)

type Pekerjaan struct {
    ID                  int64             `json:"id" bson:"_id,omitempty"`
    AlumniID            int               `json:"alumni_id" bson:"alumni_id"`
    NamaPerusahaan      string            `json:"nama_perusahaan" bson:"nama_perusahaan"`
    PosisiJabatan       string            `json:"posisi_jabatan" bson:"posisi_jabatan"`
    BidangIndustri      string            `json:"bidang_industri" bson:"bidang_industri"`
    LokasiKerja         string            `json:"lokasi_kerja" bson:"lokasi_kerja"`
    GajiRange           string            `json:"gaji_range" bson:"gaji_range"`
    TanggalMulaiKerja   time.Time         `json:"tanggal_mulai_kerja" bson:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *time.Time        `json:"tanggal_selesai_kerja,omitempty" bson:"tanggal_selesai_kerja,omitempty"`
    StatusPekerjaan     string            `json:"status_pekerjaan" bson:"status_pekerjaan"`
    DeskripsiPekerjaan  string            `json:"deskripsi_pekerjaan" bson:"deskripsi_pekerjaan"`
    CreatedAt           time.Time         `json:"created_at" bson:"created_at"`
    UpdatedAt           time.Time         `json:"updated_at" bson:"updated_at"`
    CreatedBy           int               `json:"created_by" bson:"created_by"`
    DeletedAt           *time.Time        `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type CreatePekerjaanRequest struct {
    AlumniID            int    `json:"alumni_id"`
    NamaPerusahaan      string `json:"nama_perusahaan"`
    PosisiJabatan       string `json:"posisi_jabatan"`
    BidangIndustri      string `json:"bidang_industri"`
    LokasiKerja         string `json:"lokasi_kerja"`
    GajiRange           string `json:"gaji_range"`
    TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja string `json:"tanggal_selesai_kerja,omitempty"`
    StatusPekerjaan     string `json:"status_pekerjaan"`
    DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanRequest struct {
    NamaPerusahaan      string `json:"nama_perusahaan"`
    PosisiJabatan       string `json:"posisi_jabatan"`
    BidangIndustri      string `json:"bidang_industri"`
    LokasiKerja         string `json:"lokasi_kerja"`
    GajiRange           string `json:"gaji_range"`
    TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja string `json:"tanggal_selesai_kerja,omitempty"`
    StatusPekerjaan     string `json:"status_pekerjaan"`
    DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
}

type GetTrashPekerjaan struct {
    ID             int64     `json:"id" bson:"_id,omitempty"` // ✅ Ubah ke _id untuk konsistensi
    AlumniID       int       `json:"alumni_id" bson:"alumni_id"`
    NamaPerusahaan string    `json:"nama_perusahaan" bson:"nama_perusahaan"`
    PosisiJabatan  string    `json:"posisi_jabatan" bson:"posisi_jabatan"`
    BidangIndustri string    `json:"bidang_industri" bson:"bidang_industri"`
    LokasiKerja    string    `json:"lokasi_kerja" bson:"lokasi_kerja"`
    StatusPekerjaan string   `json:"status_pekerjaan" bson:"status_pekerjaan"`
    DeletedAt      *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
    CreatedBy      int       `json:"created_by" bson:"created_by"`
}