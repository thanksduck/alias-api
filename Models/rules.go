package models

import "time"

type Rule struct {
	ID               uint32    `json:"id"`
	UserID           uint32    `json:"-"`
	Username         string    `json:"username"`
	AliasEmail       string    `json:"aliasEmail"`
	DestinationEmail string    `json:"destinationEmail"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"createdAt,omitempty"`
	UpdatedAt        time.Time `json:"updatedAt,omitempty"`
	Comment          string    `json:"comment,omitempty"`
	Name             string    `json:"name,omitempty"`
}
