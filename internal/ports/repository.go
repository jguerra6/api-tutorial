package ports

import (
	"context"

	"github.com/jguerra6/api-tutorial/internal/domain"
)

type UserRepository interface {
	Insert(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, uid string) error
}
