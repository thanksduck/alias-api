package routes

import (
	"net/http"

	destinations "github.com/thanksduck/alias-api/Controllers/Destinations"
	rules "github.com/thanksduck/alias-api/Controllers/Rules"
	user "github.com/thanksduck/alias-api/Controllers/User"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
)

func UserRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v2/user", middlewares.Protect(user.GetUser))
	mux.HandleFunc("PATCH /api/v2/user/{username}", middlewares.Protect(user.UpdateUser))
	mux.HandleFunc("GET /api/v2/user/{username}/destination", middlewares.Protect(destinations.ListDestinations))
	mux.HandleFunc("GET /api/v2/user/{username}/destinations", middlewares.Protect(destinations.ListDestinations))
	mux.HandleFunc("GET /api/v2/user/{username}/rule", middlewares.Protect(rules.ListRules))
	mux.HandleFunc("GET /api/v2/user/{username}/rules", middlewares.Protect(rules.ListRules))
	mux.HandleFunc("PATCH /api/v2/user/{username}/update-password", middlewares.Protect(user.UpdateUserPassword))
	mux.HandleFunc("DELETE /api/v2/user/{username}", middlewares.Protect(user.DeleteUser))
}
