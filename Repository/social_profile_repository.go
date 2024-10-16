package repository

import (
	"context"
	"strings"
	"time"

	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

func CreateOrUpdateSocialProfile(socialProfile *models.SocialProfile) (*models.SocialProfile, error) {
	poll := db.GetPool()

	// Prepare the query and arguments
	query := `INSERT INTO social_profiles (user_id, username, google, facebook, github, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) 
			  ON CONFLICT (user_id) DO UPDATE SET `
	args := []interface{}{socialProfile.UserID, socialProfile.Username, socialProfile.Google, socialProfile.Facebook, socialProfile.Github, time.Now(), time.Now()}

	// Check for null values and adjust the query accordingly
	updateFields := []string{}
	if socialProfile.Google != "" {
		updateFields = append(updateFields, "google = EXCLUDED.google")
	}
	if socialProfile.Facebook != "" {
		updateFields = append(updateFields, "facebook = EXCLUDED.facebook")
	}
	if socialProfile.Github != "" {
		updateFields = append(updateFields, "github = EXCLUDED.github")
	}
	updateFields = append(updateFields, "updated_at = EXCLUDED.updated_at")

	query += strings.Join(updateFields, ", ") + " RETURNING id, user_id, username, google, facebook, github, created_at, updated_at"

	// Execute the query
	err := poll.QueryRow(context.Background(), query, args...).Scan(
		&socialProfile.ID, &socialProfile.UserID, &socialProfile.Username, &socialProfile.Google, &socialProfile.Facebook, &socialProfile.Github, &socialProfile.CreatedAt, &socialProfile.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return socialProfile, nil
}

func FindSocialProfileByIDOrUsername(userID uint32, username string) (*models.SocialProfile, error) {
	pool := db.GetPool()
	var socialProfile models.SocialProfile

	var query string
	var args []interface{}
	if userID != 0 {
		query = `SELECT id, user_id, username, google, facebook, github, created_at, updated_at
				 FROM social_profiles WHERE user_id = $1`
		args = []interface{}{userID}
	} else {
		query = `SELECT id, user_id, username, google, facebook, github, created_at, updated_at
				 FROM social_profiles WHERE username = $1`
		args = []interface{}{username}
	}

	err := pool.QueryRow(context.Background(), query, args...).Scan(
		&socialProfile.ID, &socialProfile.UserID, &socialProfile.Username, &socialProfile.Google, &socialProfile.Facebook, &socialProfile.Github, &socialProfile.CreatedAt, &socialProfile.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &socialProfile, nil
}
