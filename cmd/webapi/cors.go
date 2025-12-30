package main

import (
	"github.com/gorilla/handlers"
	"net/http"
)

func applyCORSHandler(h http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type",
		}),
		handlers.AllowedHeaders([]string{"x-example-header","Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.MaxAge(1),
	)(h)
}
