package routes

import (
	"net/http"

	domains "github.com/thanksduck/alias-api/Controllers/Domains"
	premium "github.com/thanksduck/alias-api/Controllers/Premium"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
)

func PremiumRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v2/premium/domain", middlewares.Protect(domains.CreateCustomDomain))
	mux.HandleFunc("POST /api/v2/premium/star", middlewares.Protect(premium.SubscribeToStar))

}
