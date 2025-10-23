package service

import (
	"os"
	"path/filepath"
	"strconv"

	// "config"
	// "alumniproject/app/models/mongodb"
	// "alumniproject/app/repository/mongodb"
	// "alumniproject/config"
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

// -------------------- UPLOAD FOTO --------------------
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

// -------------------- GET ALL FOTO --------------------
func GetAllFoto(c *fiber.Ctx) error {
    fotoRepo := repository.NewFotoRepository(db.DB) // ✅ simpan ke variabel
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

// -------------------- GET FOTO BY ID --------------------
func GetFotoByID(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid ID format",
        })
    }

    fotoRepo := repository.NewFotoRepository(db.DB) // ✅ gunakan alias variabel
    foto, err := fotoRepo.FindByID(id)        // ✅ panggil dari variabel fotoRepo
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


// -------------------- DELETE FOTO --------------------
func DeleteFoto(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid ID format",
        })
    }

    fotoRepo := repository.NewFotoRepository(db.DB) // ✅ gunakan variabel
    foto, err := fotoRepo.FindByID(id)        // ✅ panggil lewat variabel
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Photo not found",
        })
    }

    os.Remove(foto.FilePath)
    fotoRepo.Delete(id) // ✅ gunakan variabel yang sama

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Photo deleted successfully",
    })
}
