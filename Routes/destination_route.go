package routes

import (
	"net/http"

	destinations "github.com/thanksduck/alias-api/Controllers/Destinations"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
)

func DestinationRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v2/destinations", middlewares.Protect(destinations.ListDestinations))
	mux.HandleFunc("GET /api/v2/destinations/{id}", middlewares.Protect(destinations.GetDestination))
	mux.HandleFunc("GET /api/v2/destinations/{id}/verify", middlewares.Protect(destinations.VerifyDestination))
	mux.HandleFunc("POST /api/v2/destinations", middlewares.Protect(destinations.CreateDestination))
	// mux.HandleFunc("PATCH /api/v2/destinations/{id}", middlewares.Protect(destinations.UpdateDestination))
	mux.HandleFunc("DELETE /api/v2/destinations/{id}", middlewares.Protect(destinations.DeleteDestination))
}
