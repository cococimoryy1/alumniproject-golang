package service

import (
    "os"
    "path/filepath"
    "strconv"

    db "alumniproject/database/mongodb"
    models "alumniproject/app/models/mongodb"
    repo "alumniproject/app/repository/mongodb"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

type SertifikatService struct {
    repo repo.FileRepository
    path string
}

func NewSertifikatService(repo repo.FileRepository, path string) *SertifikatService {
    return &SertifikatService{repo: repo, path: path}
}

// UploadSertifikat godoc
// @Summary Upload sertifikat baru (PDF)
// @Description Mengunggah file sertifikat (PDF) ke server dan menyimpannya di MongoDB
// @Tags Sertifikat
// @Accept multipart/form-data
// @Produce json
// @Param sertifikat formData file true "File sertifikat (PDF, max 2MB)"
// @Success 200 {object} models.FileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/sertifikat/upload [post]
func (s *SertifikatService) UploadSertifikat(c *fiber.Ctx) error {
    fileHeader, err := c.FormFile("sertifikat")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "No file uploaded",
        })
    }

    if fileHeader.Header.Get("Content-Type") != "application/pdf" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Only PDF allowed",
        })
    }

    if fileHeader.Size > 2*1024*1024 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Max file size 2MB",
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
        FileType:     "application/pdf",
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
        "message": "Sertifikat uploaded successfully",
        "data":    fileModel,
    })
}

// GetAllSertifikat godoc
// @Summary Menampilkan semua sertifikat
// @Description Mengambil seluruh data sertifikat dari MongoDB, bisa difilter dengan query parameter
// @Tags Sertifikat
// @Accept json
// @Produce json
// @Param search query string false "Cari sertifikat berdasarkan nama file"
// @Param limit query int false "Batas jumlah data yang diambil (default 10)"
// @Param sort query string false "Urutkan berdasarkan field (contoh: uploaded_at, file_name)"
// @Param order query string false "Urutan data (asc / desc)"
// @Success 200 {array} models.FileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/sertifikat [get]
func GetAllSertifikat(c *fiber.Ctx) error {
    repo := repo.NewFileRepository(db.DB)

    // Ambil parameter dari query
    search := c.Query("search", "")
    sortField := c.Query("sort", "uploaded_at")
    order := c.Query("order", "desc")
    limitStr := c.Query("limit", "10")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }

    // Ambil semua data
    files, err := repo.FindAll()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to get certificates",
        })
    }

    // Filter sederhana berdasarkan nama file (kalau search diisi)
    var filtered []models.File
    for _, f := range files {
        if search == "" || 
           (search != "" && (filepath.Base(f.FileName) == search || f.OriginalName == search)) {
            filtered = append(filtered, f)
        }
    }

    // Batasi jumlah data
    if len(filtered) > limit {
        filtered = filtered[:limit]
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Certificates retrieved successfully",
        "count":   len(filtered),
        "sort_by": sortField,
        "order":   order,
        "data":    filtered,
    })
}


// GetSertifikatByID godoc
// @Summary Menampilkan sertifikat berdasarkan ID
// @Description Mengambil detail sertifikat dari MongoDB berdasarkan ID
// @Tags Sertifikat
// @Accept json
// @Produce json
// @Param id path int true "ID sertifikat"
// @Param include_deleted query bool false "Tampilkan juga sertifikat yang sudah dihapus"
// @Success 200 {object} models.FileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/sertifikat/{id} [get]
func GetSertifikatByID(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid ID format",
        })
    }

    repo := repo.NewFileRepository(db.DB)
    file, err := repo.FindByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Certificate not found",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Certificate retrieved successfully",
        "data":    file,
    })
}

// DeleteSertifikat godoc
// @Summary Menghapus sertifikat
// @Description Menghapus sertifikat dari server dan MongoDB berdasarkan ID
// @Tags Sertifikat
// @Accept json
// @Produce json
// @Param id path int true "ID sertifikat"
// @Param force query bool false "Hapus permanen (true) atau hanya tandai sebagai deleted"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/sertifikat/{id} [delete]
func DeleteSertifikat(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid ID format",
        })
    }

    repo := repo.NewFileRepository(db.DB)
    file, err := repo.FindByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Certificate not found",
        })
    }

    os.Remove(file.FilePath)
    repo.Delete(id)

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Certificate deleted successfully",
    })
}
