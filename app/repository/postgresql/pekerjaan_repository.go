package repository

import (
	"context"
	"time"
	"database/sql"
	"fmt"
	"log"

	"alumniproject/database/postgresql"
	"alumniproject/app/models/postgresql"
)

func GetAllPekerjaan(role string, userID int) ([]models.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseQuery := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
		       created_by, created_at, updated_at
		FROM pekerjaan_alumni
	`

	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = postgresql.DB.QueryContext(ctx, baseQuery+" ORDER BY created_at DESC")
	} else {
		// Filter berdasarkan JWT userID
		rows, err = postgresql.DB.QueryContext(ctx, baseQuery+" WHERE created_by = $1 ORDER BY created_at DESC", userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Pekerjaan
	for rows.Next() {
		var p models.Pekerjaan
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
			&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja,
			&p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}


func GetPekerjaanByID(id int) (models.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p models.Pekerjaan
	err := postgresql.DB.QueryRowContext(ctx, `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE id = $1
	`, id).Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func GetPekerjaanByAlumniID(alumniID int) ([]models.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := postgresql.DB.QueryContext(ctx, `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE alumni_id = $1 ORDER BY created_at DESC
	`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Pekerjaan
	for rows.Next() {
		var p models.Pekerjaan
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func CreatePekerjaan(p *models.Pekerjaan) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := postgresql.DB.QueryRowContext(ctx, `
        INSERT INTO pekerjaan_alumni (
            alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
            gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
            deskripsi_pekerjaan, created_by, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id
    `,
        p.AlumniID,
        p.NamaPerusahaan,
        p.PosisiJabatan,
        p.BidangIndustri,
        p.LokasiKerja,
        p.GajiRange,
        p.TanggalMulaiKerja,
        p.TanggalSelesaiKerja,
        p.StatusPekerjaan,
        p.DeskripsiPekerjaan,
        p.CreatedBy,   // ⬅️ penting
        time.Now(),
        time.Now(),
    ).Scan(&p.ID)

    return err
}


func UpdatePekerjaan(p *models.Pekerjaan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postgresql.DB.ExecContext(ctx, `
		UPDATE pekerjaan_alumni SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3, lokasi_kerja = $4, gaji_range = $5, tanggal_mulai_kerja = $6, tanggal_selesai_kerja = $7, status_pekerjaan = $8, deskripsi_pekerjaan = $9, updated_at = $10
		WHERE id = $11
	`, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, time.Now(), p.ID)
	return err
}

// Backup & Soft Delete pekerjaan_alumni
// func DeletePekerjaan(id int, userID int, role string) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     // 1. Insert data ke archive sebelum dihapus (pakai query baru)
//     _, err := database.DB.ExecContext(ctx, `
//         INSERT INTO pekerjaan_alumni_archive
//         (id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
//          tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
//          created_at, updated_at, created_by, deleted_at)
//         SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
//                tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
//                created_at, updated_at, created_by, NOW()
//         FROM pekerjaan_alumni
//         WHERE id = $1
//     `, id)
//     if err != nil {
//         return err
//     }

//     // 2. Kolom yang ingin dihapus (NULL)
//     colsToNull := `
//         nama_perusahaan = NULL,
//         posisi_jabatan = NULL,
//         bidang_industri = NULL,
//         lokasi_kerja = NULL,
//         gaji_range = NULL,
//         tanggal_mulai_kerja = NULL,
//         tanggal_selesai_kerja = NULL,
//         status_pekerjaan = NULL,
//         deskripsi_pekerjaan = NULL,
//         deleted_at = NOW()
//     `

//     var res sql.Result
//     if role == "admin" {
//         res, err = database.DB.ExecContext(ctx, `UPDATE pekerjaan_alumni SET `+colsToNull+` WHERE id=$1`, id)
//     } else {
//         res, err = database.DB.ExecContext(ctx, `UPDATE pekerjaan_alumni SET `+colsToNull+` WHERE id=$1 AND created_by=$2`, id, userID)
//     }

//     if err != nil {
//         return err
//     }

//     rows, _ := res.RowsAffected()
//     if rows == 0 {
//         return fmt.Errorf("tidak boleh hapus data ini atau data tidak ditemukan")
//     }

//     return nil
// }


// Restore pekerjaan_alumni dari archive
// func RestorePekerjaan(id int) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     res, err := database.DB.ExecContext(ctx, `
//         UPDATE pekerjaan_alumni a
//         SET 
//             alumni_id = b.alumni_id,
//             nama_perusahaan = b.nama_perusahaan,
//             posisi_jabatan = b.posisi_jabatan,
//             bidang_industri = b.bidang_industri,
//             lokasi_kerja = b.lokasi_kerja,
//             gaji_range = b.gaji_range,
//             tanggal_mulai_kerja = b.tanggal_mulai_kerja,
//             tanggal_selesai_kerja = b.tanggal_selesai_kerja,
//             status_pekerjaan = b.status_pekerjaan,
//             deskripsi_pekerjaan = b.deskripsi_pekerjaan,
//             created_at = b.created_at,
//             updated_at = b.updated_at,
//             created_by = b.created_by,
//             deleted_at = NULL
//         FROM pekerjaan_alumni_archive b
//         WHERE a.id = b.id AND a.id = $1
//     `, id)

//     if err != nil {
//         return err
//     }

//     rows, _ := res.RowsAffected()
//     if rows == 0 {
//         return fmt.Errorf("data tidak ditemukan di archive atau gagal di-restore")
//     }

//     return nil
// }


func GetTrashPekerjaan() ([]models.GetTrashPekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		       lokasi_kerja, status_pekerjaan, deleted_at, created_by
		FROM pekerjaan_alumni
		WHERE deleted_at IS NOT NULL
		ORDER BY deleted_at DESC;
	`

	rows, err := postgresql.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.GetTrashPekerjaan
	for rows.Next() {
		var t models.GetTrashPekerjaan
		if err := rows.Scan(
			&t.ID,
			&t.AlumniID,
			&t.NamaPerusahaan,
			&t.PosisiJabatan,
			&t.BidangIndustri,
			&t.LokasiKerja,
			&t.StatusPekerjaan,
			&t.DeletedAt,
			&t.CreatedBy, // pastikan ada di SELECT
		); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func HardDeletePekerjaanByID(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM pekerjaan_alumni
		WHERE id = $1 
		AND deleted_at IS NOT NULL
	`

	result, err := postgresql.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("data tidak ditemukan atau belum dihapus (soft delete)")
	}

	return nil
}


func RestorePekerjaan(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    res, err := postgresql.DB.ExecContext(ctx, `
        UPDATE pekerjaan_alumni 
        SET deleted_at = NULL 
        WHERE id = $1
    `, id)
    if err != nil {
        return err
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("data tidak ditemukan atau belum dihapus")
    }

    return nil
}





// Tampilkan data yang sudah di-soft delete
// func GetTrashedBarang(db *sql.DB) ([]models.Barang, error) {
// 	rows, err := db.Query("SELECT id, nama, deskripsi, deleted_at FROM barang WHERE deleted_at IS NOT NULL")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []models.Barang
// 	for rows.Next() {
// 		var b models.Barang
// 		err := rows.Scan(&b.ID, &b.Nama, &b.Deskripsi, &b.DeletedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result = append(result, b)
// 	}
// 	return result, nil
// }

// Hard delete
func HardDeleteBarang(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM barang WHERE id = ? AND deleted_at IS NOT NULL", id)
	return err
}

func GetAlumniWithPekerjaan() ([]models.AlumniWithPekerjaan, error) {
    rows, err := postgresql.DB.Query(`
		SELECT 
			a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.no_telepon, a.alamat, 
			a.created_at AS alumni_created_at, a.updated_at AS alumni_updated_at,
			p.id AS pekerjaan_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri, p.lokasi_kerja, p.gaji_range, 
			p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, p.status_pekerjaan, p.deskripsi_pekerjaan, 
			p.created_at AS pekerjaan_created_at, p.updated_at AS pekerjaan_updated_at
		FROM alumni a
		LEFT JOIN pekerjaan_alumni p ON a.id = p.alumni_id
		ORDER BY a.id, p.tanggal_mulai_kerja DESC;

    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    alumniMap := make(map[int]*models.AlumniWithPekerjaan)

    for rows.Next() {
        var a models.Alumni
        var (
            pekerjaanID sql.NullInt64
            namaPerusahaan, posisiJabatan, bidangIndustri, lokasiKerja, gajiRange, statusPekerjaan, deskripsiPekerjaan sql.NullString
            tanggalMulai, tanggalSelesai sql.NullTime
            pCreatedAt, pUpdatedAt sql.NullTime
        )

		err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt,
			&pekerjaanID, &namaPerusahaan, &posisiJabatan, &bidangIndustri, &lokasiKerja, &gajiRange,
			&tanggalMulai, &tanggalSelesai, &statusPekerjaan, &deskripsiPekerjaan,
			&pCreatedAt, &pUpdatedAt,
		)
        if err != nil {
            return nil, err
        }

        if _, ok := alumniMap[a.ID]; !ok {
            alumniMap[a.ID] = &models.AlumniWithPekerjaan{
                ID:         a.ID,
                NIM:        a.NIM,
                Nama:       a.Nama,
                Jurusan:    a.Jurusan,
                Angkatan:   a.Angkatan,
                TahunLulus: a.TahunLulus,
                Email:      a.Email,
                NoTelepon:  a.NoTelepon,
                Alamat:     a.Alamat,
                CreatedAt:  a.CreatedAt,
                UpdatedAt:  a.UpdatedAt,
                Pekerjaan:  []models.Pekerjaan{},
            }
        }

        if pekerjaanID.Valid {
            pekerjaan := models.Pekerjaan{
                ID:               pekerjaanID.Int64, // ✅ langsung int64
                AlumniID:        a.ID,
                NamaPerusahaan:  namaPerusahaan.String,
                PosisiJabatan:   posisiJabatan.String,
                BidangIndustri:  bidangIndustri.String,
                LokasiKerja:     lokasiKerja.String,
                GajiRange:       gajiRange.String,
                TanggalMulaiKerja: tanggalMulai.Time,
                StatusPekerjaan: statusPekerjaan.String,
                DeskripsiPekerjaan: deskripsiPekerjaan.String,
                CreatedAt:       pCreatedAt.Time,
                UpdatedAt:       pUpdatedAt.Time,
            }
            if tanggalSelesai.Valid {
                pekerjaan.TanggalSelesaiKerja = &tanggalSelesai.Time
            }
            alumniMap[a.ID].Pekerjaan = append(alumniMap[a.ID].Pekerjaan, pekerjaan)
        }
    }

    var results []models.AlumniWithPekerjaan
    for _, v := range alumniMap {
        results = append(results, *v)
    }

    return results, nil
}



// FIX: Tambah helper functions untuk dereference pointers
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// GetPekerjaanPaginated -> ambil data pekerjaan dengan search, sort, paginate
func GetPekerjaanPaginated(search, sortBy, order string, limit, offset int) ([]models.Pekerjaan, error) {
    query := fmt.Sprintf(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
               tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
        FROM pekerjaan_alumni
        WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rows, err := postgresql.DB.QueryContext(ctx, query, "%"+search+"%", limit, offset)
    if err != nil {
        log.Println("Query error:", err)
        return nil, err
    }
    defer rows.Close()

    var pekerjaan []models.Pekerjaan
    for rows.Next() {
        var p models.Pekerjaan
        var tanggalSelesai sql.NullTime
        err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, 
            &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &tanggalSelesai, &p.StatusPekerjaan, 
            &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
        if err != nil {
            return nil, err
        }
        if tanggalSelesai.Valid {
            p.TanggalSelesaiKerja = &tanggalSelesai.Time
        }
        pekerjaan = append(pekerjaan, p)
    }
    return pekerjaan, nil
}

// CountPekerjaan -> hitung total
func CountPekerjaan(search string) (int, error) {
    var total int
    query := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1`
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err := postgresql.DB.QueryRowContext(ctx, query, "%"+search+"%").Scan(&total)
    return total, err
}

func SoftDeletePekerjaan(id int, userID int, role string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var res sql.Result
    var err error

    if role == "admin" {
        // Admin boleh delete semua
        res, err = postgresql.DB.ExecContext(ctx,
            "UPDATE pekerjaan_alumni SET deleted_at = NOW() WHERE id=$1", id)
    } else {
        // User hanya boleh delete miliknya sendiri
        res, err = postgresql.DB.ExecContext(ctx,
            "UPDATE pekerjaan_alumni SET deleted_at = NOW() WHERE id=$1 AND created_by=$2", id, userID)
    }

    if err != nil {
        return err
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("tidak ada data yang dihapus, mungkin bukan milik user ini")
    }

    return nil
}
