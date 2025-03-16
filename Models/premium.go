package models

import (
	"time"
)

type PlanType string

const (
	FreePlan   PlanType = "free"
	StarPlan   PlanType = "star"
	GalaxyPlan PlanType = "galaxy"
)

// Payment represents a payment transaction
type Payment struct {
	ID      uint32 `json:"id" db:"id"`
	UserID  uint32 `json:"user_id" db:"user_id"`
	Type    string `json:"type" db:"type"`
	Gateway string `json:"gateway" db:"gateway"`
	// Mobile     string    `json:"mobile" db:"mobile_number"`
	TxnID     string    `json:"txn_id" db:"txn_id"`
	Amount    int64     `json:"amount" db:"amount"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Credit represents a user's credit balance
type Credit struct {
	ID      uint32 `json:"id" db:"id"`
	UserID  uint32 `json:"user_id" db:"user_id"`
	Balance int64  `json:"balance" db:"balance"`
	// Mobile    string    `json:"mobile" db:"mobile_number"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Subscription represents a user's subscription
type Subscription struct {
	ID        uint32    `json:"id" db:"id"`
	UserID    uint32    `json:"user_id" db:"user_id"`
	CreditID  uint32    `json:"credit_id" db:"credit_id"`
	Plan      PlanType  `json:"plan" db:"plan"`
	Price     uint32    `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Status    string    `json:"status" db:"status"`
}

type PaymentRequest struct {
	Plan   PlanType `json:"plan"`
	Months int      `json:"months"`
	// Mobile string   `json:"mobile"`
}
