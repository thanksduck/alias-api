package routes

import (
	"net/http"

	premium "github.com/thanksduck/alias-api/Controllers/Premium"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
)

func PremiumRouter(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/v2/premium/init", middlewares.Protect(premium.CreatePayment))

}
