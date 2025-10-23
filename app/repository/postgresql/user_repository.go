package repository

import (
	"context"
	"time"
	"log"
	"database/sql"
	"fmt"


	"alumniproject/app/models/postgresql"
	"alumniproject/database/postgresql"
)

// GetUserByUsernameOrEmail retrieves user and password hash for login
func GetUserByUsernameOrEmail(usernameOrEmail string) (models.User, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	var passwordHash string

	err := postgresql.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, role, created_at 
		FROM users 
		WHERE username = $1 OR email = $1
	`, usernameOrEmail).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.CreatedAt,
	)

	return user, passwordHash, err
}
// GetUsersRepo -> ambil data users dari DB
// GetUsersRepo -> ambil data users dari DB
func GetUsersRepo(search, sortBy, order string, limit, offset int) ([]models.User, error) {
    query := fmt.Sprintf(`
        SELECT id, name, email, created_at
        FROM users
        WHERE name ILIKE $1 OR email ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := postgresql.DB.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        log.Println("Query error:", err)
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

// CountUsersRepo -> hitung total data untuk pagination
func CountUsersRepo(search string) (int, error) {
    var total int
    countQuery := `SELECT COUNT(*) FROM users WHERE name ILIKE $1 OR email ILIKE $1`
    err := postgresql.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
    if err != nil && err != sql.ErrNoRows {
        return 0, err
    }
    return total, nil
}