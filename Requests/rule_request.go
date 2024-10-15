package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var (
	ruleURL    string
	ruleAPIKey string
	once       sync.Once
)

func initializeRuleConfig() {
	ruleURL = os.Getenv("RULE_URL_PREFIX") + "/rules"
	ruleAPIKey = os.Getenv("RULE_API_KEY")
}

type RuleData struct {
	Alias       string `json:"alias"`
	Destination string `json:"destination"`
	Username    string `json:"username"`
	Domain      string `json:"domain"`
	Comment     string `json:"comment"`
}

func CreateRuleRequest(method, alias, destination, username, domain string) error {
	once.Do(initializeRuleConfig)

	data := RuleData{
		Alias:       alias,
		Destination: destination,
		Username:    username,
		Domain:      domain,
		Comment:     fmt.Sprintf("Created from alias-api v2 by %s", username),
	}

	var url string
	switch method {
	case "POST":
		url = ruleURL
	case "PATCH":
		url = fmt.Sprintf("%s/%s/%s/flip", ruleURL, domain, alias)
	default:
		url = fmt.Sprintf("%s/%s/%s", ruleURL, domain, alias)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ruleAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
