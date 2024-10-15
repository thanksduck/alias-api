package cfconfig

import (
	"encoding/json"
	"strings"
	"sync"
)

type CloudflareConfig struct {
	AccountID string
	APIKey    string
	ZoneID    string
	Email     string
}

var (
	domainConfigs map[string]CloudflareConfig
	mu            sync.RWMutex
)

// These variables will be set at build time
var (
	allowedDomains string
	configJSON     string
)

func init() {
	initializeConfigs()
}

func initializeConfigs() {
	mu.Lock()
	defer mu.Unlock()

	if domainConfigs != nil {
		return
	}

	domainConfigs = make(map[string]CloudflareConfig)

	// Parse the JSON configuration
	var configs map[string]CloudflareConfig
	err := json.Unmarshal([]byte(configJSON), &configs)
	if err != nil {
		panic("Failed to parse config JSON: " + err.Error())
	}

	// Initialize the domain configs
	for _, domain := range strings.Split(allowedDomains, ",") {
		domain = strings.TrimSpace(domain)
		if domain == "" {
			continue
		}
		if config, ok := configs[domain]; ok {
			domainConfigs[domain] = config
		}
	}
}

func SelectDomain(domain string) (CloudflareConfig, bool) {
	mu.RLock()
	defer mu.RUnlock()

	// First, try to match the exact domain
	if config, ok := domainConfigs[domain]; ok {
		return config, true
	}

	// If not found, try to match a subdomain
	baseDomain := getBaseDomain(domain)
	if config, ok := domainConfigs[baseDomain]; ok {
		return config, true
	}

	return CloudflareConfig{}, false
}

func getBaseDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return domain
}
