package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func CORS(allowedDomains []string) func(http.Handler) http.Handler {
	allowedHeaders := []string{
		authHeader,
		"content-type",
		"Authorization",
	}
	allowedMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodOptions,
		http.MethodDelete,
	}

	return func(next http.Handler) http.Handler {
		return handlers.CORS(
			handlers.AllowedOrigins(allowedDomains),
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods(allowedMethods),
		)(next)
	}

}
