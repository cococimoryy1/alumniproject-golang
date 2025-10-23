package service

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
	"fmt"
	"log"

	"alumniproject/app/models/postgresql"
	"alumniproject/app/repository/postgresql"
	"alumniproject/database/postgresql"
	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaanService(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    role := c.Locals("role").(string)

    list, err := repository.GetAllPekerjaan(role, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(list)
}

func GetPekerjaanByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/pekerjaan/%d", username, id)

	p, err := repository.GetPekerjaanByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    p,
		"message": "Data pekerjaan berhasil diambil",
	})
}

func GetAllAlumniWithPekerjaan(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/alumni-pekerjaan", username)

	list, err := repository.GetAlumniWithPekerjaan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengambil data: " + err.Error(),
		})
	}

	if len(list) == 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"count":   0,
			"data":    []models.AlumniWithPekerjaan{},
			"message": "Belum ada data alumni",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(list),
		"data":    list,
		"message": "Data alumni beserta pekerjaan berhasil diambil",
	})
}

func GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Alumni ID tidak valid"})
	}

	username := c.Locals("username").(string)
	log.Printf("Admin %s mengakses GET /api/pekerjaan/alumni/%d", username, alumniID)

	list, err := repository.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data pekerjaan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    list,
		"message": "Data pekerjaan berhasil diambil",
	})
}

func CreatePekerjaanService(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int) // ambil dari JWT
    var req models.CreatePekerjaanRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    tMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Tanggal mulai tidak valid"})
    }

    var tSelesai *time.Time
    if req.TanggalSelesaiKerja != "" {
        t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Tanggal selesai tidak valid"})
        }
        tSelesai = &t
    }

    pekerjaan := models.Pekerjaan{
        AlumniID:            req.AlumniID,
        NamaPerusahaan:      req.NamaPerusahaan,
        PosisiJabatan:       req.PosisiJabatan,
        BidangIndustri:      req.BidangIndustri,
        LokasiKerja:         req.LokasiKerja,
        GajiRange:           req.GajiRange,
        TanggalMulaiKerja:   tMulai,
        TanggalSelesaiKerja: tSelesai,
        StatusPekerjaan:     req.StatusPekerjaan,
        DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
        CreatedBy:           userID, // otomatis dari JWT
    }

    if err := repository.CreatePekerjaan(&pekerjaan); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan data"})
    }

    return c.JSON(fiber.Map{
        "message": "Pekerjaan berhasil ditambahkan",
        "data":    pekerjaan,
    })
}

func UpdatePekerjaanService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	var req models.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi field wajib
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Field wajib diisi",
		})
	}

	// Ambil data pekerjaan berdasarkan ID
	data, err := repository.GetPekerjaanByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Data tidak ditemukan",
		})
	}

	// Jika bukan admin, pastikan user hanya bisa ubah datanya sendiri
	if role != "admin" && data.CreatedBy != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Kamu tidak punya akses untuk update data ini",
		})
	}

	// Update field
	data.NamaPerusahaan = req.NamaPerusahaan
	data.PosisiJabatan = req.PosisiJabatan
	data.BidangIndustri = req.BidangIndustri
	data.LokasiKerja = req.LokasiKerja
	data.GajiRange = req.GajiRange
	data.StatusPekerjaan = req.StatusPekerjaan
	data.DeskripsiPekerjaan = req.DeskripsiPekerjaan

	// Tanggal mulai kerja
	if req.TanggalMulaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Format tanggal mulai kerja tidak valid",
			})
		}
		data.TanggalMulaiKerja = t
	}

	// Tanggal selesai kerja (optional)
	if req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Format tanggal selesai kerja tidak valid",
			})
		}
		data.TanggalSelesaiKerja = &t
	} else {
		data.TanggalSelesaiKerja = nil
	}

	// Simpan ke database
	if err := repository.UpdatePekerjaan(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memperbarui data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data pekerjaan berhasil diperbarui",
		"data":    data,
	})
}

// func DeletePekerjaanService(c *fiber.Ctx) error {
//     userID := c.Locals("user_id").(int)
//     role := c.Locals("role").(string)
//     id, err := strconv.Atoi(c.Params("id"))
//     if err != nil {
//         return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
//     }

