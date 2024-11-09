package repository

import (
	"context"
	"fmt"

	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

func FindDestinationByID(id uint32) (*models.Destination, error) {
	pool := db.GetPool()
	var destination models.Destination
	err := pool.QueryRow(context.Background(),
		`SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, verified 
		FROM destinations WHERE id = $1`, id).Scan(
		&destination.ID,
		&destination.UserID,
		&destination.Username,
		&destination.DestinationEmail,
		&destination.Domain,
		&destination.CloudflareDestinationID,
		&destination.Verified,
	)
	if err != nil {
		return nil, err
	}
	return &destination, nil
}

func FindDestinationByEmail(email string) (*models.Destination, error) {
	pool := db.GetPool()
	var destination models.Destination
	err := pool.QueryRow(context.Background(),
		`SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, verified 
		FROM destinations WHERE destination_email = $1`, email).Scan(
		&destination.ID,
		&destination.UserID,
		&destination.Username,
		&destination.DestinationEmail,
		&destination.Domain,
		&destination.CloudflareDestinationID,
		&destination.Verified,
	)
	if err != nil {
		return nil, err
	}
	return &destination, nil
}

func FindDestinationsByUserID(userID uint32) ([]models.Destination, error) {
	pool := db.GetPool()
	rows, err := pool.Query(context.Background(),
		`SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, verified 
		FROM destinations WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var destinations []models.Destination
	for rows.Next() {
		var destination models.Destination
		err := rows.Scan(
			&destination.ID,
			&destination.UserID,
			&destination.Username,
			&destination.DestinationEmail,
			&destination.Domain,
			&destination.CloudflareDestinationID,
			&destination.Verified,
		)
		if err != nil {
			return nil, err
		}
		destinations = append(destinations, destination)
	}
	return destinations, nil
}
func FindDestinationByEmailAndUsername(email, username string) (*models.Destination, error) {
	pool := db.GetPool()
	var destination models.Destination
	err := pool.QueryRow(context.Background(),
		`SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, verified 
		FROM destinations WHERE destination_email = $1 AND username = $2`, email, username).Scan(
		&destination.ID,
		&destination.UserID,
		&destination.Username,
		&destination.DestinationEmail,
		&destination.Domain,
		&destination.CloudflareDestinationID,
		&destination.Verified,
	)
	if err != nil {
		return nil, err
	}
	return &destination, nil
}

func FindDestinationsByUsername(username string) ([]models.Destination, error) {
	pool := db.GetPool()
	rows, err := pool.Query(context.Background(),
		`SELECT id, user_id, username, destination_email, domain, cloudflare_destination_id, verified 
		FROM destinations WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var destinations []models.Destination
	for rows.Next() {
		var destination models.Destination
		err := rows.Scan(
			&destination.ID,
			&destination.UserID,
			&destination.Username,
			&destination.DestinationEmail,
			&destination.Domain,
			&destination.CloudflareDestinationID,
			&destination.Verified,
		)
		if err != nil {
			return nil, err
		}
		destinations = append(destinations, destination)
	}
	return destinations, nil
}

func CreateDestination(destination *models.Destination) (*models.Destination, error) {
	pool := db.GetPool()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(),
		`INSERT INTO destinations (user_id, username, destination_email, domain, cloudflare_destination_id, verified)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, user_id, username, destination_email, domain, cloudflare_destination_id, verified`,
		destination.UserID,
		destination.Username,
		destination.DestinationEmail,
		destination.Domain,
		destination.CloudflareDestinationID,
		destination.Verified,
	).Scan(
		&destination.ID,
		&destination.UserID,
		&destination.Username,
		&destination.DestinationEmail,
		&destination.Domain,
		&destination.CloudflareDestinationID,
		&destination.Verified,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE users SET destination_count = destination_count + 1 WHERE id = $1`, destination.UserID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return destination, nil
}

func UpdateDestination(destination *models.Destination) (*models.Destination, error) {
	pool := db.GetPool()
	err := pool.QueryRow(context.Background(),
		`UPDATE destinations SET username = $1, destination_email = $2, domain = $3, cloudflare_destination_id = $4, verified = $5 WHERE id = $6 RETURNING id, user_id, username, destination_email, domain, cloudflare_destination_id, verified`,
		destination.Username,
		destination.DestinationEmail,
		destination.Domain,
		destination.CloudflareDestinationID,
		destination.Verified,
		destination.ID,
	).Scan(
		&destination.ID,
		&destination.UserID,
		&destination.Username,
		&destination.DestinationEmail,
		&destination.Domain,
		&destination.CloudflareDestinationID,
		&destination.Verified,
	)
	if err != nil {
		return nil, err
	}
	return destination, nil
}

func VerifyDestinationByID(id uint32) error {
	pool := db.GetPool()
	_, err := pool.Exec(context.Background(),
		`UPDATE destinations SET verified = true WHERE id = $1`, id)
	return err
}

func DeleteDestinationByID(id uint32, userID uint32) error {
	pool := db.GetPool()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	result, err := tx.Exec(context.Background(),
		`DELETE FROM destinations WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}

	// Check if any row was actually deleted
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no destination found with id %d for user %d", id, userID)
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE users SET destination_count = destination_count - 1 WHERE id = $1`, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
