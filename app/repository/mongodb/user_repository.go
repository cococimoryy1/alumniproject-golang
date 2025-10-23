package repository

import (
	"context"
	"log"
	"time"

	"alumniproject/app/models/mongodb"
	db "alumniproject/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ambil user berdasarkan username atau email
func GetUserByUsernameOrEmail(usernameOrEmail string) (models.User, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	filter := bson.M{
		"$or": []bson.M{
			{"username": usernameOrEmail},
			{"email": usernameOrEmail},
		},
	}

	err := db.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, "", err
		}
		log.Println("FindOne error:", err)
		return models.User{}, "", err
	}

	return user, user.Password, nil
}

// Ambil daftar user dengan fitur search, sort, pagination
func GetUsersRepo(search, sortBy, order string, limit, offset int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}

	opts := options.Find().
		SetSort(bson.M{sortBy: sortOrder}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := db.DB.Collection("users").Find(ctx, filter, opts)
	if err != nil {
		log.Println("Find error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return users, nil
}

// Hitung total user untuk pagination
func CountUsersRepo(search string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	count, err := db.DB.Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		log.Println("CountDocuments error:", err)
		return 0, err
	}

	return int(count), nil
}
