package service

import (
	"os"
	"path/filepath"
	"strconv"

	models "alumniproject/app/models/mongodb"
	repository "alumniproject/app/repository/mongodb"
	db "alumniproject/database/mongodb"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FotoService struct {
	repo repository.FotoRepository
	path string
}

func NewFotoService(repo repository.FotoRepository, path string) *FotoService {
	return &FotoService{repo: repo, path: path}
}

// UploadFoto godoc
// @Summary Upload foto baru
// @Description Mengunggah file foto ke server dan menyimpannya di MongoDB
// @Tags Foto
// @Accept multipart/form-data
// @Produce json
// @Param foto formData file true "File foto (jpeg/jpg/png, max 1MB)"
// @Param kategori formData string false "Kategori foto (contoh: profil, dokumen, event)"
// @Param deskripsi formData string false "Deskripsi singkat foto"
// @Param uploader_id formData int false "ID pengguna yang mengunggah (opsional)"
// @Success 200 {object} models.FileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/foto/upload [post]
func (s *FotoService) UploadFoto(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("foto")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
		})
	}

	allowed := map[string]bool{"image/jpeg": true, "image/jpg": true, "image/png": true}
	if !allowed[fileHeader.Header.Get("Content-Type")] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Only jpeg/jpg/png allowed",
		})
	}
	if fileHeader.Size > 1*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Max file size 1MB",
		})
	}

	os.MkdirAll(s.path, os.ModePerm)
	newFileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	filePath := filepath.Join(s.path, newFileName)

	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file",
		})
	}

	fileModel := &models.File{
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     fileHeader.Header.Get("Content-Type"),
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save metadata",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Foto uploaded successfully",
		"data":    fileModel,
	})
}

// GetAllFoto godoc
// @Summary Menampilkan semua foto
// @Description Mengambil seluruh data foto yang tersimpan di MongoDB dengan opsi filter, search, pagination, dan sorting
// @Tags Foto
// @Accept json
// @Produce json
// @Param search query string false "Cari berdasarkan nama file atau nama asli foto"
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah data per halaman (default: 10)"
// @Param sort_by query string false "Kolom untuk sorting (default: uploaded_at)"
// @Param order query string false "Urutan sort (asc/desc)"
// @Param file_type query string false "Filter berdasarkan tipe file (contoh: image/jpeg, image/png)"
// @Success 200 {array} models.FileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/foto [get]
func GetAllFoto(c *fiber.Ctx) error {
	fotoRepo := repository.NewFotoRepository(db.DB)
	fotos, err := fotoRepo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get photos",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Photos retrieved successfully",
		"data":    fotos,
	})
}

// GetFotoByID godoc
// @Summary Menampilkan foto berdasarkan ID
// @Description Mengambil detail satu foto berdasarkan ID dan bisa menentukan format respon
// @Tags Foto
// @Accept json
// @Produce json
// @Param id path int true "ID foto"
// @Param include_metadata query bool false "Tampilkan metadata tambahan (true/false)"
// @Param view_mode query string false "Mode tampilan (contoh: thumbnail/full)"
// @Success 200 {object} models.FileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/foto/{id} [get]
func GetFotoByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid ID format",
		})
	}

	fotoRepo := repository.NewFotoRepository(db.DB)
	foto, err := fotoRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Photo not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Photo retrieved successfully",
		"data":    foto,
	})
}

// DeleteFoto godoc
// @Summary Menghapus foto
// @Description Menghapus file foto dari server dan database berdasarkan ID
// @Tags Foto
// @Accept json
// @Produce json
// @Param id path int true "ID foto"
// @Param permanent query bool false "Hapus permanen (true) atau soft delete (false)"
// @Param reason query string false "Alasan penghapusan foto (opsional)"
// @Param admin_id query int false "ID admin yang menghapus (opsional)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/foto/{id} [delete]
func DeleteFoto(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid ID format",
		})
	}

	fotoRepo := repository.NewFotoRepository(db.DB)
	foto, err := fotoRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Photo not found",
		})
	}

	os.Remove(foto.FilePath)
	fotoRepo.Delete(id)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Photo deleted successfully",
	})
}
