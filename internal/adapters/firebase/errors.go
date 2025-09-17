package firebase

import (
	"firebase.google.com/go/v4/errorutils"

	"github.com/jguerra6/api-tutorial/internal/ports"
)

func ParseFirebaseAuthError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errorutils.IsAlreadyExists(err):
		return ports.NewAppError(ports.CodeConflict, err.Error())
	case errorutils.IsNotFound(err):
		return ports.NewNotFoundError("user not found")
	default:
		return ports.NewAppError(ports.CodeInternal, "failed to create user")
	}

}
