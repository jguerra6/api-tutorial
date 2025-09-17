package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	fbApp "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"

	"github.com/jguerra6/api-tutorial/config"
	"github.com/jguerra6/api-tutorial/internal/adapters/firebase"
	"github.com/jguerra6/api-tutorial/internal/adapters/postgres"
	"github.com/jguerra6/api-tutorial/internal/app/users"
	"github.com/jguerra6/api-tutorial/internal/transport/http"
)

const (
	servicePrefix      = "service"
	serviceName        = "api-tutorial"
	appShutdownTimeout = 30 * time.Second
	readTimeout        = 30 * time.Second
	writeTimeout       = 30 * time.Second
	idleTimeout        = 120 * time.Second
	handlersTimeOut    = 30 * time.Second
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var (
		console = zerolog.NewConsoleWriter()
		ctx     = context.Background()
	)
	logger := zerolog.New(console).With().
		Str(servicePrefix, serviceName).
		Timestamp().
		Logger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	opt := option.WithCredentialsFile(cfg.FirebaseConfigFile)
	authApp, err := fbApp.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	authClient, err := authApp.Auth(ctx)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	authAdapter := firebase.NewAuthAdapter(authClient)
	db, err := postgres.NewDb(ctx, "", cfg.DbHost, cfg.DbUserName, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	userRepo := postgres.NewRepository(db)
	userSvc := users.NewService(authAdapter, userRepo, nil)

	router := transporthttp.NewRouter(cfg, *userSvc, nil, authAdapter, &logger)

	startServer(router, &logger, cfg.Port)

}

func startServer(restServer *mux.Router, logger *zerolog.Logger, port string) {
	logger.Info().Str("addr", port).Msg("starting server")
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        restServer,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Info().
		Str("server address", port).
		Msg("server up")

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatal().
				Err(err).
				Msg("HTTP server closed")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info().
		Msgf("received terminate %s, graceful shutdown", sig)

	shutdownServer(srv, logger)
}

func shutdownServer(restServer *http.Server, logger *zerolog.Logger) {

	ctx, cancel := context.WithTimeout(context.Background(), appShutdownTimeout)
	defer cancel()
	err := restServer.Shutdown(ctx)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("error shutting down the HTTP server")
	}
}
