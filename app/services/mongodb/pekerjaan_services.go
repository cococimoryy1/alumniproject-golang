package service

import (
	"fmt"
	"log"
	"strings"
	"time"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"alumniproject/app/models/mongodb"
	"alumniproject/app/repository/mongodb"
)
// GetAllPekerjaanService godoc
// @Summary Menampilkan semua data pekerjaan
// @Description Mengambil semua data pekerjaan dari MongoDB (dapat difilter dan diurutkan)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param search query string false "Cari berdasarkan nama perusahaan atau posisi"
// @Param sort_by query string false "Urutkan berdasarkan kolom (contoh: nama_perusahaan, tanggal_mulai_kerja)"
// @Param order query string false "Urutan data (asc/desc)"
// @Param limit query int false "Jumlah maksimum data (default: 10)"
// @Success 200 {array} models.Pekerjaan
// @Router /api/pekerjaan [get]
func GetAllPekerjaanService(c *fiber.Ctx) error {
    // implementasi asli kamu


	repo := repository.New()
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

	var filtered []map[string]interface{}
	for _, p := range list {
		filtered = append(filtered, map[string]interface{}{
			"id":                    p.ID.Hex(),
			"alumni_id":             p.AlumniID,
			"nama_perusahaan":       p.NamaPerusahaan,
			"posisi_jabatan":        p.PosisiJabatan,
			"bidang_industri":       p.BidangIndustri,
			"lokasi_kerja":          p.LokasiKerja,
			"gaji_range":            p.GajiRange,
			"tanggal_mulai_kerja":   p.TanggalMulaiKerja,
			"tanggal_selesai_kerja": p.TanggalSelesaiKerja,
			"status_pekerjaan":      p.StatusPekerjaan,
			"deskripsi_pekerjaan":   p.DeskripsiPekerjaan,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    filtered,
		"count":   len(filtered),
	})
}


// GetPekerjaanByID godoc
// @Summary Menampilkan detail pekerjaan berdasarkan ID
// @Description Mengambil satu data pekerjaan berdasarkan ID dari MongoDB
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Param include_deleted query bool false "Tampilkan juga jika pekerjaan sudah dihapus (soft delete)"
// @Success 200 {object} models.Pekerjaan
// @Failure 404 {object} map[string]string
// @Router /api/pekerjaan/{id} [get]
func GetPekerjaanByID(c *fiber.Ctx) error {
	repo := repository.New()
	id := c.Params("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	username, _ := c.Locals("username").(string)
	userID, _ := c.Locals("user_id").(int)
	role, _ := c.Locals("role").(string)

	log.Printf("👤 User %s (ID: %d, role: %s) akses GET /api/pekerjaan/%s", username, userID, role, id)

	p, err := repo.GetByID(role, userID, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if p == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}

	filtered := map[string]interface{}{
		"id":                    p.ID.Hex(),
		"alumni_id":             p.AlumniID,
		"nama_perusahaan":       p.NamaPerusahaan,
		"posisi_jabatan":        p.PosisiJabatan,
		"bidang_industri":       p.BidangIndustri,
		"lokasi_kerja":          p.LokasiKerja,
		"gaji_range":            p.GajiRange,
		"tanggal_mulai_kerja":   p.TanggalMulaiKerja,
		"tanggal_selesai_kerja": p.TanggalSelesaiKerja,
		"status_pekerjaan":      p.StatusPekerjaan,
		"deskripsi_pekerjaan":   p.DeskripsiPekerjaan,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    filtered,
	})
}

// GetAllAlumniWithPekerjaan godoc
// @Summary Menampilkan semua alumni beserta pekerjaan
// @Description Mengambil daftar alumni dan pekerjaan mereka (khusus admin). Dapat difilter dan diurutkan menggunakan parameter query.
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param search query string false "Cari alumni berdasarkan nama atau perusahaan"
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah data per halaman (default: 10)"
// @Param sort_by query string false "Kolom untuk sorting (contoh: nama, perusahaan)"
// @Param order query string false "Urutan data (asc/desc)"
// @Success 200 {array} models.AlumniPekerjaan
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pekerjaan/alumni-pekerjaan [get]
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

// GetPekerjaanByAlumniID godoc
// @Summary Menampilkan pekerjaan berdasarkan Alumni ID
// @Description Mengambil semua pekerjaan berdasarkan ID alumni (khusus admin). Dapat difilter dan diurutkan dengan query parameter.
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param alumni_id path int true "ID Alumni"
// @Param search query string false "Cari berdasarkan nama perusahaan atau posisi jabatan"
// @Param sort_by query string false "Urutkan berdasarkan kolom (contoh: nama_perusahaan, tanggal_mulai_kerja)"
// @Param order query string false "Urutan data (asc/desc)"
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah data per halaman (default: 10)"
// @Success 200 {array} models.Pekerjaan
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pekerjaan/alumni/{alumni_id} [get]
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


// CreatePekerjaanService godoc
// @Summary Tambah data pekerjaan baru
// @Description Membuat data pekerjaan baru di MongoDB
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param validate query bool false "Validasi data sebelum insert (true/false)"
// @Param body body models.CreatePekerjaanRequest true "Data pekerjaan baru"
// @Success 201 {object} models.ResponsePekerjaan
// @Router /api/pekerjaan [post]
func CreatePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New()
	userID := c.Locals("user_id").(int)
	var req models.CreatePekerjaanRequest

	// Parsing request body JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// Parsing tanggal mulai
	tMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tanggal mulai tidak valid"})
	}

	// Parsing tanggal selesai (opsional)
	var tSelesai *time.Time
	if req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Tanggal selesai tidak valid"})
		}
		tSelesai = &t
	}

	// Membuat object pekerjaan
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
		CreatedBy:           userID,
	}

	// Insert ke repository
	if err := repo.CreatePekerjaan(p); err != nil {
		log.Printf("❌ Create error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat pekerjaan"})
	}

	// Mapping ke ResponsePekerjaan
	res := models.ResponsePekerjaan{
		ID:                  p.ID.Hex(),
		AlumniID:            p.AlumniID,
		NamaPerusahaan:      p.NamaPerusahaan,
		PosisiJabatan:       p.PosisiJabatan,
		BidangIndustri:      p.BidangIndustri,
		LokasiKerja:         p.LokasiKerja,
		GajiRange:           p.GajiRange,
		TanggalMulaiKerja:   p.TanggalMulaiKerja,
		TanggalSelesaiKerja: p.TanggalSelesaiKerja,
		StatusPekerjaan:     p.StatusPekerjaan,
		DeskripsiPekerjaan:  p.DeskripsiPekerjaan,
	}

	// Response API
	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}


