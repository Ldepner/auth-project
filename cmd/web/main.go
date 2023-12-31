package main

import (
	"github.com/Ldepner/auth-project/internal/authenticator"
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/drivers"
	"github.com/Ldepner/auth-project/internal/handlers"
	"github.com/Ldepner/auth-project/internal/helpers"
	"log"
	"net/http"
)

const PORT = ":8080"

var app config.AppConfig

func main() {
	// connect to DB
	db := drivers.DBConnect()

	// instantiate new Repo and handlers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	authenticator.NewAuthenticator(&repo.DB)

	log.Println("starting app on port", PORT, "...")

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}
