package routes

import (
    "github.com/gofiber/fiber/v2"
    service "alumniproject/app/services/mongodb"
    middleware "alumniproject/middleware/mongodb"
	repo "alumniproject/app/repository/mongodb"
    db "alumniproject/database/mongodb"
)

func SetupMongoRoutes(app *fiber.App) {
    api := app.Group("/api")

    // =============================
    // AUTH & LOGIN
    // =============================
    api.Post("/login", service.Login)

    // =============================
    // PEKERJAAN ROUTES
    // =============================
    pekerjaan := api.Group("/pekerjaan")

    pekerjaan.Get("/", middleware.AuthRequired(), service.GetAllPekerjaanService)
    pekerjaan.Get("/:id", middleware.AuthRequired(), service.GetPekerjaanByID)
    pekerjaan.Get("/alumni/:alumni_id", middleware.AuthRequired(), middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
    pekerjaan.Get("/paginated", middleware.AuthRequired(), service.GetPekerjaanPaginated)
    pekerjaan.Get("/alumni-pekerjaan", middleware.AuthRequired(), middleware.AdminOnly(), service.GetAllAlumniWithPekerjaan)

    pekerjaan.Post("/", middleware.AuthRequired(), service.CreatePekerjaanService)
    pekerjaan.Put("/:id", middleware.AuthRequired(), service.UpdatePekerjaanService)
    pekerjaan.Delete("/:id", middleware.AuthRequired(), middleware.AdminOrOwner(), service.DeletePekerjaanService)
    pekerjaan.Get("/trash", middleware.AuthRequired(), service.GetTrashPekerjaanService)
    pekerjaan.Post("/:id/restore", middleware.AuthRequired(), service.RestorePekerjaanService)
    pekerjaan.Delete("/:id/hard", middleware.AuthRequired(), middleware.AdminOnly(), service.HardDeletePekerjaanService)

    // UPLOAD FOTO & SERTIFIKAT
    // =============================
    fotoRepo := repo.NewFotoRepository(db.DB)
    sertifikatRepo := repo.NewFileRepository(db.DB)

    fotoService := service.NewFotoService(fotoRepo, "/uploads/foto")
    sertifikatService := service.NewSertifikatService(sertifikatRepo, "./uploads/sertifikat")

    foto := api.Group("/foto", middleware.AuthRequired())
    sertifikat := api.Group("/sertifikat", middleware.AuthRequired())

    // Foto routes
    foto.Post("/upload", middleware.AuthRequired(), fotoService.UploadFoto)
    foto.Get("/", middleware.AuthRequired(), service.GetAllFoto)
    foto.Get("/:id", middleware.AuthRequired(), service.GetFotoByID)
    foto.Delete("/:id", middleware.AuthRequired(), middleware.AdminOrOwner(), service.DeleteFoto)

    // Sertifikat routes
    sertifikat.Post("/upload", middleware.AuthRequired(), sertifikatService.UploadSertifikat)
    sertifikat.Get("/", middleware.AuthRequired(), service.GetAllSertifikat)
    sertifikat.Get("/:id", middleware.AuthRequired(), service.GetSertifikatByID)
    sertifikat.Delete("/:id", middleware.AuthRequired(), middleware.AdminOrOwner(), service.DeleteSertifikat)
}