// UpdatePekerjaanService godoc
// @Summary Update data pekerjaan
// @Description Mengubah data pekerjaan berdasarkan ID
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Param notify query bool false "Kirim notifikasi ke alumni (true/false)"
// @Param body body models.UpdatePekerjaanRequest true "Data pekerjaan yang akan diupdate"
// @Success 200 {object} map[string]interface{}
// @Router /api/pekerjaan/{id} [put]
func UpdatePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New()
	id := c.Params("id")
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var req models.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
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
	}

	if err := repo.UpdatePekerjaan(id, p, role, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil diupdate",
	})
}

// DeletePekerjaanService godoc
// @Summary Hapus data pekerjaan (soft delete)
// @Description Menghapus pekerjaan dengan menandainya sebagai terhapus. Bisa diatur mode penghapusan dengan parameter query.
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Param permanent query bool false "Hapus permanen (true) atau soft delete (false)"
// @Param reason query string false "Alasan penghapusan (opsional)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pekerjaan/{id} [delete]
func DeletePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New()
	id := c.Params("id")
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	err := repo.SoftDeletePekerjaan(id, userID, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Riwayat pekerjaan berhasil dihapus (soft delete)",
	})
}

// RestorePekerjaanService godoc
// @Summary Restore data pekerjaan
// @Description Mengembalikan data pekerjaan yang sebelumnya dihapus (soft delete)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Param notify query bool false "Kirim notifikasi ke user terkait (true/false)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pekerjaan/{id}/restore [post]
func RestorePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New()
	id := c.Params("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	if err := repo.RestorePekerjaan(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data berhasil direstore"})
}

