package models

import "time"

type SocialProfile struct {
	ID        uint32    `json:"-"`
	UserID    uint32    `json:"-"`
	Username  string    `json:"username,omitempty"`
	Google    string    `json:"google,omitempty"`
	Facebook  string    `json:"facebook,omitempty"`
	Github    string    `json:"github,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
