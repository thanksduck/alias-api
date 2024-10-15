package models

import "time"

type Destination struct {
	ID                      uint32    `json:"destinationID"`
	UserID                  uint32    `json:"-"`
	Username                string    `json:"username"`
	DestinationEmail        string    `json:"destinationEmail"`
	CloudflareDestinationID string    `json:"-"`
	Domain                  string    `json:"domain"`
	Verified                bool      `json:"verified"`
	CreatedAt               time.Time `json:"-"`
	UpdatedAt               time.Time `json:"-"`
}
