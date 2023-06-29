package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// The function returns an HTTP handler that loads and saves session data.
// The NoSurf function sets a base cookie with certain properties and returns a handler that uses the
// nosurf package to prevent cross-site request forgery attacks.
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
