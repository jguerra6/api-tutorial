package postgres

import (
	"time"

	"github.com/google/uuid"
)

type UserRow struct {
	ID          *uuid.UUID `db:"id"`
	Email       string     `db:"email"`
	Role        string     `db:"role"`
	PhoneNumber string     `db:"phone_number"`
	DisplayName string     `db:"display_name"`
	IsActive    bool       `db:"is_active"`
	IsVerified  bool       `db:"is_verified"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
