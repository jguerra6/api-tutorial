package ports

import "context"

type Auth interface {
	CreateAuthUser(ctx context.Context, id, email, password, displayName, phone, role string) (providerUID string, err error)
	DeleteUser(ctx context.Context, id string) error
	VerifyIDToken(ctx context.Context, token string) (userID string, role string, err error)
}
