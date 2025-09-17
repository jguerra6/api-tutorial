package firebase

import (
	"context"
	"strings"
	"time"

	fbauth "firebase.google.com/go/v4/auth"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
	"github.com/jguerra6/api-tutorial/internal/ports"
)

var _ ports.Auth = &AuthAdapter{}

const (
	bearerPrefix = "bearer "
	roleKey      = "role"
	roleUser     = "user"
	roleAdmin    = "admin"
)

type AuthAdapter struct {
	Client *fbauth.Client
}

func (a *AuthAdapter) DeleteUser(ctx context.Context, id string) error {
	logger := ctxutils.Logger(ctx).With().
		Str("component", "authAdapter.user").
		Str("op", "delete").
		Str("user_id", id).
		Logger()

	start := time.Now()

	if err := a.Client.DeleteUser(ctx, id); err != nil {
		logger.Error().Dur("dur", time.Since(start)).Err(err).Msg("unable to delete user from Auth")
		return ParseFirebaseAuthError(err)
	}

	return nil
}

func NewAuthAdapter(client *fbauth.Client) *AuthAdapter {
	return &AuthAdapter{
		Client: client,
	}
}

func (a *AuthAdapter) CreateAuthUser(ctx context.Context, id, email, password, displayName, phone, role string) (string, error) {
	params := (&fbauth.UserToCreate{}).
		UID(id).
		Email(email).
		Password(password).
		DisplayName(displayName).
		PhoneNumber(phone)

	u, err := a.Client.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}

	if role == "" {
		role = roleUser
	}

	claims := createClaims(role)
	if err = a.Client.SetCustomUserClaims(ctx, u.UID, claims); err != nil {
		_ = a.Client.DeleteUser(ctx, u.UID)
		return "", err
	}

	token, err := a.Client.CustomTokenWithClaims(ctx, u.UID, claims)
	if err != nil {
		_ = a.Client.DeleteUser(ctx, u.UID)
		return "", err
	}

	return token, nil

}

func createClaims(role string) map[string]interface{} {
	return map[string]interface{}{
		roleKey: role,
	}
}

func (a *AuthAdapter) SetUserRole(ctx context.Context, uid, role string) error {

	claims := map[string]interface{}{
		roleKey: role,
	}
	return a.Client.SetCustomUserClaims(ctx, uid, claims)
}

func (a *AuthAdapter) VerifyIDToken(ctx context.Context, token string) (string, string, error) {

	if strings.HasPrefix(strings.ToLower(token), bearerPrefix) {
		token = strings.TrimSpace(token[7:])
	}

	tok, err := a.Client.VerifyIDToken(ctx, token)
	if err != nil {
		return "", "", err
	}

	var (
		uid  = tok.UID
		role = roleUser
	)

	if v, ok := tok.Claims[roleKey]; ok {
		if s, ok := v.(string); ok && s != "" {
			role = strings.ToLower(s)
		}
	}

	return uid, role, nil
}
