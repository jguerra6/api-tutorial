package users

import (
	"context"

	"github.com/google/uuid"

	"github.com/jguerra6/api-tutorial/internal/ports"
	"github.com/jguerra6/api-tutorial/internal/transport/http/dto"
)

func (s *Service) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	var (
		id  = uuid.New()
		now = s.Clock.Now()
	)

	token, err := s.Auth.CreateAuthUser(ctx, id.String(), req.Email, req.Password, req.DisplayName, req.PhoneNumber, req.Role)
	if err != nil {
		return nil, ports.NewExternalError("failed to create auth user: " + err.Error())
	}

	u := dto.MapReqToDomain(&id, req, &now)

	if err = s.Users.Insert(ctx, u); err != nil {
		_ = s.Auth.DeleteUser(ctx, id.String())
		return nil, err
	}

	return dto.MapResp(u, token), nil
}
