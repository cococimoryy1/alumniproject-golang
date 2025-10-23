package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"alumniproject/utils/mongodb"
)

// AuthRequired memverifikasi JWT dan menyimpan data user di Locals
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token akses diperlukan"})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Format token tidak valid"})
		}

		claims, err := mongodbutils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid atau sudah expired"})
		}

		userIDStr := claims.UserID
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "User ID tidak valid"})
		}

		c.Locals("user_id", userID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// AdminOnly membatasi akses hanya untuk admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan."})
		}
		return c.Next()
	}
}

// AdminOrOwner membatasi akses hanya untuk admin (check owner dihandle di service untuk akurasi)
func AdminOrOwner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role == "admin" {
			return c.Next()
		}
		// Untuk non-admin, biarkan service check owner
		return c.Next()
	}
}