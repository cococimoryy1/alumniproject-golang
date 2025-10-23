package repository

import (
    "context"
    "fmt"
    "log"
    "strconv"
    "strings"
    "time"
  

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    
    "go.mongodb.org/mongo-driver/mongo/options"
    "alumniproject/app/models/mongodb" // Asumsi path benar
    "alumniproject/database/mongodb"

)

type PekerjaanMongoRepo struct{}

// GetAllPekerjaan: Sudah OK, tapi tambah log
func (r *PekerjaanMongoRepo) GetAllPekerjaan(role string, userID int) ([]*models.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"deleted_at": nil}
    if role != "admin" {
        filter["created_by"] = userID
    }

    log.Printf("🔍 GetAllPekerjaan filter: %+v (role: %s, userID: %d)", filter, role, userID)

    opts := options.Find().SetSort(bson.D{{"created_at", -1}})
    cursor, err := database.PekerjaanCollection.Find(ctx, filter, opts)
    if err != nil {
        log.Printf("❌ Find error: %v", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []*models.Pekerjaan
    if err = cursor.All(ctx, &results); err != nil {
        log.Printf("❌ Decode error: %v", err)
        return nil, err
    }

    log.Printf("📊 Found %d pekerjaan records", len(results))
    return results, nil
}

// GetPekerjaanByID: Tambah role/userID untuk filter owner
func (r *PekerjaanMongoRepo) GetPekerjaanByID(id string, role string, userID int) (*models.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    idInt, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid ID format: %v", err)
    }

    filter := bson.M{"_id": idInt, "deleted_at": nil} // ✅ _id
    if role != "admin" {
        filter["created_by"] = userID
    }

    log.Printf("🔍 GetByID filter: %+v (role: %s, userID: %d)", filter, role, userID)

    var p models.Pekerjaan
    err = database.PekerjaanCollection.FindOne(ctx, filter).Decode(&p)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("mongo: no documents in result")
        }
        log.Printf("❌ Error decoding pekerjaan with ID %s: %v", id, err)
        return nil, fmt.Errorf("failed to fetch pekerjaan: %v", err)
    }

    log.Printf("✅ Data ditemukan: %+v", p)
    return &p, nil
}

// GetPekerjaanByAlumniID: Tambah filter role/userID (asumsi admin only, tapi selaraskan)
func (r *PekerjaanMongoRepo) GetPekerjaanByAlumniID(alumniID int, role string, userID int) ([]*models.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"alumni_id": alumniID, "deleted_at": nil}
    // Untuk non-admin, filter tambah created_by? Asumsi admin only, tapi tambah jika perlu
    if role != "admin" {
        filter["created_by"] = userID
    }

    opts := options.Find().SetSort(bson.D{{"created_at", -1}})
    cursor, err := database.PekerjaanCollection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []*models.Pekerjaan
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    return results, nil
}

func (r *PekerjaanMongoRepo) GetNextSequence(sequenceName string) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id": sequenceName}
    update := bson.M{"$inc": bson.M{"seq": 1}}
    opts := options.FindOneAndUpdate().
        SetUpsert(true).
        SetReturnDocument(options.After)

    var updatedDoc bson.M
    err := database.CountersCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDoc)
    if err != nil {
        log.Printf("❌ Error GetNextSequence %s: %v", sequenceName, err)
        return 0, fmt.Errorf("gagal increment sequence: %v", err)
    }

    // Konversi aman ke int64
    var seq int64
    switch v := updatedDoc["seq"].(type) {
    case int32:
        seq = int64(v)
    case int64:
        seq = v
    case float64:
        seq = int64(v)
    default:
        return 0, fmt.Errorf("sequence value invalid")
    }

    log.Printf("✅ Next sequence for %s: %d", sequenceName, seq)
    return seq, nil
}

// CreatePekerjaan: Pakai sequential ID dari GetNextSequence (sudah include init)
func (r *PekerjaanMongoRepo) CreatePekerjaan(p *models.Pekerjaan) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // ✅ Get next sequential ID (auto-init counter ke 10 jika belum)
    nextID, err := r.GetNextSequence("pekerjaan_id")
    if err != nil {
        return err
    }
    p.ID = nextID // Set ID sequential (mulai 11)

    // Set timestamps
    now := time.Now()
    p.CreatedAt = now
    p.UpdatedAt = now
    p.DeletedAt = nil // Explicit null

    // Insert dengan ID custom
    _, err = database.PekerjaanCollection.InsertOne(ctx, p)
    if err != nil {
        log.Printf("❌ Insert error: %v", err)
        return fmt.Errorf("gagal membuat pekerjaan: %v", err)
    }

    log.Printf("✅ Insert berhasil, ID sequential: %d", nextID)
    return nil
}

// UpdatePekerjaan: Tambah role/userID untuk filter owner
func (r *PekerjaanMongoRepo) UpdatePekerjaan(id string, p *models.Pekerjaan, role string, userID int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    p.UpdatedAt = time.Now()

    idInt, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return fmt.Errorf("invalid ID format: %v", err)
    }

    filter := bson.M{"_id": idInt, "deleted_at": nil} // ✅ _id
    if role != "admin" {
        filter["created_by"] = userID
    }

    // Update hanya fields dari p, tapi exclude ID, timestamps (sudah set manual)
    update := bson.M{"$set": bson.M{
        "nama_perusahaan":      p.NamaPerusahaan,
        "posisi_jabatan":       p.PosisiJabatan,
        // ... tambah fields lain dari p
        "updated_at": p.UpdatedAt,
    }}

    _, err = database.PekerjaanCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Printf("❌ Update error: %v", err)
    }
    return err
}