//     err = repository.DeletePekerjaan(id, userID, role)
//     if err != nil {
//         return c.Status(403).JSON(fiber.Map{"error": err.Error()})
//     }

//     return c.JSON(fiber.Map{"message": "Riwayat pekerjaan berhasil dihapus (soft delete)"})
// }

// func RestorePekerjaanService(c *fiber.Ctx) error {
//     id, err := strconv.Atoi(c.Params("id"))
//     if err != nil {
//         return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
//     }

//     err = repository.RestorePekerjaan(id)
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{"error": err.Error()})
//     }

//     return c.JSON(fiber.Map{"message": "Riwayat pekerjaan berhasil direstore"})
// }


func DeletePekerjaanService(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    role := c.Locals("role").(string)
    id, _ := strconv.Atoi(c.Params("id"))

    var res sql.Result
    var err error

    if role == "admin" {
        res, err = postgresql.DB.Exec("UPDATE pekerjaan_alumni SET deleted_at=NOW() WHERE id=$1", id)
    } else {
        res, err = postgresql.DB.Exec("UPDATE pekerjaan_alumni SET deleted_at=NOW() WHERE id=$1 AND created_by=$2", id, userID)
    }

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "DB error"})
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh hapus data ini"})
    }

    return c.JSON(fiber.Map{"message": "Riwayat pekerjaan berhasil dihapus (soft delete)"})
}

func GetTrashPekerjaanService(c *fiber.Ctx) error {
	userRole := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	allTrash, err := repository.GetTrashPekerjaan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data trash"})
	}

		if userRole == "user" {
			var filtered []models.GetTrashPekerjaan
			for _, p := range allTrash {
				if p.CreatedBy == userID {
					filtered = append(filtered, p)
				}
			}
			return c.JSON(fiber.Map{
				"success": true,
				"count":   len(filtered),
				"data":    filtered,
			})
		}


	// Admin dapat semua
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(allTrash),
		"data":    allTrash,
	})
}


func RestorePekerjaanService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	userRole := c.Locals("role")
	userID := c.Locals("user_id").(int)

	// Ambil semua trash dulu
	trash, err := repository.GetTrashPekerjaan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}

	// Pastikan user hanya restore miliknya
	for _, t := range trash {
		if int(t.ID) == id {
			if userRole == "user" && t.CreatedBy != userID {
				return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh mengakses data milik orang lain"})
			}
			err = repository.RestorePekerjaan(id)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			return c.JSON(fiber.Map{"message": "Data berhasil direstore"})
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
}

func HardDeletePekerjaanService(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	userRole := c.Locals("role")
	userID := c.Locals("user_id").(int)

	// Ambil semua trash
	trash, err := repository.GetTrashPekerjaan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}

for _, t := range trash {
    if int64(t.ID) == id { // konversi t.ID ke int64
        if userRole == "user" && t.CreatedBy != userID {
            return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh menghapus data milik orang lain"})
        }
        err = repository.HardDeletePekerjaanByID(id)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"message": fmt.Sprintf("Data dengan ID %d dihapus permanen", id)})
    }
}

	return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
}



// GetPekerjaanPaginated -> ambil data pekerjaan dengan pagination, sorting, dan search
func GetPekerjaanPaginated(pageStr, limitStr, sortBy, order, search string) (*models.PekerjaanResponse, error) {
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 100 { // Batasi max 100
		limit = 5 // Default limit menjadi 5
	}

	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{
		"id":                   true,
		"nama_perusahaan":      true,
		"posisi_jabatan":       true,
		"tanggal_mulai_kerja": true,
		"created_at":           true,
	}
	if _, ok := sortByWhitelist[sortBy]; !ok {
		sortBy = "id"
	}

	order = strings.ToLower(order)
	if order != "desc" {
		order = "asc"
	}

	pekerjaan, err := repository.GetPekerjaanPaginated(search, sortBy, order, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := repository.CountPekerjaan(search)
	if err != nil {
		return nil, err
	}

	response := &models.PekerjaanResponse{
		Data: pekerjaan,
		Meta: &models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return response, nil
}

func SoftDeletePekerjaan(id int, userID int, role string) error {
    return repository.SoftDeletePekerjaan(id, userID, role)
}
