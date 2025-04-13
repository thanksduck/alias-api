-- queries/social_profiles.sql

-- name: CreateOrUpdateSocialProfile :one
INSERT INTO social_profiles (user_id, username, google, facebook, github, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id) DO UPDATE SET
    google = EXCLUDED.google,
    facebook = EXCLUDED.facebook,
    github = EXCLUDED.github,
    updated_at = EXCLUDED.updated_at
RETURNING id, user_id, username, google, facebook, github, created_at, updated_at;

-- name: FindSocialProfileByUserID :one
SELECT id, user_id, username, google, facebook, github, created_at, updated_at
FROM social_profiles
WHERE user_id = $1;

-- name: FindSocialProfileByUsername :one
SELECT id, user_id, username, google, facebook, github, created_at, updated_at
FROM social_profiles
WHERE username = $1;
