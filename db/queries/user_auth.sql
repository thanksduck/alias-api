-- queries/user_auth.sql

-- name: SavePasswordResetToken :exec
INSERT INTO user_auth (user_id, username, password_reset_token, password_reset_expires)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
SET password_reset_token = $3,
    password_reset_expires = $4;

-- name: RemovePasswordResetToken :exec
UPDATE user_auth
SET password_reset_token = NULL,
    password_reset_expires = NULL
WHERE user_id = $1;

-- name: FindUserByPasswordResetToken :one
SELECT u.id,
       u.username,
       u.name,
       u.email,
       ua.password_reset_expires
FROM users u
    JOIN user_auth ua ON u.id = ua.user_id
WHERE ua.password_reset_token = $1;

-- name: GetPasswordResetTokenExpiry :one
SELECT password_reset_expires
FROM user_auth
WHERE password_reset_token = $1;

-- name: GetActivePasswordResetTokenExpiry :one
SELECT password_reset_expires
FROM user_auth
WHERE user_id = $1
  AND password_reset_token IS NOT NULL;