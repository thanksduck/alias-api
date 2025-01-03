package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thanksduck/alias-api/cfconfig"
)

type DestinationResult struct {
	Created  time.Time `json:"created"`
	Email    string    `json:"email"`
	ID       string    `json:"id"`
	Modified time.Time `json:"modified"`
	Tag      string    `json:"tag"`
	Verified time.Time `json:"verified"`
}

type DestinationResponse struct {
	Errors   []interface{}     `json:"errors"`
	Messages []interface{}     `json:"messages"`
	Success  bool              `json:"success"`
	Result   DestinationResult `json:"result"`
}

func DestinationRequest(method, domain, destination, cfId string) (*DestinationResponse, error) {
	config, ok := cfconfig.SelectDomain(domain)
	if !ok {
		return nil, fmt.Errorf("no configuration found for domain: %s", domain)
	}

	urlPrefix := "https://api.cloudflare.com/client/v4"
	url := fmt.Sprintf("%s/accounts/%s/email/routing/addresses", urlPrefix, config.AccountID)
	if cfId != "" {
		url += "/" + cfId
	}

	payload := map[string]string{"email": destination}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var destResponse DestinationResponse
	err = json.Unmarshal(body, &destResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &destResponse, nil
}
