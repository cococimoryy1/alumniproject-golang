package repository

import (
	"context"
	"time"
	"log"
	"fmt"
	"database/sql"
	
	"alumniproject/database/postgresql"
	"alumniproject/app/models/postgresql"
)


func GetAllAlumni(role string, userID int) ([]models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseQuery := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at, created_by
		FROM alumni
	`

	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = postgresql.DB.QueryContext(ctx, baseQuery + " ORDER BY created_at DESC")
	} else {
		rows, err = postgresql.DB.QueryContext(ctx, baseQuery + " WHERE created_by = $1 ORDER BY created_at DESC", userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Alumni
	for rows.Next() {
		var a models.Alumni
		err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
			&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt, &a.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}


func GetAlumniByID(id int) (models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var a models.Alumni
	query := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, 
		       created_at, updated_at, created_by, deleted_at
		FROM alumni
		WHERE id = $1
	`

	err := postgresql.DB.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.NIM,
		&a.Nama,
		&a.Jurusan,
		&a.Angkatan,
		&a.TahunLulus,
		&a.Email,
		&a.NoTelepon,
		&a.Alamat,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.CreatedBy,
		&a.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return a, fmt.Errorf("data alumni tidak ditemukan")
		}
		return a, err
	}

	return a, nil
}


// CreateAlumni sama saja, tambahkan created_by
func CreateAlumni(a *models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return postgresql.DB.QueryRowContext(ctx, `
		INSERT INTO alumni (
			nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_by, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id
	`, a.NIM, a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus, a.Email, a.NoTelepon, a.Alamat, a.CreatedBy, time.Now(), time.Now()).Scan(&a.ID)
}

// UpdateAlumni dengan role-based
func UpdateAlumni(a *models.Alumni, userID int, role string) error {
    var query string
    var args []interface{}

    if role == "admin" {
        query = `
            UPDATE alumni 
            SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, email=$5, no_telepon=$6, alamat=$7, updated_at=$8
            WHERE id=$9`
        args = []interface{}{a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus, a.Email, a.NoTelepon, a.Alamat, a.UpdatedAt, a.ID}
    } else {
        // hanya boleh update data miliknya sendiri
        query = `
            UPDATE alumni 
            SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, email=$5, no_telepon=$6, alamat=$7, updated_at=$8
            WHERE id=$9 AND created_by=$10`
        args = []interface{}{a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus, a.Email, a.NoTelepon, a.Alamat, a.UpdatedAt, a.ID, userID}
    }

    _, err := postgresql.DB.Exec(query, args...)
    return err
}


// DeleteAlumni (soft delete)
func DeleteAlumni(id int, userID int, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res sql.Result
	var err error

	if role == "admin" {
		res, err = postgresql.DB.ExecContext(ctx, `
			UPDATE alumni SET deleted_at=$1 WHERE id=$2
		`, time.Now(), id)
	} else {
		res, err = postgresql.DB.ExecContext(ctx, `
			UPDATE alumni SET deleted_at=$1 WHERE id=$2 AND created_by=$3
		`, time.Now(), id, userID)
	}

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak boleh hapus data ini")
	}
	return nil
}

func RestoreAlumni(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := postgresql.DB.ExecContext(ctx, `
		UPDATE alumni 
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


// func DeleteAlumni(id int, userID int, role string) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     // 1. Salin data ke tabel alumni_archive
//     _, err := database.DB.ExecContext(ctx, `
//         INSERT INTO alumni_archive
//         (id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at, created_by, deleted_at)
//         SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at, created_by, NOW()
//         FROM alumni
//         WHERE id = $1
//     `, id)
//     if err != nil {
//         return err
//     }

//     // 2. Kolom yang ingin dihapus
//     colsToNull := `
//         nim = NULL,
//         nama = NULL,
//         jurusan = NULL,
//         angkatan = NULL,
//         tahun_lulus = NULL,
//         email = NULL,
//         no_telepon = NULL,
//         alamat = NULL,
//         deleted_at = $1
//     `

//     var res sql.Result

//     if role == "admin" {
//         // Admin bisa hapus semua data
//         res, err = database.DB.ExecContext(ctx, `
//             UPDATE alumni SET `+colsToNull+` WHERE id = $2
//         `, time.Now(), id)
//     } else {
//         // User hanya bisa hapus data miliknya sendiri
//         res, err = database.DB.ExecContext(ctx, `
//             UPDATE alumni SET `+colsToNull+` WHERE id = $2 AND created_by = $3
//         `, time.Now(), id, userID)
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


// func RestoreAlumni(id int) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     // Update tabel alumni dari archive
//     res, err := database.DB.ExecContext(ctx, `
//         UPDATE alumni a
//         SET nim = b.nim,
//             nama = b.nama,
//             jurusan = b.jurusan,
//             angkatan = b.angkatan,
//             tahun_lulus = b.tahun_lulus,
//             email = b.email,
//             no_telepon = b.no_telepon,
//             alamat = b.alamat,
//             created_at = b.created_at,
//             updated_at = b.updated_at,
//             created_by = b.created_by,
//             deleted_at = NULL
//         FROM alumni_archive b
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


// GetAlumniPaginated -> ambil data alumni dengan search, sort, paginate (adapt dari GetUsersRepo)
func GetAlumniRepo(search, sortBy, order string, limit, offset int) ([]models.Alumni, error) {
    query := fmt.Sprintf(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := postgresql.DB.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        log.Println("Query error:", err)
        return nil, err
    }
    defer rows.Close()

    var alumni []models.Alumni
    for rows.Next() {
        var a models.Alumni
        if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
            return nil, err
        }
        alumni = append(alumni, a)
    }
    return alumni, nil
}

func CountAlumniRepo(search string) (int, error) {
    var total int
    query := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1`
    err := postgresql.DB.QueryRow(query, "%"+search+"%").Scan(&total)
    if err != nil && err != sql.ErrNoRows {
        return 0, err
    }
    return total, nil
}