package repository

import (
	"context"

	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
)

func CreateNewSubDomain(customDomain *models.CustomDomain) (uint32, error) {
	dbpool := db.GetPool()
	var id uint32
	err := dbpool.QueryRow(context.Background(), `INSERT INTO custom_domains (user_id, domain, username) VALUES ($1, $2, $3) RETURNING id`, customDomain.UserID, customDomain.Domain, customDomain.Username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func CreateNewDNSRecord(dnsRecord *models.CustomDomainDNSRecord) error {
	dbpool := db.GetPool()
	_, err := dbpool.Exec(context.Background(), `INSERT INTO custom_domain_dns_records (custom_domain_id, cloudflare_id, type, name, content, ttl, priority) VALUES ($1, $2, $3, $4, $5, $6, $7)`, dnsRecord.CustomDomainID, dnsRecord.CloudflareID, dnsRecord.Type, dnsRecord.Name, dnsRecord.Content, dnsRecord.TTL, dnsRecord.Priority)
	if err != nil {
		return err
	}
	return nil
}
