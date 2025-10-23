package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	// "go.mongodb.org/mongo-driver/mongo"
	"alumniproject/app/models/mongodb"
	"alumniproject/app/repository/mongodb"
)



func GetAllPekerjaanService(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
    userID, ok := c.Locals("user_id").(int)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak ditemukan"})
    }

    role, ok := c.Locals("role").(string)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Role tidak ditemukan"})
    }

    fmt.Printf("👤 Authenticated user: %d (role: %s)\n", userID, role)

    list, err := repo.GetAllPekerjaan(role, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    list,
        "count":   len(list),
    })
}

func GetPekerjaanByID(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
    id := c.Params("id")
    username, _ := c.Locals("username").(string)
    username = "anonymous" // Default jika error
    userID, _ := c.Locals("user_id").(int)
    role, _ := c.Locals("role").(string)

    log.Printf("User %s (ID: %d, role: %s) mengakses GET /api/pekerjaan/%s", username, userID, role, id)

    p, err := repo.GetPekerjaanByID(id, role, userID) // ✅ Pass role/userID
    if err != nil {
        if err.Error() == "mongo: no documents in result" {
            return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
        }
        log.Printf("Error fetching pekerjaan ID %s: %v", id, err)
        return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    p,
        "message": "Data pekerjaan berhasil diambil",
    })
}

func GetAllAlumniWithPekerjaan(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
	username := c.Locals("username").(string)
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)
	log.Printf("User %s (ID: %d, role: %s) mengakses GET /api/pekerjaan/alumni-pekerjaan", username, userID, role)

	isAdmin := role == "admin"
	list, err := repo.GetAlumniWithPekerjaan(userID, isAdmin)
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
	repo := repository.New() // ✅ pindahkan ke sini
	alumniIDStr := c.Params("alumni_id")
	alumniID, err := strconv.Atoi(alumniIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Alumni ID tidak valid"})
	}

	username := c.Locals("username").(string)
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)
	log.Printf("User %s (ID: %d, role: %s) mengakses GET /api/pekerjaan/alumni/%d", username, userID, role, alumniID)

	// ✅ Check admin (karena route sudah AdminOnly, tapi double-check)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Hanya admin yang boleh akses"})
	}

	list, err := repo.GetPekerjaanByAlumniID(alumniID, role, userID)
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
	repo := repository.New() // ✅ pindahkan ke sini
	userID := c.Locals("user_id").(int)
	// ✅ Hapus: role := c.Locals("role").(string) – unused
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

	// ✅ Buat struct Pekerjaan dari req, set CreatedBy
	p := &models.Pekerjaan{
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
		CreatedBy:           userID, // ✅ Set dari user
	}

	if err := repo.CreatePekerjaan(p); err != nil {
		log.Printf("❌ Create error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat pekerjaan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    p,
		"message": "Pekerjaan berhasil dibuat",
	})
}
func UpdatePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
	id := c.Params("id")
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)
	var req models.UpdatePekerjaanRequest

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

	// ✅ Buat partial update struct
	p := &models.Pekerjaan{
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tMulai,
		TanggalSelesaiKerja: tSelesai,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		// CreatedBy tidak diubah
	}

	if err := repo.UpdatePekerjaan(id, p, role, userID); err != nil {
		if err.Error() == "mongo: no documents in result" || err.Error() == "invalid ID format: " {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update pekerjaan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil diupdate",
	})
}

// DeletePekerjaanService: Soft delete
func DeletePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
	id := c.Params("id")
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	err := repo.SoftDeletePekerjaan(id, userID, role)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Riwayat pekerjaan berhasil dihapus (soft delete)",
	})
}

func GetTrashPekerjaanService(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	allTrash, err := repo.GetTrashPekerjaan(userID, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil trash"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(allTrash),
		"data":    allTrash,
	})
}

func RestorePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
    idStr := c.Params("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }

    userID := c.Locals("user_id").(int)
    userRole := c.Locals("role").(string)

    trash, err := repo.GetTrashPekerjaan(userID, userRole)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
    }

    for _, t := range trash {
        if t.ID == id {
            if userRole != "admin" && t.CreatedBy != userID { // ✅ Check owner
                return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh mengakses data milik orang lain"})
            }

            err = repo.RestorePekerjaan(id)
            if err != nil {
                return c.Status(500).JSON(fiber.Map{"error": err.Error()})
            }

            return c.JSON(fiber.Map{"message": "Data berhasil direstore"})
        }
    }

    return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
}

func HardDeletePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
    idStr := c.Params("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }

    userID := c.Locals("user_id").(int)
    userRole := c.Locals("role").(string)

    trash, err := repo.GetTrashPekerjaan(userID, userRole)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
    }

    for _, t := range trash {
        if t.ID == id {
            if userRole != "admin" && t.CreatedBy != userID { // ✅ Check owner
                return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh menghapus data milik orang lain"})
            }

            err = repo.HardDeletePekerjaanByID(id)
            if err != nil {
                return c.Status(500).JSON(fiber.Map{"error": err.Error()})
            }

            return c.JSON(fiber.Map{"message": fmt.Sprintf("Data dengan ID %d dihapus permanen", id)})
        }
    }

    return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
}

func GetPekerjaanPaginated(c *fiber.Ctx) error {
	repo := repository.New() // ✅ pindahkan ke sini
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "5")
	sortBy := c.Query("sort_by", "created_at")
	order := c.Query("order", "desc")
	search := c.Query("search", "")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 100 {
		limit = 5
	}

	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{
		"id":                  true,
		"nama_perusahaan":     true,
		"posisi_jabatan":      true,
		"tanggal_mulai_kerja": true,
		"created_at":          true,
	}
	if _, ok := sortByWhitelist[sortBy]; !ok {
		sortBy = "id"
	}

	order = strings.ToLower(order)
	if order != "desc" {
		order = "asc"
	}

	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	pekerjaan, err := repo.GetPekerjaanPaginated(search, sortBy, order, limit, offset, role, userID) // ✅ Pass role/userID
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	total64, err := repo.CountPekerjaan(search, role, userID) // ✅ Pass
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	total := int(total64)

	response := &models.PekerjaanResponse{ // Asumsi struct ada
		Data: pekerjaan,
		Meta: &models.MetaInfo{ // Asumsi struct ada
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Login: Asumsi sudah ada di service, tapi jika perlu tambah
// func Login(c *fiber.Ctx) error { ... } // Implement jika belum