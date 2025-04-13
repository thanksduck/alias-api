package models

import (
	"time"
)

type User struct {
	ID                uint32    `json:"-"`
	Username          string    `json:"username,omitempty"`
	Name              string    `json:"name,omitempty"`
	Email             string    `json:"email,omitempty"`
	EmailVerified     bool      `json:"emailVerified"`
	AliasCount        int       `json:"aliasCount"`
	DestinationCount  int       `json:"destinationCount"`
	IsPremium         bool      `json:"isPremium"`
	Password          string    `json:"-"`
	Provider          string    `json:"-"`
	Avatar            string    `json:"avatar,omitempty"`
	PasswordChangedAt time.Time `json:"-"`
	Active            bool      `json:"-"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

type SafeUser struct {
	Username         string   `json:"username,omitempty"`
	Name             string   `json:"name,omitempty"`
	Email            string   `json:"email,omitempty"`
	IsEmailVerified  bool     `json:"isEmailVerified"`
	AliasCount       int64    `json:"aliasCount"`
	DestinationCount int64    `json:"destinationCount"`
	IsPremium        bool     `json:"isPremium"`
	Avatar           string   `json:"avatar,omitempty"`
	Plan             PlanType `json:"plan,omitempty"`
}
