package postgresql

import (
	service "alumniproject/app/services/postgresql"
	"alumniproject/middleware/postgresql"
	"github.com/gofiber/fiber/v2"
)

// SetupPostgresRoutes mengatur semua endpoint API untuk PostgreSQL
func SetupPostgresRoutes(app *fiber.App) {
	api := app.Group("/api")

	// --- Public route ---
	api.Post("/login", service.Login)

	// --- Protected routes ---
	protected := api.Group("", middleware.AuthRequired())

	// === ALUMNI ROUTES ===
	alumni := protected.Group("/alumni")
	alumni.Get("/", service.GetAllAlumni)
	alumni.Get("/all", service.GetAlumniService)
	alumni.Get("/:id", service.GetAlumniByIDService)
	alumni.Post("/", service.CreateAlumniService)
	alumni.Put("/:id", service.UpdateAlumniService)
	alumni.Delete("/:id", service.DeleteAlumniService)
	alumni.Put("/restore/:id", service.RestoreAlumniService)

	// === PEKERJAAN ROUTES ===
	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/trash", service.GetTrashPekerjaanService)
	pekerjaan.Get("/", service.GetAllPekerjaanService)
	pekerjaan.Get("/:id", service.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
	pekerjaan.Post("/", service.CreatePekerjaanService)
	pekerjaan.Delete("/:id", service.DeletePekerjaanService)
	pekerjaan.Delete("/hard-delete/:id", service.HardDeletePekerjaanService)
	pekerjaan.Put("/restore/:id", service.RestorePekerjaanService)
	pekerjaan.Put("/:id", service.UpdatePekerjaanService)

	// === ALUMNI + PEKERJAAN COMBINED ===
	alumniPekerjaan := protected.Group("/alumni-pekerjaan")
	alumniPekerjaan.Get("/", service.GetAllAlumniWithPekerjaan)
	alumniPekerjaan.Get("/long-term", service.GetAlumniWithLongTermJobs)
	alumniPekerjaan.Get("/status/:status", service.GetAlumniByStatusPekerjaan)
}




// func GetProfile(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(int)
// 	username := c.Locals("username").(string)
// 	role := c.Locals("role").(string)
// 	return c.JSON(fiber.Map{
// 		"success": true,
// 		"message": "Profile berhasil diambil",
// 		"data": fiber.Map{
// 			"user_id":  userID,
// 			"username": username,
// 			"role":     role,
// 		},
// 	})
// }


// func countTotalJobs(alumniList []models.AlumniWithPekerjaan) int {
// 	total := 0
// 	for _, item := range alumniList {
// 		total += len(item.Pekerjaan)
// 	}
// 	return total
// }

