package main

import (
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/handlers"
	"log"
	"net/http"
)

const PORT = ":8080"

var app config.AppConfig

func main() {
	// instantiate new Repo and handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	log.Println("starting app on port", PORT, "...")

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}
