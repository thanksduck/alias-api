package routes

import (
	"net/http"

	rules "github.com/thanksduck/alias-api/Controllers/Rules"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
)

func RuleRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v2/rules", middlewares.Protect(rules.ListRules))
	mux.HandleFunc("GET /api/v2/rules/{id}", middlewares.Protect(rules.GetRule))
	mux.HandleFunc("POST /api/v2/rules", middlewares.Protect(rules.CreateRule))
	mux.HandleFunc("PATCH /api/v2/rules/{id}", middlewares.Protect(rules.UpdateRule))
	mux.HandleFunc("PATCH /api/v2/rules/{id}/toggle", middlewares.Protect(rules.ToggleRule))
	mux.HandleFunc("DELETE /api/v2/rules/{id}", middlewares.Protect(rules.DeleteRule))
}
