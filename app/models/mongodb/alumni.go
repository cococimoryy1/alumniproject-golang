package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struktur utama untuk data Alumni
type Alumni struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"` // gunakan ObjectID untuk MongoDB
	NIM        string             `json:"nim" bson:"nim"`
	Nama       string             `json:"nama" bson:"nama"`
	Jurusan    string             `json:"jurusan" bson:"jurusan"`
	Angkatan   int                `json:"angkatan" bson:"angkatan"`
	TahunLulus int                `json:"tahun_lulus" bson:"tahun_lulus"`
	Email      string             `json:"email" bson:"email"`
	NoTelepon  string             `json:"no_telepon" bson:"no_telepon"`
	Alamat     string             `json:"alamat" bson:"alamat"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// Struktur untuk request pembuatan data Alumni (tanpa ID & timestamp)
type CreateAlumniRequest struct {
	NIM        string `json:"nim" bson:"nim"`
	Nama       string `json:"nama" bson:"nama"`
	Jurusan    string `json:"jurusan" bson:"jurusan"`
	Angkatan   int    `json:"angkatan" bson:"angkatan"`
	TahunLulus int    `json:"tahun_lulus" bson:"tahun_lulus"`
	Email      string `json:"email" bson:"email"`
	NoTelepon  string `json:"no_telepon" bson:"no_telepon"`
	Alamat     string `json:"alamat" bson:"alamat"`
}

// Struktur untuk update data Alumni
type UpdateAlumniRequest struct {
	Nama       string `json:"nama" bson:"nama"`
	Jurusan    string `json:"jurusan" bson:"jurusan"`
	Angkatan   int    `json:"angkatan" bson:"angkatan"`
	TahunLulus int    `json:"tahun_lulus" bson:"tahun_lulus"`
	Email      string `json:"email" bson:"email"`
	NoTelepon  string `json:"no_telepon" bson:"no_telepon"`
	Alamat     string `json:"alamat" bson:"alamat"`
}

// Struktur gabungan alumni + pekerjaan
type AlumniWithPekerjaan struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	NIM        string             `json:"nim" bson:"nim"`
	Nama       string             `json:"nama" bson:"nama"`
	Jurusan    string             `json:"jurusan" bson:"jurusan"`
	Angkatan   int                `json:"angkatan" bson:"angkatan"`
	TahunLulus int                `json:"tahun_lulus" bson:"tahun_lulus"`
	Email      string             `json:"email" bson:"email"`
	NoTelepon  string             `json:"no_telepon" bson:"no_telepon"`
	Alamat     string             `json:"alamat" bson:"alamat"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	Pekerjaan  []Pekerjaan        `json:"pekerjaan" bson:"pekerjaan"`
}
