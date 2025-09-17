package users

import (
	"time"

	"github.com/jguerra6/api-tutorial/internal/ports"
)

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

type Service struct {
	Auth  ports.Auth
	Users ports.UserRepository
	Clock Clock
}

func NewService(auth ports.Auth, repo ports.UserRepository, clock Clock) *Service {
	if clock == nil {
		clock = realClock{}
	}
	return &Service{
		Auth:  auth,
		Users: repo,
		Clock: clock,
	}
}
