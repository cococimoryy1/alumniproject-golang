package service

import (
	"alumniproject/app/models/mongodb"
	"alumniproject/app/repository/mongodb"
	mongodbutils "alumniproject/utils/mongodb"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
    fmt.Println("✅ Route /login terpanggil")

    var req models.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        fmt.Println("❌ BodyParser error:", err)
        return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
    }

    fmt.Println("🔍 Input dari user => Username:", req.Username, "Password:", req.Password)

    if req.Username == "" || req.Password == "" {
        fmt.Println("⚠️ Username atau password kosong")
        return c.Status(400).JSON(fiber.Map{"error": "Username dan password harus diisi"})
    }

    user, passwordHash, err := repository.GetUserByUsernameOrEmail(req.Username)
    if err != nil {
        fmt.Println("❌ User tidak ditemukan:", err)
        return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
    }

    fmt.Println("✅ User ditemukan di DB:", user.Username)
    fmt.Println("🔐 Password Hash dari DB:", passwordHash)

    if !mongodbutils.CheckPassword(req.Password, passwordHash) {
        fmt.Println("❌ Password tidak cocok")
        return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
    }

    token, err := mongodbutils.GenerateToken(user)
    if err != nil {
        fmt.Println("❌ Gagal generate token:", err)
        return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
    }

    fmt.Println("✅ Login berhasil untuk user:", user.Username)

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Login berhasil",
        "data": fiber.Map{
            "user":  user,
            "token": token,
        },
    })
}
