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

	// Handle specific status codes
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil
	case http.StatusConflict: // 409 Conflict
		body, _ := io.ReadAll(resp.Body)
		var errorResp struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}

		// Try to parse the error response
		jsonErr := json.Unmarshal(body, &errorResp)
		if jsonErr == nil && !errorResp.Success {
			return fmt.Errorf(errorResp.Error)
		}

		fmt.Println("Error response:", string(body))
		return fmt.Errorf("something went wrong")
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
}
