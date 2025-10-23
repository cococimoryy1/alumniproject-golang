package main

import (
    "log"
    "os"

    "alumniproject/config"
    "alumniproject/database/mongodb"
    "alumniproject/database/postgresql"
    mongoRoutes "alumniproject/routes/mongodb"
    pgRoutes "alumniproject/routes/postgresql"

    "github.com/gofiber/fiber/v2"
)

func main() {
    // Load environment variables
    config.LoadEnv()
    log.Println("📦 DB_TYPE from os.Getenv:", os.Getenv("DB_TYPE"))

    // Set up logger
    config.SetupLogger()

    // Setup Fiber app
    app := config.SetupApp()

    // Tentukan DB_TYPE dari .env
    dbType := os.Getenv("DB_TYPE")
    if dbType == "" {
        dbType = "postgres"
    }

    // Pilih database berdasarkan DB_TYPE
    switch dbType {
    case "mongodb":
        database.ConnectMongo()
        mongoRoutes.SetupMongoRoutes(app)
        log.Println("✅ MongoDB Connected and Routes Registered")
    case "postgres":
        postgresql.ConnectPostgres()
        pgRoutes.SetupPostgresRoutes(app)
        log.Println("✅ PostgreSQL Connected and Routes Registered")
    default:
        log.Fatalf("❌ Unknown DB_TYPE: %s", dbType)
    }

    // Root endpoint
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Selamat datang di API Alumni dan Pekerjaan",
            "db_used": dbType,
        })
    })

    // Jalankan server
    if err := app.Listen(":3000"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
// package main

// import (
// 	"log"

// 	"alumniproject/config"
// 	"alumniproject/database"
// 	"alumniproject/routes"
// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	// Load environment variables
// 	config.LoadEnv()

// 	// Setup logger
// 	config.SetupLogger()

// 	// Connect to MongoDB
// 	database.ConnectDB()
// 	// Tidak perlu defer database.DB.Close() karena driver Mongo handle sendiri

// 	// Setup Fiber app
// 	app := config.SetupApp()

// 	// Root endpoint
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{"message": "Selamat datang di API Alumni dan Pekerjaan (MongoDB)"})
// 	})

// 	// Register routes
// 	routes.SetupRoutes(app)

// 	// Start server
// 	if err := app.Listen(":3000"); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }