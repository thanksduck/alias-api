package middlewares

import (
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	usernameRegex  = regexp.MustCompile(`^[a-z][a-z0-9-_\.]{3,15}$`)
	nameRegex      = regexp.MustCompile(`^[a-zA-Z\s]{3,64}$`)
	emailRegex     = regexp.MustCompile(`^[a-zA-Z0-9.-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex  = regexp.MustCompile(`^.{8,20}$`)
	subDomainRegex = regexp.MustCompile(`^[a-z0-9]{1,8}$`)
)

var (
	allowedDomains map[string]struct{}
	once           sync.Once
)

func initAllowedDomains() {
	once.Do(func() {
		allowedDomains = make(map[string]struct{})
		domains := os.Getenv("ALLOWED_DOMAINS")
		for _, domain := range strings.Split(domains, ",") {
			allowedDomains[strings.TrimSpace(domain)] = struct{}{}
		}
	})
}

type validBody struct{}

var ValidBody validBody

func (v validBody) IsValidUsername(name string) bool {
	return usernameRegex.MatchString(name)
}

func (v validBody) IsValidName(name string) bool {
	return nameRegex.MatchString(name)
}

func (v validBody) IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (v validBody) IsValidPassword(password string) bool {
	return passwordRegex.MatchString(password)
}

func (v validBody) IsValidDomain(domain string) bool {
	return subDomainRegex.MatchString(domain)
}

func (v validBody) IsAllowedDomain(domain string) bool {
	initAllowedDomains()
	_, exists := allowedDomains[domain]
	return exists
}
