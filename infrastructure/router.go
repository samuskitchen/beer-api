package infrastructure

import (
	v1 "beer-api/domain/beer/application/v1"
	"beer-api/infrastructure/database"
	"github.com/go-chi/chi"
	"net/http"
)

// Routes returns the API V1 Handler with configuration.
func Routes(conn *database.Data) http.Handler {
	router := chi.NewRouter()

	br := v1.NewBeerHandler(conn)
	router.Mount("/beers", routesBeer(br))

	return router
}

// routesUser returns beer router with each endpoint.
func routesBeer(handler *v1.BeersRouter) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handler.GetAllBeersHandler)
	router.Get("/{beerID}", handler.GetOneHandler)
	router.Get("/{beerID}/boxprice", handler.GetOneBoxPriceHandler)
	router.Post("/", handler.CreateHandler)

	return router
}