// HardDeletePekerjaanService godoc
// @Summary Hapus permanen data pekerjaan
// @Description Menghapus data pekerjaan secara permanen dari database. Gunakan dengan hati-hati!
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Param confirm query bool true "Konfirmasi hapus permanen (true untuk melanjutkan)"
// @Param admin_reason query string false "Alasan admin menghapus data ini (opsional)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pekerjaan/{id}/hard [delete]
func HardDeletePekerjaanService(c *fiber.Ctx) error {
	repo := repository.New()
	id := c.Params("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	if err := repo.HardDeletePekerjaanByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Data dengan ID %s dihapus permanen", id),
	})
}

// GetTrashPekerjaanService godoc
// @Summary Menampilkan daftar pekerjaan yang dihapus (trash)
// @Description Menampilkan semua data pekerjaan yang sudah dihapus secara soft delete. Dapat difilter dengan parameter query.
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah data per halaman (default: 10)"
// @Param search query string false "Kata kunci pencarian (nama perusahaan / posisi)"
// @Param order query string false "Urutan data (asc/desc)"
// @Success 200 {array} models.GetTrashPekerjaan
// @Failure 500 {object} map[string]string
// @Router /api/pekerjaan/trash [get]
func GetTrashPekerjaanService(c *fiber.Ctx) error {
	repo := repository.New()
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	allTrash, err := repo.GetTrashPekerjaan(userID, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil trash"})
	}

	// Konversi ObjectID ke string (hex)
	var result []map[string]interface{}
	for _, t := range allTrash {
		result = append(result, map[string]interface{}{
			"id":               t.ID.Hex(),
			"alumni_id":        t.AlumniID,
			"nama_perusahaan":  t.NamaPerusahaan,
			"posisi_jabatan":   t.PosisiJabatan,
			"bidang_industri":  t.BidangIndustri,
			"lokasi_kerja":     t.LokasiKerja,
			"status_pekerjaan": t.StatusPekerjaan,
			"deleted_at":       t.DeletedAt,
			"created_by":       t.CreatedBy,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(result),
		"data":    result,
	})
}



// GetPekerjaanPaginated godoc
// @Summary Menampilkan data pekerjaan dengan pagination
// @Description Mengambil data pekerjaan dengan pagination, sorting, dan search
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah data per halaman (default: 5)"
// @Param sort_by query string false "Kolom untuk sorting (default: created_at)"
// @Param order query string false "Urutan sort asc/desc (default: desc)"
// @Param search query string false "Kata kunci pencarian"
// @Param only_active query bool false "Tampilkan hanya pekerjaan aktif"
// @Success 200 {object} models.PekerjaanResponse
// @Router /api/pekerjaan/paginated [get]
func GetPekerjaanPaginated(c *fiber.Ctx) error {
	repo := repository.New()
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

	// whitelist kolom yang boleh di-sort
	sortByWhitelist := map[string]bool{
		"_id":                 true,
		"nama_perusahaan":     true,
		"posisi_jabatan":      true,
		"tanggal_mulai_kerja": true,
		"created_at":          true,
	}
	if _, ok := sortByWhitelist[sortBy]; !ok {
		sortBy = "_id"
	}

	order = strings.ToLower(order)
	if order != "desc" {
		order = "asc"
	}

	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	// ambil data dari repository
	pekerjaanList, err := repo.GetPekerjaanPaginated(search, sortBy, order, limit, offset, role, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	total, err := repo.CountPekerjaan(search, role, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// langsung pakai ObjectID tanpa convert
	response := &models.PekerjaanResponse{
		Data: pekerjaanList, // []*models.Pekerjaan
		Meta: &models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  int(total),
			Pages:  (int(total) + limit - 1) / limit,
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