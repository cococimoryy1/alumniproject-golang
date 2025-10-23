package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectPostgres() {
	var err error
	dsn := "host=localhost user=postgres password=2255 dbname=alumnidb port=5432 sslmode=disable"
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}
	fmt.Println("Berhasil terhubung ke database PostgreSQL")
}

// package database

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var DB *mongo.Database

// func ConnectDB() {
// 	// Ambil URI MongoDB dari .env (kalau ada)
// 	mongoURI := os.Getenv("MONGO_URI")
// 	if mongoURI == "" {
// 		mongoURI = "mongodb://localhost:27017"
// 	}

// 	dbName := os.Getenv("DATABASE_NAME")
// 	if dbName == "" {
// 		dbName = "alumnidb" // default sama seperti yang kamu pakai di Compass
// 	}

// 	// Setup koneksi
// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal konek ke MongoDB: %v", err)
// 	}

// 	// Tes koneksi
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal ping ke MongoDB: %v", err)
// 	}

// 	fmt.Println("✅ Berhasil terhubung ke MongoDB!")
// 	DB = client.Database(dbName)
// }
