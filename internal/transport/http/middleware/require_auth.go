package middleware

import (
	"net/http"
	"strings"

	"github.com/jguerra6/api-tutorial/config"
	"github.com/jguerra6/api-tutorial/internal/domain"
	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
	"github.com/jguerra6/api-tutorial/internal/ports"
	"github.com/jguerra6/api-tutorial/internal/transport/http/writer"
)

func validateStaticToken(static string, adminTokens []string) bool {

	if static == "" {
		return false
	}

	for _, t := range adminTokens {
		if static == t {
			return true
		}
	}

	return false
}

func Auth(auth ports.Auth, cfg *config.Config) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			static := r.Header.Get(authHeader)
			if validateStaticToken(static, cfg.AdminTokens) {
				r = r.WithContext(ctxutils.WithUser(r.Context(), staticAdmin, string(domain.RoleAdmin)))
				next.ServeHTTP(w, r)
				return
			}

			h := r.Header.Get("Authorization")
			if len(h) < 8 || !strings.EqualFold(h[:7], "Bearer ") {
				writer.WriteAppError(w, r, ports.NewUnauthenticatedError("missing bearer token"))
				return
			}

			token := strings.TrimSpace(h[7:])

			uid, role, err := auth.VerifyIDToken(r.Context(), token)
			if err != nil {
				writer.WriteAppError(w, r, ports.NewUnauthenticatedError("invalid token"))
				return
			}

			r = r.WithContext(ctxutils.WithUser(r.Context(), uid, role))
			next.ServeHTTP(w, r)
		})
	}
}
