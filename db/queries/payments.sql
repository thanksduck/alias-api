-- queries/payments.sql

-- name: InitialisePayment :exec
INSERT INTO payments (user_id, type, gateway, txn_id, amount, status)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: FindPaymentByTxnID :one
SELECT id, user_id, type, gateway, txn_id, amount, status, created_at, updated_at
FROM payments
WHERE txn_id = $1;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: FindCreditByUserID :one
SELECT id, user_id, balance, is_active, created_at, updated_at
FROM credits
WHERE user_id = $1;

-- name: CreateCredit :one
INSERT INTO credits (user_id, balance, is_active)
VALUES ($1, $2, true)
RETURNING id;

-- name: UpdateCreditBalance :exec
UPDATE credits
SET balance = balance + $1,
    updated_at = NOW()
WHERE id = $2;

-- name: CreateSubscription :exec
INSERT INTO subscriptions (user_id, credit_id, plan, price, expires_at, status)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUserToPremium :exec
UPDATE users
SET is_premium = true
WHERE id = $1;

-- name: GetPlanByUserID :one
SELECT plan
FROM subscriptions
WHERE user_id = $1;

-- name: GetSubscriptionByUserID :one
SELECT plan, price, created_at, updated_at, expires_at, status
FROM subscriptions
WHERE user_id = $1;