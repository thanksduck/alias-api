package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thanksduck/alias-api/cfconfig"
)

type DNSRecord struct {
	Content  string   `json:"content"`
	Name     string   `json:"name"`
	Proxied  bool     `json:"proxied"`
	Type     string   `json:"type"`
	Comment  string   `json:"comment"`
	ID       string   `json:"id,omitempty"`
	Tags     []string `json:"tags"`
	Priority int      `json:"priority,omitempty"`
	TTL      int      `json:"ttl"`
}

type BatchRequest struct {
	Posts []DNSRecord `json:"posts"`
}

type BatchResponse struct {
	Result struct {
		Posts []struct {
			ID        string `json:"id"`
			ZoneID    string `json:"zone_id"`
			ZoneName  string `json:"zone_name"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			Content   string `json:"content"`
			Priority  uint16 `json:"priority,omitempty"`
			Proxiable bool   `json:"proxiable"`
			Proxied   bool   `json:"proxied"`
			TTL       uint16 `json:"ttl"`
			CreatedOn string `json:"created_on"`
		} `json:"posts"`
	} `json:"result"`
	Success bool `json:"success"`
}

type RecordResult struct {
	Type     string `json:"type,omitempty"`
	Content  string `json:"content,omitempty"`
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	TTL      uint16 `json:"ttl,omitempty"`
	Priority uint16 `json:"priority,omitempty"`
}

func CreateNewDomain(subdomain, domain string, config cfconfig.CloudflareConfig) ([]RecordResult, error) {
	records := []DNSRecord{
		{
			Content:  "route1.mx.cloudflare.net",
			Name:     fmt.Sprintf("%s.%s", subdomain, domain),
			Proxied:  false,
			Type:     "MX",
			Comment:  "mail subdomain dns Record",
			Tags:     []string{},
			Priority: 11,
			TTL:      3600,
		},
		{
			Content:  "route2.mx.cloudflare.net",
			Name:     fmt.Sprintf("%s.%s", subdomain, domain),
			Proxied:  false,
			Type:     "MX",
			Comment:  "mail subdomain dns Record",
			Tags:     []string{},
			Priority: 52,
			TTL:      3600,
		},
		{
			Content:  "route3.mx.cloudflare.net",
			Name:     fmt.Sprintf("%s.%s", subdomain, domain),
			Proxied:  false,
			Type:     "MX",
			Comment:  "mail subdomain dns Record",
			Tags:     []string{},
			Priority: 85,
			TTL:      3600,
		},
		{
			Content: `"v=spf1 include:_spf.mx.cloudflare.net ~all"`,
			Name:    fmt.Sprintf("%s.%s", subdomain, domain),
			Proxied: false,
			Type:    "TXT",
			Comment: "mail subdomain dns Record",
			Tags:    []string{},
			TTL:     3600,
		},
	}

	batchRequest := BatchRequest{Posts: records}

	jsonData, err := json.Marshal(batchRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/batch", config.ZoneID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", config.Email)
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var batchResponse BatchResponse
	err = json.Unmarshal(body, &batchResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if !batchResponse.Success {
		return nil, fmt.Errorf("batch request was not successful")
	}

	var result []RecordResult
	for _, record := range batchResponse.Result.Posts {
		result = append(result, RecordResult{
			Type:     record.Type,
			Content:  record.Content,
			ID:       record.ID,
			Name:     record.Name,
			TTL:      record.TTL,
			Priority: record.Priority,
		})
	}

	return result, nil
}
