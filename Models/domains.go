package models

import "time"

type CustomDomain struct {
	ID       uint32 `json:"id"`
	UserID   uint32 `json:"-"`
	Username string `json:"username"`
	Domain   string `json:"domain"`
}

type CustomDomainDNSRecord struct {
	ID             uint32    `json:"id"`
	CustomDomainID uint32    `json:"custom_domain_id"`
	CloudflareID   string    `json:"cloudflare_id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Content        string    `json:"content"`
	TTL            uint16    `json:"ttl"`
	Priority       uint16    `json:"priority,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
