package repository

import (
    "context"
    "alumniproject/app/models/mongodb"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

type FotoRepository interface {
    Create(file *models.File) error
    FindAll() ([]models.File, error)
    FindByID(id int64) (*models.File, error)
    Delete(id int64) error
    GetNextID() (int64, error)
}

type fotoRepository struct {
    collection *mongo.Collection
}

func NewFotoRepository(db *mongo.Database) FotoRepository {
    return &fotoRepository{
        collection: db.Collection("fotos"),
    }
}

func (r *fotoRepository) GetNextID() (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

     opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
    opts.Sort = bson.D{{Key: "id", Value: -1}}

    var lastFoto models.File
    err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&lastFoto)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return 1, nil
        }
        return 0, err
    }

    return lastFoto.ID + 1, nil
}

func (r *fotoRepository) Create(file *models.File) error {
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

func (r *fotoRepository) FindAll() ([]models.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var fotos []models.File
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &fotos); err != nil {
        return nil, err
    }

    return fotos, nil
}

func (r *fotoRepository) FindByID(id int64) (*models.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var foto models.File
    err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&foto)
    if err != nil {
        return nil, err
    }

    return &foto, nil
}

func (r *fotoRepository) Delete(id int64) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
    return err
}
