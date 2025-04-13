-- queries/users.sql

-- name: CreateUser :exec
INSERT INTO users (username, name, email, password, avatar)
VALUES ($1, $2, $3, $4, $5);

-- name: CreateOrUpdateUser :one
INSERT INTO users (email, username, name, is_email_verified, provider, avatar, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (email) DO UPDATE SET avatar = $6, password = $7, updated_at = $9
RETURNING id, username, name, email, alias_count, destination_count, provider, avatar;

-- name: FindUserByID :one
SELECT id, username, name, email, alias_count, destination_count, is_premium,
       provider, avatar, password_changed_at, active,  is_email_verified, created_at,updated_at
FROM users
WHERE id = $1;

-- name: FindUserByUsername :one
SELECT id, username, name, email, alias_count, destination_count, is_premium,
       provider, avatar, password_changed_at, active,  is_email_verified, created_at,updated_at
FROM users
WHERE username = $1;

-- name: FindPasswordById :one
SELECT password
FROM users
WHERE id = $1;

-- name: FindUserByUsernameOrEmail :one
SELECT id, username, name, email, alias_count, destination_count, is_premium,
       provider, avatar, password_changed_at, active, password, is_email_verified, created_at, updated_at
FROM users
WHERE username = $1 OR email = $2;

-- name: UpdatePasswordUser :exec
UPDATE users
SET password = $1,
    password_changed_at = $2
WHERE id = $3;

-- name: UpdatePasswordAuth :exec
UPDATE user_auth
SET password_reset_token = NULL,
    password_reset_expires = NULL
WHERE user_id = $1;

-- name: HasNoActiveResetToken :one
SELECT id
FROM user_auth
WHERE user_id = $1 AND password_reset_expires > now();

-- name: FindUserByValidResetToken :one
SELECT user_id
FROM user_auth
WHERE password_reset_token = $1
  AND password_reset_expires > now();

-- name: CreateNewPasswordResetToken :exec
INSERT INTO user_auth (user_id, username, password_reset_token)
VALUES ($1, $2, $3);

-- name: UpdateUser :exec
UPDATE users
SET name = $2,
    email = $3,
    avatar = $4,
    username = $5,
    provider = $6
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: VerifyEmailByID :exec
UPDATE users
SET is_email_verified = true
WHERE id = $1;

-- name: UpdateProviderByID :exec
UPDATE users
SET provider = $1
WHERE id = $2;
