package main

import (
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Post("/logout", handlers.Repo.PostLogout)
	mux.Get("/register", handlers.Repo.Register)
	mux.Post("/register", handlers.Repo.PostRegister)
	mux.Get("/registration/begin", handlers.Repo.BeginRegistration)
	mux.Post("/registration/finish", handlers.Repo.FinishRegistration)
	mux.Get("/login/begin", handlers.Repo.BeginLogin)
	mux.Post("/login/finish", handlers.Repo.FinishLogin)

	mux.Route("/", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/", handlers.Repo.Home)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
