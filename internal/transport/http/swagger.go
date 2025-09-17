package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/jguerra6/api-tutorial/internal/transport/http/docs"
)

func InitSwaggerBindings(r *mux.Router) {
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	r.PathPrefix("/docs").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.URL("/docs/doc.json"),
	))

}
