package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
)

func Logging(base *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &responseRecorder{ResponseWriter: w, status: 0, bytes: 0}

			reqLogger := base.With().
				Str("request_id", ctxutils.RequestID(r.Context())).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Logger()
			ctx := reqLogger.WithContext(r.Context())
			r = r.WithContext(ctx)

			next.ServeHTTP(rec, r)

			lat := time.Since(start)
			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}

			routeName, routeTpl := currentRouteInfo(r)
			if routeName == "healthz" || routeName == "readyz" || r.Method == http.MethodOptions {
				return
			}

			logger := ctxutils.Logger(r.Context())
			evt := logger.Info()

			switch {
			case status >= 500:
				evt = logger.Error()
			case status >= 400:
				evt = logger.Warn()
			}

			evt.
				Str("route", routeName).
				Str("route_tpl", routeTpl).
				Int("status", status).
				Int64("duration_ms", lat.Milliseconds()).
				Int("bytes", rec.bytes).
				Str("request_id", ctxutils.RequestID(r.Context())).
				Str("remote_ip", remoteIP(r)).
				Str("user_agent", safeUserAgent(r.UserAgent())).
				Msg("request finished")
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (rw *responseRecorder) WriteHeader(code int) {
	if rw.status == 0 {
		rw.status = code
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseRecorder) Write(p []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(p)
	rw.bytes += n
	return n, err
}

func currentRouteInfo(r *http.Request) (name, tpl string) {
	cr := mux.CurrentRoute(r)
	if cr == nil {
		return "", ""
	}
	if n := cr.GetName(); n != "" {
		name = n
	}
	if t, err := cr.GetPathTemplate(); err == nil {
		tpl = t
	}
	return
}

func remoteIP(r *http.Request) string {
	if xff := r.Header.Get(forwardedForHeader); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func safeUserAgent(ua string) string {

	if len(ua) > uaLen {
		return ua[:uaLen] + "…"
	}
	return ua
}
