package users

import (
	"context"

	"github.com/jguerra6/api-tutorial/internal/transport/http/dto"
)

func (s *Service) GetUser(ctx context.Context, uid string) (*dto.UserResponse, error) {
	//s.Auth.GetUser()

	return nil, nil
}
