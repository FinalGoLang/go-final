package main

import (
	"context"
	"films/config"
	"films/internal/service"
	"films/internal/storage/email"
	"films/internal/storage/psql"
	"films/internal/transport/rest/handler"
	"films/pkg/db"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// init config
	logger := hclog.Default()
	err := config.Init()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// connect to db
	dbConn, err := db.ConnectDB(viper.GetString("app.db_path"))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// migrate to db
	err = db.MigrateToDB(viper.GetString("app.db_path"), logger.Named("migration"))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// setup storages
	emailClient := email.NewEmailClient(viper.GetString("app.emailUsername"))
	userStorage := psql.NewAccount(dbConn)
	filmStorage := psql.NewFilm(dbConn)
	subStorage := psql.NewSubscribe(dbConn)

	// setup services
	userService := service.NewUserService(userStorage)
	authService := service.NewAuthService(
		emailClient,
		userStorage, viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetDuration("auth.token_ttl"))
	filmService := service.NewFilmService(filmStorage)
	subService := service.NewSubscribeService(subStorage)

	// setup handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	filmHandler := handler.NewFilmHandler(filmService)
	subHandler := handler.NewSubscribeHandler(subService)

	// init routes
	restHandler := handler.New(handler.Deps{
		AuthHandler: *authHandler,
		UserHandler: *userHandler,
		FilmHandler: *filmHandler,
		SubscribeHandler: *subHandler,
	}).InitRoutes(authService)

	// starting server
	httpServer := &http.Server{
		Addr:    ":" + viper.GetString("port"),
		Handler: restHandler,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Failed to shutdown: %+v", err)
	}
}
