package models

import "time"

type UserAuth struct {
	ID                   uint32    `json:"-"`
	UserID               uint32    `json:"-"`
	Username             string    `json:"username,omitempty"`
	PasswordResetToken   string    `json:"-"`
	PasswordResetExpires time.Time `json:"-"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"-"`
}
