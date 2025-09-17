package writer

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
	"github.com/jguerra6/api-tutorial/internal/ports"
)

func WriteAppError(w http.ResponseWriter, r *http.Request, err error) {
	status, kind := pickStatus(err)

	JSON(w, status, map[string]any{
		"error":   kind,
		"message": safeMessage(err),
	})

	rid := requestID(r)
	method := r.Method
	path := r.URL.Path

	routeName := ""
	if cr := mux.CurrentRoute(r); cr != nil {
		if name := cr.GetName(); name != "" {
			routeName = name
		} else if tpl, _ := cr.GetPathTemplate(); tpl != "" {
			routeName = tpl
		}
	}

	logger := ctxutils.Logger(r.Context())
	logger.Error().
		Str("kind", kind).
		Int("status", status).
		Str("method", method).
		Str("path", path).
		Str("route", routeName).
		Str("request_id", rid).
		Err(err).
		Msg("request error")
}

func pickStatus(err error) (int, string) {
	switch {
	case ports.IsCode(err, ports.CodeValidation):
		return http.StatusBadRequest, "validation_error"
	case ports.IsCode(err, ports.CodeUnauthenticated):
		return http.StatusUnauthorized, "unauthenticated"
	case ports.IsCode(err, ports.CodeUnauthorized):
		return http.StatusForbidden, "forbidden"
	case ports.IsCode(err, ports.CodeNotFound):
		return http.StatusNotFound, "not_found"
	case ports.IsCode(err, ports.CodeConflict):
		return http.StatusConflict, "conflict"
	case ports.IsCode(err, ports.CodeTooManyRequests):
		return http.StatusTooManyRequests, "rate_limited"
	case ports.IsCode(err, ports.CodeExternal):
		return http.StatusBadGateway, "external_error"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}

func safeMessage(err error) string { return err.Error() }
