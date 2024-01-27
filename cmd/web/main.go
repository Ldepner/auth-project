package main

import (
	"context"
	"github.com/Ldepner/auth-project/internal/authenticator"
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/drivers"
	"github.com/Ldepner/auth-project/internal/handlers"
	"github.com/Ldepner/auth-project/internal/helpers"
	session_manager "github.com/Ldepner/auth-project/internal/session_manager"
	"github.com/go-webauthn/webauthn/webauthn"
	"log"
	"net/http"
)

const PORT = ":8080"

var app config.AppConfig

func main() {
	// connect to DB
	db := drivers.DBConnect()
	defer func() {
		if err := db.Client.Disconnect(context.TODO()); err != nil {
			log.Println("error connecting to DB")
			panic(err)
		}
	}()

	// initialize webAuthN
	wConfig := &webauthn.Config{
		RPDisplayName: "Go Webauthn",                     // Display Name for your site
		RPID:          "localhost",                       // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:8080"}, // The origin URLs allowed for WebAuthn requests
	}

	var err error
	if app.WebAuthn, err = webauthn.New(wConfig); err != nil {
		log.Fatal(err)
	}

	// instantiate new Repo and handlers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	authenticator.NewAuthenticator(&repo.DB)
	session_manager.NewSessionManager(repo.DB)

	log.Println("starting app on port", PORT, "...")

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}
