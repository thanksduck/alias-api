package routes

import (
	"net/http"

	domains "github.com/thanksduck/alias-api/Controllers/Domains"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
)

func PremiumRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v2/premium/domain", middlewares.Protect(domains.CreateCustomDomain))

}
