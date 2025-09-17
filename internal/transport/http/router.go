package transporthttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"github.com/jguerra6/api-tutorial/config"
	"github.com/jguerra6/api-tutorial/internal/app/users"
	"github.com/jguerra6/api-tutorial/internal/ports"
	"github.com/jguerra6/api-tutorial/internal/transport/http/handlers"
	"github.com/jguerra6/api-tutorial/internal/transport/http/middleware"
)

const (
	apiVersion = 1
	timeOut    = 30 * time.Second
)

func NewRouter(cfg *config.Config, userSvc users.Service, dbPing handlers.DBPinger, auth ports.Auth, logger *zerolog.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(
		middleware.CORS(cfg.AllowedOrigins),
		middleware.RequestID,
		middleware.ExtractParams,
		middleware.Logging(logger),
		middleware.Recovery(),
		middleware.Timeout(timeOut),
	)

	hh := handlers.NewHealthHandler(dbPing)
	r.HandleFunc("/healthz", hh.Healthz).Methods(http.MethodGet).Name("healthz")
	r.HandleFunc("/readyz", hh.Readyz).Methods(http.MethodGet).Name("readyz")
	InitSwaggerBindings(r)

	uh := handlers.NewUserHandler(&userSvc)

	version := fmt.Sprintf("/v%d", apiVersion)
	v1Routes := r.PathPrefix(version).Subrouter()
	publicRoutes(v1Routes, uh)

	adminRouter := v1Routes.PathPrefix("/admin").Subrouter()
	adminRouter.Use(
		middleware.Auth(auth, cfg),
	)

	initGet(adminRouter, uh)
	initPost(adminRouter, uh)
	initDelete(adminRouter, uh)

	return r
}

func publicRoutes(r *mux.Router, uh *handlers.UsersHandler) {
	var (
		//getRouter = r.Methods(http.MethodGet).Subrouter()
		postRouter = r.Methods(http.MethodPost).Subrouter()
	)

	postRouter.HandleFunc("/users", uh.PublicCreateUser).
		Name("users.public_create")
}

func initGet(r *mux.Router, uh *handlers.UsersHandler) {
	var (
	//router = r.Methods(http.MethodGet).Subrouter()
	)
}

func initPost(r *mux.Router, uh *handlers.UsersHandler) {
	var (
		router = r.Methods(http.MethodPost).Subrouter()
	)

	router.HandleFunc("/users", uh.AdminCreateUser).
		Name("users.admin_create")
}

func initDelete(r *mux.Router, uh *handlers.UsersHandler) {
	var (
		router = r.Methods(http.MethodDelete).Subrouter()
	)

	router.HandleFunc("/users/{userId}", uh.DeleteUser).
		Name("users.admin_delete")
}
