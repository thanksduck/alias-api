-- name: InsertWebAuthnCredential :one
INSERT INTO webauthn_credentials (
    user_id,
    username,
    credential_id,
    public_key,
    sign_count,
    transports,
    authenticator_aaguid,
    is_backup
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, user_id, username, credential_id, public_key, sign_count, transports, authenticator_aaguid, is_backup, created_at, updated_at;

-- name: GetWebAuthnCredentialByID :one
SELECT id, user_id, username, credential_id, public_key, sign_count, transports, authenticator_aaguid, is_backup, created_at, updated_at
FROM webauthn_credentials
WHERE id = $1;

-- name: FindWebAuthnCredentialByUserAndCredID :one
SELECT id, user_id, username, credential_id, public_key, sign_count, transports, authenticator_aaguid, is_backup, created_at, updated_at
FROM webauthn_credentials
WHERE credential_id = $1 AND user_id=$2;

-- name: GetWebAuthnCredentialByCredentialID :one
SELECT id, user_id, username, credential_id, public_key, sign_count, transports, authenticator_aaguid, is_backup, created_at, updated_at
FROM webauthn_credentials
WHERE credential_id = $1;

-- name: ListWebAuthnCredentialsByUserID :many
SELECT id, user_id, username, credential_id, public_key, sign_count, transports, authenticator_aaguid, is_backup, created_at, updated_at
FROM webauthn_credentials
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListWebAuthnCredentialsByUsername :many
SELECT id, user_id, username, credential_id, public_key, sign_count, transports, authenticator_aaguid, is_backup, created_at, updated_at
FROM webauthn_credentials
WHERE username = $1
ORDER BY created_at DESC;

-- name: DeleteWebAuthnCredential :exec
DELETE FROM webauthn_credentials
WHERE id = $1;

-- name: UpdateWebAuthnSignCount :exec
UPDATE webauthn_credentials
SET sign_count = $2, updated_at = now()
WHERE credential_id = $1;