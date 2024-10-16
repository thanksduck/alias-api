package models

import (
	"time"
)

type User struct {
	ID                   uint32    `json:"-"`
	Username             string    `json:"username,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Email                string    `json:"email,omitempty"`
	EmailVerified        bool      `json:"emailVerified"`
	AliasCount           int       `json:"aliasCount"`
	DestinationCount     int       `json:"destinationCount"`
	IsPremium            bool      `json:"-"`
	Password             string    `json:"-"`
	Provider             string    `json:"-"`
	Avatar               string    `json:"avatar,omitempty"`
	PasswordChangedAt    time.Time `json:"-"`
	PasswordResetToken   string    `json:"-"`
	PasswordResetExpires time.Time `json:"-"`
	Active               bool      `json:"-"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"-"`
}

func (u *User) IsPasswordChangedAfter(unixTime int64) bool {
	if u.PasswordChangedAt.IsZero() {
		return false
	}
	// Adjust for the timezone offset of 5 hours and 30 minutes (19800 seconds)
	passwordChangedAtUnix := u.PasswordChangedAt.Unix() - 19800

	return passwordChangedAtUnix > unixTime
}
