-- queries/destinations.sql

-- name: FindDestinationByID :one
SELECT id as destination_id, username, destination_email, domain, is_verified
FROM destinations
WHERE id = $1;

-- name: GetCloudflareDestinationID :one
SELECT cloudflare_destination_id,domain,destination_email,is_verified
FROM destinations
where id = $1 and user_id = $2;

-- name: FindDestinationByEmail :one
SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, is_verified
FROM destinations
WHERE destination_email = $1;

-- name: FindDestinationByEmailAndDomain :one
SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, is_verified
FROM destinations
WHERE destination_email = $1 AND domain = $2;

-- name: FindDestinationByEmailAndDomainAndUserID :one
SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, is_verified
FROM destinations
WHERE destination_email = $1 AND domain = $2 AND user_id = $3;

-- name: FindDestinationsByUserID :many
SELECT id as destination_id, username, destination_email, domain, is_verified
FROM destinations
WHERE user_id = $1;

-- name: FindDestinationByEmailAndUsername :one
SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, is_verified
FROM destinations
WHERE destination_email = $1 AND username = $2;

-- name: CreateDestination :exec
INSERT INTO destinations (user_id, username, destination_email, domain, cloudflare_destination_id, is_verified)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: IncrementUserDestinationCount :exec
UPDATE users
SET destination_count = destination_count + 1
WHERE id = $1;

-- name: UpdateDestination :exec
UPDATE destinations
SET username = $2,
    destination_email = $3,
    domain = $4,
    cloudflare_destination_id = $5,
    is_verified = $6
WHERE id = $1;

-- name: VerifyDestinationByID :exec
UPDATE destinations
SET is_verified = true
WHERE id = $1;

-- name: DeleteDestinationByID :exec
DELETE FROM destinations
WHERE id = $1;

-- name: DecrementUserDestinationCount :exec
UPDATE users
SET destination_count = destination_count - 1
WHERE id = $1;
