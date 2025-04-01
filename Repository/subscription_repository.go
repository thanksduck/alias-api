package repository

import (
	"context"

	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

func InitialisePayment(payment *models.Payment) error {
	pool := db.GetPool()
	_, err := pool.Exec(context.Background(),
		`INSERT INTO payments (user_id, type, gateway, txn_id, amount, status) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		payment.UserID,
		payment.Type,
		payment.Gateway,
		payment.TxnID,
		payment.Amount,
		payment.Status,
	)
	if err != nil {
		return err
	}
	return nil
}

func FindPaymentByTxnID(txnID string) (*models.Payment, error) {
	pool := db.GetPool()
	var payment models.Payment
	err := pool.QueryRow(context.Background(),
		`SELECT id, user_id, type, gateway, txn_id,  amount, status, created_at, updated_at 
		FROM payments WHERE txn_id = $1`, txnID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Type,
		&payment.Gateway,
		&payment.TxnID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func UpdatePaymentStatusCreditAndCreateSubscription(subscription *models.Subscription, payment *models.Payment, status string) error {
	pool := db.GetPool()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Update payment status
	_, err = tx.Exec(context.Background(),
		`UPDATE payments SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, payment.ID)
	if err != nil {
		return err
	}

	// Check if credit record exists for user
	var creditID uint32
	var exists bool
	err = tx.QueryRow(context.Background(),
		`SELECT id, true FROM credits WHERE user_id = $1`,
		payment.UserID).Scan(&creditID, &exists)

	if err != nil && err.Error() != "no rows in result set" {
		return err
	}

	// Create or update credits
	if exists {
		// Update existing credit
		_, err = tx.Exec(context.Background(),
			`UPDATE credits SET balance = balance + $1, updated_at = NOW() 
			WHERE id = $2`,
			payment.Amount, creditID)
		subscription.CreditID = creditID
	} else {
		// Create new credit
		err = tx.QueryRow(context.Background(),
			`INSERT INTO credits (user_id, balance, is_active) 
			VALUES ($1, $2, true) RETURNING id`,
			payment.UserID, payment.Amount).Scan(&subscription.CreditID)
	}
	if err != nil {
		return err
	}

	// Create subscription
	_, err = tx.Exec(context.Background(),
		`INSERT INTO subscriptions (user_id, credit_id, plan, price, expires_at, status) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		subscription.UserID,
		subscription.CreditID,
		subscription.Plan,
		subscription.Price,
		subscription.ExpiresAt,
		subscription.Status)
	if err != nil {
		return err
	}

	// Update user to premium status
	_, err = tx.Exec(context.Background(),
		`UPDATE users SET is_premium = true WHERE id = $1`,
		subscription.UserID)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func GetPlanByUserID(UserID uint32) (*models.PlanType, error) {
	pool := db.GetPool()
	var plan models.PlanType
	err := pool.QueryRow(context.Background(), `SELECT plan from subscriptions where user_id = $1`, UserID).Scan(&plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func GetSubscriptionByUserID(UserID uint32) (*models.Subscription, error) {
	pool := db.GetPool()
	var subscription models.Subscription
	err := pool.QueryRow(context.Background(),
		`SELECT  plan, price, created_at, updated_at, expires_at, status 
		 FROM subscriptions 
		 WHERE user_id = $1 
`,
		UserID).Scan(
		&subscription.Plan,
		&subscription.Price,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.ExpiresAt,
		&subscription.Status,
	)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}
