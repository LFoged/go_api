package common

import (
	"net/http"
	"time"

	"api/cmd/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const maxBytes = 512000

func Routes() http.Handler {
	router := chi.NewRouter()

	// middleware - sec. headers, cors, reqID, maxBodyBytes, logger, recoverer, timeout
	router.Use(
		securityHeadersMiddleware,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: false,
			MaxAge:           300, // 5min
		}),
		middleware.RequestID,
		middleware.Logger,
		middleware.RequestSize(maxBytes),
		middleware.Timeout(25*time.Second),
		middleware.Recoverer,
	)

	// versioned routes
	v1Router := chi.NewRouter()
	v1Router.Get("/ping", handlers.Ping)
	v1Router.Get("/err", handlers.Err)
	v1Router.Post("/post", handlers.Post)

	// mount versioned routes onto main router
	router.Mount("/v1", v1Router)

	return router
}
