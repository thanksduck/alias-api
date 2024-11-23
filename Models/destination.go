package models

import "time"

type Destination struct {
	ID                      uint32    `json:"destinationID" db:"id"`
	UserID                  uint32    `json:"-" db:"user_id"`
	Username                string    `json:"username" db:"username"`
	DestinationEmail        string    `json:"destinationEmail" db:"destination_email"`
	CloudflareDestinationID string    `json:"-" db:"cloudflare_destination_id"`
	Domain                  string    `json:"domain" db:"domain"`
	Verified                bool      `json:"verified" db:"verified"`
	CreatedAt               time.Time `json:"-" db:"created_at"`
	UpdatedAt               time.Time `json:"-" db:"updated_at"`
}
