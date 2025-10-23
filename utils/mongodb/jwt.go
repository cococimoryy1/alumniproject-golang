package mongodbutils

import (
    "fmt"
    "alumniproject/app/models/mongodb"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-min-32-characters-long") // Ganti dengan kunci rahasia yang aman, sebaiknya dari variabel lingkungan

// GenerateToken menghasilkan token JWT untuk user
func GenerateToken(user models.User) (string, error) {
	claims := models.JWTClaims{
		UserID:   fmt.Sprintf("%d", user.ID),  // <- konversi int ke string
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ValidateToken memvalidasi token JWT dan mengembalikan claims
func ValidateToken(tokenString string) (*models.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrInvalidKey
}