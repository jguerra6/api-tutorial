package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	roleUser       = "user"
	roleAdmin      = "admin"
	RoleUser  Role = roleUser
	RoleAdmin Role = roleAdmin
)

type User struct {
	ID          *uuid.UUID
	Email       string
	Role        Role
	PhoneNumber string
	Password    string
	DisplayName string
	IsActive    bool
	IsVerified  bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func ParseRole(role string) Role {
	switch role {
	case roleUser:
		return RoleUser
	case roleAdmin:
		return RoleAdmin
	default:
		return roleUser
	}

}
