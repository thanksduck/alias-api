package models

import "time"

type Rule struct {
	ID               uint32    `json:"ruleId"`
	UserID           uint32    `json:"-"`
	Username         string    `json:"username"`
	AliasEmail       string    `json:"aliasEmail"`
	DestinationEmail string    `json:"destinationEmail"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
	Comment          string    `json:"comment,omitempty"`
	Name             string    `json:"name,omitempty"`
}
