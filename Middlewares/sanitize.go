package middlewares

import (
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-z][a-z0-9-_\.]{3,15}$`)
	nameRegex     = regexp.MustCompile(`^[a-zA-Z\s]{3,64}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9.-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex = regexp.MustCompile(`^.{8,20}$`)
)

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
