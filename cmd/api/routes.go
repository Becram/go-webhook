package main

import (
	"net/http"

	"github.com/Becram/go-webhook/internal/config"
	"github.com/Becram/go-webhook/internal/handlers"
	"github.com/go-chi/chi"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsoleNext)
	// mux.Use(NoSurf)
	// mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Post("/github/webhook", handlers.Repo.GHWebhook)
	return mux
}
