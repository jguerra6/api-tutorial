package postgres

import (
	"fmt"

	"github.com/jguerra6/api-tutorial/internal/domain"
)

func userFromDomain(u *domain.User) (*UserRow, error) {
	if u == nil {
		return nil, fmt.Errorf("user is nil")
	}

	return &UserRow{
		ID:          u.ID,
		Email:       u.Email,
		Role:        string(u.Role),
		PhoneNumber: u.PhoneNumber,
		DisplayName: u.DisplayName,
		IsActive:    u.IsActive,
		IsVerified:  u.IsVerified,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}, nil
}