// RestorePekerjaan: Fix bson.M dengan $exists untuk hindari nil type issue
func (r *PekerjaanMongoRepo) RestorePekerjaan(id int64) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // ✅ Fix: Gunakan $exists: false untuk match deleted_at null (hindari $ne: nil type conflict)
    filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": true}} // hanya yang ada deleted_at
    update := bson.M{"$unset": bson.M{"deleted_at": ""}}               // hapus field


    result, err := database.PekerjaanCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("gagal memperbarui data: %v", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("data tidak ditemukan")
    }

    log.Printf("✅ Restore berhasil untuk ID %d", id)
    return nil
}

// SoftDeletePekerjaan: Fix sama, gunakan $exists: false untuk deleted_at: nil
func (r *PekerjaanMongoRepo) SoftDeletePekerjaan(id string, userID int, role string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    idInt, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return fmt.Errorf("id tidak valid: %v", err)
    }

    filter := bson.M{"_id": idInt, "deleted_at": bson.M{"$exists": false}} // Tidak ada deleted_at (active)
    if role != "admin" {
        filter["created_by"] = userID
    }

    update := bson.M{"$set": bson.M{"deleted_at": time.Now()}} // ✅ benar

    result, err := database.PekerjaanCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("gagal soft delete: %v", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("data tidak ditemukan atau sudah dihapus")
    }

    log.Printf("✅ Soft delete berhasil untuk ID %d", idInt)
    return nil
}

// HardDeletePekerjaanByID: OK, _id
func (r *PekerjaanMongoRepo) HardDeletePekerjaanByID(id int64) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id": id, "deleted_at": bson.M{"$ne": nil}} // ✅ _id
    result, err := database.PekerjaanCollection.DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("gagal hard delete: %v", err)
    }
    if result.DeletedCount == 0 {
        return fmt.Errorf("data tidak ditemukan")
    }

    log.Printf("✅ Hard delete berhasil untuk ID %d", id)
    return nil
}

// GetTrashPekerjaan: OK
func (r *PekerjaanMongoRepo) GetTrashPekerjaan(userID int, role string) ([]*models.GetTrashPekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
    if role != "admin" {
        filter["created_by"] = userID
    }

    log.Printf("🔍 GetTrash filter: %+v", filter)

    opts := options.Find().SetSort(bson.D{{"deleted_at", -1}})
    cursor, err :=database.PekerjaanCollection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []*models.GetTrashPekerjaan
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    return results, nil
}

// GetPekerjaanPaginated: Tambah role/userID
func (r *PekerjaanMongoRepo) GetPekerjaanPaginated(search, sortBy, order string, limit, offset int, role string, userID int) ([]*models.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{
        "deleted_at": nil,
        "$or": []bson.M{
            {"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
            {"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
        },
    }
    if role != "admin" {
        filter["created_by"] = userID
    }

    sort := bson.D{{sortBy, getMongoOrder(order)}}
    opts := options.Find().
        SetSort(sort).
        SetLimit(int64(limit)).
        SetSkip(int64(offset))

    cursor, err := database.PekerjaanCollection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []*models.Pekerjaan
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    return results, nil
}

// CountPekerjaan: Tambah role/userID
func (r *PekerjaanMongoRepo) CountPekerjaan(search string, role string, userID int) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{
        "deleted_at": nil,
        "$or": []bson.M{
            {"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
            {"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
        },
    }
    if role != "admin" {
        filter["created_by"] = userID
    }

    return database.PekerjaanCollection.CountDocuments(ctx, filter)
}

// GetAlumniWithPekerjaan: OK (sudah pass isAdmin)
func (r *PekerjaanMongoRepo) GetAlumniWithPekerjaan(userID int, isAdmin bool) ([]*models.AlumniWithPekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    pipeline := mongo.Pipeline{
        bson.D{{"$match", bson.D{{"deleted_at", nil}}}},
        bson.D{{"$lookup", bson.D{
            {"from", "alumni"},
            {"localField", "alumni_id"},
            {"foreignField", "_id"},
            {"as", "alumni_data"},
        }}},
    }

    if !isAdmin {
        pipeline = append(pipeline, bson.D{{"$match", bson.D{{"created_by", userID}}}})
    }

    pipeline = append(pipeline, bson.D{{"$group", bson.D{
        {"_id", "$alumni_data._id"},
        {"alumni", bson.D{{"$first", "$alumni_data"}}},
        {"pekerjaan", bson.D{{"$push", "$$ROOT"}}},
    }}})

    cursor, err := database.PekerjaanCollection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []*models.AlumniWithPekerjaan
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    return results, nil
}

// Helper
func getMongoOrder(order string) int {
    if strings.ToLower(order) == "desc" {
        return -1
    }
    return 1
}

func New() *PekerjaanMongoRepo {
    return &PekerjaanMongoRepo{}
}