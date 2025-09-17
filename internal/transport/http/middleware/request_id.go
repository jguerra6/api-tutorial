package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(requestIDHeader)
		if rid != "" {
			if _, err := uuid.Parse(rid); err != nil {
				rid = ""
			}
		}
		if rid == "" {
			rid = uuid.NewString()
		}

		w.Header().Set("X-Request-ID", rid)

		ctx := ctxutils.WithRequestID(r.Context(), rid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
