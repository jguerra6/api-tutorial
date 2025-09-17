package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/jguerra6/api-tutorial/internal/domain"
)

type CreateUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role,omitempty"`
}

type UserResponse struct {
	ID          *uuid.UUID `json:"id"`
	Email       string     `json:"email"`
	DisplayName string     `json:"displayName"`
	PhoneNumber string     `json:"phone_number"`
	Role        string     `json:"role"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	JWT         string     `json:"jwt,omitempty"`
}

func MapReqToDomain(id *uuid.UUID, req *CreateUserRequest, now *time.Time) *domain.User {
	if req == nil {
		return &domain.User{}
	}

	return &domain.User{
		ID:          id,
		Email:       req.Email,
		Role:        domain.ParseRole(req.Role),
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		IsActive:    true,
		IsVerified:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

}

func MapResp(u *domain.User, token string) *UserResponse {

	return &UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		PhoneNumber: u.PhoneNumber,
		Role:        string(u.Role),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		JWT:         token,
	}
}
