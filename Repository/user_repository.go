package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

func CreateUser(user *models.User) (*models.User, error) {
	pool := db.GetPool()
	err := pool.QueryRow(context.Background(),
		`INSERT INTO users (username, name, email, password, created_at, updated_at, provider, avatar)
		 VALUES ($1, $2, $3, $4, $5, $6,$7,$8) RETURNING id, username, name, email, alias_count, destination_count, provider, avatar`,
		user.Username, user.Name, user.Email, user.Password, time.Now(), time.Now(), user.Provider, user.Avatar).Scan(
		&user.ID, &user.Username, &user.Name, &user.Email, &user.AliasCount,
		&user.DestinationCount, &user.Provider, &user.Avatar)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return user, nil
}

func CreateOrUpdateUser(user *models.User) (*models.User, error) {
	pool := db.GetPool()
	err := pool.QueryRow(context.Background(),
		`INSERT INTO users (email, username, name, email_verified, provider, avatar, password, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (email) DO UPDATE SET avatar = $6, password = $7, updated_at = $9 RETURNING id, username, name, email, alias_count, destination_count, provider, avatar`,
		user.Email, user.Username, user.Name, user.EmailVerified, user.Provider, user.Avatar, user.Password, time.Now(), time.Now()).Scan(
		&user.ID, &user.Username, &user.Name, &user.Email, &user.AliasCount,
		&user.DestinationCount, &user.Provider, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUserByID(id uint32) (*models.User, error) {
	pool := db.GetPool()
	var user models.User
	var passwordChangedAt sql.NullTime

	err := pool.QueryRow(context.Background(),
		`SELECT id, username, name, email, alias_count, destination_count, is_premium,
         provider, avatar, password_changed_at, active, password, email_verified
         FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.Username, &user.Name, &user.Email, &user.AliasCount,
		&user.DestinationCount, &user.IsPremium, &user.Provider, &user.Avatar,
		&passwordChangedAt, &user.Active, &user.Password, &user.EmailVerified)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	if passwordChangedAt.Valid {
		user.PasswordChangedAt = passwordChangedAt.Time
	} else {
		user.PasswordChangedAt = time.Time{}
	}

	return &user, nil
}

func FindUserByUsernameOrEmail(username string, email string) (*models.User, error) {
	pool := db.GetPool()
	var user models.User
	var passwordChangedAt sql.NullTime

	err := pool.QueryRow(context.Background(),
		`SELECT id, username, name, email, alias_count, destination_count, is_premium,
         provider, avatar, password_changed_at, active, password, email_verified
         FROM users WHERE username = $1 OR email = $2`, username, email).Scan(
		&user.ID, &user.Username, &user.Name, &user.Email, &user.AliasCount,
		&user.DestinationCount, &user.IsPremium, &user.Provider, &user.Avatar,
		&passwordChangedAt, &user.Active, &user.Password, &user.EmailVerified)
	if err != nil {
		return nil, err
	}

	if passwordChangedAt.Valid {
		user.PasswordChangedAt = passwordChangedAt.Time
	} else {
		user.PasswordChangedAt = time.Time{}
	}

	return &user, nil
}

func UpdatePassword(id uint32, password string) error {
	pool := db.GetPool()

	// Start a transaction since we're updating two tables
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Update password in users table
	_, err = tx.Exec(context.Background(),
		`UPDATE users SET password = $1, password_changed_at = $2 WHERE id = $3`,
		password, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	// Update password related fields in user_auth table
	_, err = tx.Exec(context.Background(),
		`UPDATE user_auth SET  password_reset_token = NULL, password_reset_expires = NULL 
		 WHERE user_id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("error updating password auth fields: %w", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// func DeletePasswordResetToken(id uint32) error {
// 	pool := db.GetPool()
// 	_, err := pool.Exec(context.Background(),
// 		`UPDATE users SET password_reset_token = '', password_reset_expires = NULL WHERE id = $1`, id)
// 	if err != nil {
// 		return fmt.Errorf("error deleting password reset token: %w", err)
// 	}
// 	return nil
// }

func UpdateUser(id uint32, user *models.User) (*models.User, error) {
	pool := db.GetPool()
	_, err := pool.Exec(context.Background(),
		`UPDATE users SET name = $1, email = $2, avatar = $3, username = $4, provider = $5, email_verified = $6 WHERE id = $7`,
		user.Name, user.Email, user.Avatar, user.Username, user.Provider, user.EmailVerified, id)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}

func DeleteUser(id uint32) error {
	pool := db.GetPool()
	_, err := pool.Exec(context.Background(),
		`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func VerifyEmailByID(id uint32) error {
	pool := db.GetPool()
	_, err := pool.Exec(context.Background(),
		`UPDATE users SET email_verified = true WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error verifying email: %w", err)
	}
	return nil
}

func UpdateProviderByID(id uint32, provider string) error {
	pool := db.GetPool()
	_, err := pool.Exec(context.Background(),
		`UPDATE users SET provider = $1 WHERE id = $2`, provider, id)
	if err != nil {
		return fmt.Errorf("error updating provider: %w", err)
	}
	return nil
}
