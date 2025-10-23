package service

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"alumniproject/app/models/postgresql"
	"alumniproject/app/repository/postgresql"
)

// AlumniPekerjaanResponse wraps the response for alumni by status pekerjaan
type AlumniPekerjaanResponse struct {
	Data       []models.AlumniPekerjaan `json:"data"`
	TotalCount int                      `json:"total_count"`
}

// GetAlumniByStatusPekerjaan retrieves alumni by job status
func GetAlumniByStatusPekerjaan(c *fiber.Ctx) error {
	status := c.Params("status")
	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/alumni-pekerjaan/status dengan status: %s", username, status)

	response, err := repository.GetAlumniByStatusPekerjaan(status)
	if err != nil {
		log.Printf("Error mengambil data alumni: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengambil data alumni berdasarkan status pekerjaan: " + err.Error(),
		})
	}

	if len(response) == 0 {
		return c.JSON(fiber.Map{
			"success": false,
			"count":   0,
			"data":    []models.AlumniPekerjaan{},
			"message": "Tidak ada data alumni dengan status pekerjaan " + status,
		})
	}

	totalCount := response[0].TotalBekerjaLebih1Tahun

	return c.JSON(fiber.Map{
		"success": true,
		"count":   totalCount,
		"data":    response,
		"message": "Data alumni dengan status pekerjaan " + status + " berhasil diambil",
	})
}

// GetAlumniWithLongTermJobs retrieves alumni with active jobs lasting more than 1 year
func GetAlumniWithLongTermJobs(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/alumni-pekerjaan/long-term", username)

	response, err := repository.GetAlumniWithLongTermJobs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengambil data alumni dengan pekerjaan lebih dari 1 tahun",
		})
	}

	if len(response) == 0 {
		return c.JSON(fiber.Map{
			"success": false,
			"count":   0,
			"data":    []models.AlumniPekerjaan{},
			"message": "Tidak ada data alumni dengan pekerjaan aktif lebih dari 1 tahun",
		})
	}

	totalCount := response[0].TotalBekerjaLebih1Tahun

	return c.JSON(fiber.Map{
		"success": true,
		"count":   totalCount,
		"data":    response,
		"message": "Data alumni dengan pekerjaan aktif lebih dari 1 tahun berhasil diambil",
	})
}