package infrastructure

import (
	"beer-api/infrastructure/database"
	"github.com/go-chi/chi"
	"net/http"
)

// Routes returns the API V1 Handler with configuration.
func Routes(conn *database.Data) http.Handler {
	router := chi.NewRouter()

	return router
}

// routesUser returns beer router with each endpoint.
func routesBeer() http.Handler {
	router := chi.NewRouter()

	return router
}
