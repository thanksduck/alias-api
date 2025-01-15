package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

func GetActivePremiumByUsername(username string) (*models.Premium, error) {
	pool := db.GetPool()
	var premium models.Premium
	err := pool.QueryRow(context.Background(),
		`SELECT id, user_id, username, subscription_id, plan, mobile, status, gateway, created_at, updated_at 
		 FROM premium WHERE username = $1 AND status = 'active'`, username).Scan(
		&premium.ID, &premium.UserID, &premium.Username, &premium.SubscriptionID,
		&premium.Plan, &premium.Mobile, &premium.Status, &premium.Gateway,
		&premium.CreatedAt, &premium.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no active subscription found for user %s", username)
		}
		return nil, fmt.Errorf("error querying premium: %w", err)
	}
	return &premium, nil
}

func GetPlanByUsername(username string) (models.PlanType, error) {
	pool := db.GetPool()
	var plan models.PlanType
	err := pool.QueryRow(context.Background(),
		`SELECT plan FROM premium WHERE username = $1 AND status = 'active'`, username).Scan(&plan)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.FreePlan, nil
		}
		return "", fmt.Errorf("error querying plan: %w", err)
	}
	return plan, nil
}

func CreateSubscription(username string, userID uint32, subscriptionID string, plan models.PlanType, mobile string, gateway models.GatewayType) error {
	pool := db.GetPool()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// First update any existing active subscriptions to inactive
	_, err = tx.Exec(context.Background(),
		`UPDATE premium SET status = 'inactive', updated_at = $1 
		 WHERE username = $2 AND status = 'active'`,
		time.Now(), username)
	if err != nil {
		return fmt.Errorf("error updating existing subscriptions: %w", err)
	}

	// Create new subscription
	_, err = tx.Exec(context.Background(),
		`INSERT INTO premium (username, user_id, subscription_id, plan, mobile, status, gateway)
		 VALUES ($1, $2, $3, $4, $5, 'active', $6)`,
		username, userID, subscriptionID, plan, mobile, gateway)
	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}

	// Update user's premium status
	_, err = tx.Exec(context.Background(),
		`UPDATE users SET is_premium = true WHERE username = $1`,
		username)
	if err != nil {
		return fmt.Errorf("error updating user premium status: %w", err)
	}

	return tx.Commit(context.Background())
}

func CancelSubscription(username string) error {
	pool := db.GetPool()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	result, err := tx.Exec(context.Background(),
		`UPDATE premium SET status = 'inactive', updated_at = $1 
		 WHERE username = $2 AND status = 'active'`,
		time.Now(), username)
	if err != nil {
		return fmt.Errorf("error canceling subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no active subscription found for user %s", username)
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE users SET is_premium = false WHERE username = $1`,
		username)
	if err != nil {
		return fmt.Errorf("error updating user premium status: %w", err)
	}

	return tx.Commit(context.Background())
}

func GetSubscriptionStatus(username string) (models.StatusType, error) {
	pool := db.GetPool()
	var status models.StatusType
	err := pool.QueryRow(context.Background(),
		`SELECT status FROM premium WHERE username = $1 ORDER BY created_at DESC LIMIT 1`,
		username).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("no subscription found for user %s", username)
		}
		return "", fmt.Errorf("error querying subscription status: %w", err)
	}
	return status, nil
}
