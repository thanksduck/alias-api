package repository

import (
	"context"
	"fmt"
	"time"

	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

// SavePasswordResetToken stores or updates a password reset token
// Requires user_auth table to have a PRIMARY KEY or UNIQUE constraint on user_id
func SavePasswordResetToken(userID uint32, username string, token string) error {
	pool := db.GetPool()
	expiresAt := time.Now().Add(10 * time.Minute)

	_, err := pool.Exec(context.Background(),
		`INSERT INTO user_auth (user_id, username, password_reset_token, password_reset_expires) 
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id) DO UPDATE 
		 SET password_reset_token = $3, password_reset_expires = $4`,
		userID, username, token, expiresAt)
	if err != nil {
		return fmt.Errorf("error updating password reset token: %w", err)
	}

	return nil
}

func RemovePasswordResetToken(userID uint32) error {
	pool := db.GetPool()

	_, err := pool.Exec(context.Background(),
		`UPDATE user_auth 
		 SET password_reset_token = NULL, password_reset_expires = NULL 
		 WHERE user_id = $1`,
		userID)
	if err != nil {
		return fmt.Errorf("error removing password reset token: %w", err)
	}

	return nil
}

func FindUserByPasswordResetToken(token string) (*models.User, error) {
	pool := db.GetPool()
	var user models.User
	var passwordResetExpires time.Time

	err := pool.QueryRow(context.Background(),
		`SELECT u.id, u.username, u.name, u.email, ua.password_reset_expires 
		 FROM users u 
		 JOIN user_auth ua ON u.id = ua.user_id 
		 WHERE ua.password_reset_token = $1`, token).Scan(
		&user.ID, &user.Username, &user.Name, &user.Email, &passwordResetExpires)
	if err != nil {
		return nil, fmt.Errorf("error querying user by password reset token: %w", err)
	}

	return &user, nil
}

func IsPasswordResetTokenExpired(token string) error {
	pool := db.GetPool()
	var expiresAt time.Time

	err := pool.QueryRow(context.Background(),
		`SELECT password_reset_expires FROM user_auth WHERE password_reset_token = $1`,
		token).Scan(&expiresAt)
	if err != nil {
		return fmt.Errorf("error querying token expiration: %w", err)
	}

	if time.Now().After(expiresAt) {
		return fmt.Errorf("password reset token has expired")
	}

	return nil
}

func HasNoActiveResetToken(userID uint32) error {
	pool := db.GetPool()
	var expiresAt time.Time

	err := pool.QueryRow(context.Background(),
		`SELECT password_reset_expires FROM user_auth 
		 WHERE user_id = $1 AND password_reset_token IS NOT NULL`,
		userID).Scan(&expiresAt)

	if err != nil {
		return nil // No active token found
	}

	if time.Now().Before(expiresAt) {
		return fmt.Errorf("active reset token exists, please wait until it expires")
	}

	return nil
}
