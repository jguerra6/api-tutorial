package middleware

import (
	"errors"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
	"github.com/jguerra6/api-tutorial/internal/transport/http/writer"
)

func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rec := &headerRecorder{ResponseWriter: w}

			defer func() {
				if p := recover(); p != nil {
					if isBrokenPipe(p) {
						logPanic(p, r, 0, true)
						return
					}

					if !rec.wroteHeader {
						writer.JSON(rec, http.StatusInternalServerError, map[string]any{
							"error":   "internal_error",
							"message": "unexpected error",
						})
					}

					logPanic(p, r, http.StatusInternalServerError, false)
				}
			}()

			next.ServeHTTP(rec, r)
		})
	}
}

type headerRecorder struct {
	http.ResponseWriter
	wroteHeader bool
}

func (hr *headerRecorder) WriteHeader(code int) {
	if !hr.wroteHeader {
		hr.wroteHeader = true
	}
	hr.ResponseWriter.WriteHeader(code)
}

func logPanic(p any, r *http.Request, status int, clientGone bool) {
	routeName, routeTpl := currentRouteInfo(r)
	rid := r.Header.Get("X-Request-ID")

	kind := "panic"
	if clientGone {
		kind = "client_gone"
	}

	logger := ctxutils.Logger(r.Context())
	logger.Error().
		Str("kind", kind).
		Int("status", status).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("route", routeName).
		Str("route_tpl", routeTpl).
		Str("request_id", rid).
		Interface("panic", p).
		Str("stack", string(debug.Stack())).
		Msg("panic recovered")
}

func isBrokenPipe(p any) bool {
	err, ok := p.(error)
	if !ok {
		if s, ok := p.(string); ok {
			ls := strings.ToLower(s)
			return strings.Contains(ls, "broken pipe") || strings.Contains(ls, "connection reset by peer")
		}
		return false
	}

	var ne *net.OpError
	if errors.As(err, &ne) {
		var se *os.SyscallError
		if errors.As(ne, &se) {
			return errors.Is(se.Err, syscall.EPIPE) || errors.Is(se.Err, syscall.ECONNRESET)
		}
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer")
}
