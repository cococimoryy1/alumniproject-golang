package config

import "github.com/gofiber/fiber/v2"

func SetupApp() *fiber.App {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        },
    })
    return app
}