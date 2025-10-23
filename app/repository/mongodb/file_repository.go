package repository

import (
	models "alumniproject/app/models/mongodb"
	// "alumniproject/models/mongodb"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileRepository interface {
    Create(file *models.File) error
    FindAll() ([]models.File, error)
    FindByID(id int64) (*models.File, error)
    Delete(id int64) error
    GetNextID() (int64, error)
}

type fileRepository struct {
    collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) FileRepository {
    return &fileRepository{
        collection: db.Collection("files"),
    }
}

// Mendapatkan ID terakhir (auto increment)
func (r *fileRepository) GetNextID() (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

     opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
    opts.Sort = bson.D{{Key: "id", Value: -1}}

    var lastFile models.File
    err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&lastFile)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return 1, nil
        }
        return 0, err
    }

    return lastFile.ID + 1, nil
}

func (r *fileRepository) Create(file *models.File) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    nextID, err := r.GetNextID()
    if err != nil {
        return err
    }

    file.ID = nextID
    file.UploadedAt = time.Now()

    _, err = r.collection.InsertOne(ctx, file)
    return err
}

func (r *fileRepository) FindAll() ([]models.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var files []models.File
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &files); err != nil {
        return nil, err
    }

    return files, nil
}

func (r *fileRepository) FindByID(id int64) (*models.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var file models.File
    err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&file)
    if err != nil {
        return nil, err
    }

    return &file, nil
}

func (r *fileRepository) Delete(id int64) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
    return err
}
