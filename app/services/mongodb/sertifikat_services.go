package service

import (
    "os"
    "path/filepath"
    "strconv"

    db "alumniproject/database/mongodb"          // 🔹 alias db
    models "alumniproject/app/models/mongodb"    // 🔹 alias models
    repo"alumniproject/app/repository/mongodb"  // 🔹 alias repo
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

// -------------------- UPLOAD --------------------
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

// -------------------- GET ALL --------------------
func GetAllSertifikat(c *fiber.Ctx) error {
    repo := repo.NewFileRepository(db.DB)
    files, err := repo.FindAll()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to get certificates",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Certificates retrieved successfully",
        "data":    files,
    })
}

// -------------------- GET BY ID --------------------
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

// -------------------- DELETE --------------------
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
