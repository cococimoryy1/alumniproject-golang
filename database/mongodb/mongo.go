package database

import (
    "context"
    "log"
    "os"
    "time"

	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    DB                  *mongo.Database
    MongoClient         *mongo.Client
    PekerjaanCollection *mongo.Collection
	UsersCollection     *mongo.Collection // ✅ Tambahan untuk koleksi user
	CountersCollection  *mongo.Collection // ✅ Tambahan untuk koleksi counters (auto-increment ID)
    
)

func ConnectMongo() {
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }

    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal("Mongo connection failed:", err)
    }

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Mongo ping failed:", err)
    }

    MongoClient = client
    DB = client.Database("alumni") // Sesuaikan nama database
    PekerjaanCollection = DB.Collection("pekerjaan")
    UsersCollection = DB.Collection("users")
    CountersCollection = DB.Collection("counters")

    // ✅ Fix: Set counter ke 10 explicit (tanpa inc, hindari conflict)
	ctxInit, cancelInit := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelInit()
    filter := bson.M{"_id": "pekerjaan_id"}
    update := bson.M{"$set": bson.M{"seq": 10}}
    opts := options.Update().SetUpsert(true)
    _, err = CountersCollection.UpdateOne(ctxInit, filter, update, opts)
    if err != nil {
        log.Printf("⚠️ Warning: Gagal set initial counter: %v", err)
    } else {
        log.Printf("✅ Counter pekerjaan di-set ke 10 (next ID: 11)")
    }

    log.Println("✅ MongoDB Connected - All Collections Ready!")
}