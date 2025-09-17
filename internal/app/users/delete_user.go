package users

import (
	"context"
)

func (s *Service) DeleteUser(ctx context.Context, uid string) error {

	if err := s.Users.Delete(ctx, uid); err != nil {
		return err
	}

	if err := s.Auth.DeleteUser(ctx, uid); err != nil {
		return err
	}

	return nil
}
