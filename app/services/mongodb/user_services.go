package service

import (
	"alumniproject/app/models/mongodb"
	"alumniproject/app/repository/mongodb"
	mongodbutils "alumniproject/utils/mongodb"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Login godoc
// @Summary Login user
// @Description Melakukan autentikasi user berdasarkan username/email dan password, lalu mengembalikan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param lang query string false "Bahasa respon (contoh: id atau en)"
// @Param remember query bool false "Login dengan mode remember (persistent)"
// @Param body body models.LoginRequest true "Data login user"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string "Request body tidak valid atau kosong"
// @Failure 401 {object} map[string]string "Username atau password salah"
// @Failure 500 {object} map[string]string "Gagal generate token"
// @Router /api/login [post]
func Login(c *fiber.Ctx) error {
	fmt.Println("✅ Route /login terpanggil")

	// ambil query parameter opsional
	lang := c.Query("lang", "id")
	remember := c.Query("remember", "false")
	fmt.Printf("🌐 Language: %s | Remember: %s\n", lang, remember)

	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username dan password harus diisi"})
	}

	user, passwordHash, err := repository.GetUserByUsernameOrEmail(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	if !mongodbutils.CheckPassword(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	token, err := mongodbutils.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	// bisa pakai parameter remember untuk memperpanjang expiry token (kalau kamu implementasikan)
	// bisa pakai parameter lang untuk ubah pesan respon (id/en)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"lang":    lang,
		"remember": remember,
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}
