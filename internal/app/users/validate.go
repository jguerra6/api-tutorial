package users

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/jguerra6/api-tutorial/internal/ports"
	"github.com/jguerra6/api-tutorial/internal/transport/http/dto"
)

var (
	emailRx        = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	phoneRx        = regexp.MustCompile(`^\+?[0-9\-.\s()]{7,20}$`)
	minPasswordLen = 8
)

func validateCreate(cmd *dto.CreateUserRequest) error {
	if cmd.Email == "" || !emailRx.MatchString(cmd.Email) {
		return ports.NewValidationError("invalid email")
	}

	if probs := validatePassword(cmd.Password); probs != "" {
		return ports.NewValidationError(probs)
	}

	if cmd.DisplayName == "" {
		return ports.NewValidationError("displayName is required")
	}

	if cmd.PhoneNumber != "" && !phoneRx.MatchString(cmd.PhoneNumber) {
		return ports.NewValidationError("invalid phone")
	}

	return nil
}

func validatePassword(password string) string {
	var probs []string
	if len(password) < minPasswordLen {
		probs = append(probs, "must be at least 8 characters")
	}

	var up, lo, di, sy bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			up = true
		case unicode.IsLower(r):
			lo = true
		case unicode.IsDigit(r):
			di = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			sy = true
		}
	}
	if !up {
		probs = append(probs, "must include an uppercase letter")
	}
	if !lo {
		probs = append(probs, "must include a lowercase letter")
	}
	if !di {
		probs = append(probs, "must include a number")
	}
	if !sy {
		probs = append(probs, "must include a symbol")
	}

	if len(probs) == 0 {
		return ""
	}

	return "Password requirements: " + strings.Join(probs, ", ")
}
