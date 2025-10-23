package models

import (
	"github.com/golang-jwt/jwt/v5"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// Struktur data yang tersimpan di MongoDB
type User struct {
	ID       int    			`json:"id" bson:"id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

// Request body yang dikirim dari client saat login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response saat login berhasil
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// Struktur claim untuk JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`  // <- ubah dari primitive.ObjectID ke string
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
