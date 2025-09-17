package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
)

func ExtractParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		params := make(map[string]string)
		for k, v := range vars {
			params[k] = v
		}
		query := r.URL.Query()
		for k, v := range query {
			params[k] = v[0]
		}

		ctx := ctxutils.WithParams(r.Context(), params)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
