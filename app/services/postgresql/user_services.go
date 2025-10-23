package service

import (
	// "errors"

	"github.com/gofiber/fiber/v2"
	"alumniproject/app/models/postgresql"
	"alumniproject/app/repository/postgresql"
	"alumniproject/utils/postgresql"
)

// Login handles user authentication and token generation
func Login(c *fiber.Ctx) error {
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

	if !postgresutils.CheckPassword(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	token, err := postgresutils.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}
