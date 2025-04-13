-- queries/rules.sql

-- name: FindRuleByID :one
SELECT id as rule_id, username, alias_email, destination_email, is_active, comment, name
FROM rules
WHERE id = $1;

-- name: FindRulesByUserID :many
SELECT id as rule_id, username, alias_email, destination_email, is_active, comment, name
FROM rules
WHERE user_id = $1;

-- name: FindRulesByDestinationEmail :many
SELECT id, user_id, username, alias_email, destination_email, is_active, comment, name
FROM rules
WHERE destination_email = $1;

-- name: FindActiveRulesByDestinationEmail :many
SELECT id, user_id, username, alias_email, destination_email, is_active, comment, name
FROM rules
WHERE destination_email = $1 AND is_active = true;

-- name: MakeAllRuleInactiveByDestinationEmail :exec
UPDATE rules
SET is_active = false
WHERE destination_email = $1;

-- name: FindRuleByAliasEmail :one
SELECT id, user_id, username, alias_email, destination_email, is_active, comment
FROM rules
WHERE alias_email = $1;

-- name: CreateNewRule :exec
INSERT INTO rules (user_id, username, alias_email, destination_email, comment, name)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: IncrementUserAliasCount :exec
UPDATE users
SET alias_count = alias_count + 1
WHERE id = $1;

-- name: UpdateRuleByID :exec
UPDATE rules
SET alias_email = $1,
    destination_email = $2,
    comment = $3,
    name = $4
WHERE id = $5;

-- name: ToggleRuleByID :exec
UPDATE rules
SET is_active = NOT is_active
WHERE id = $1;

-- name: DeleteRuleByID :exec
DELETE FROM rules
WHERE id = $1;

-- name: DecrementUserAliasCount :exec
UPDATE users
SET alias_count = alias_count - 1
WHERE id = $1;
