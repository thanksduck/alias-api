package app

import (
	"fmt"
	"net/http"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	routes "github.com/thanksduck/alias-api/Routes"
)

func Init() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(`GET /api/v2/health`, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `I'm Healthy Boi`)
	})
	routes.AuthRouter(mux)
	routes.UserRouter(mux)
	routes.RuleRouter(mux)
	routes.DestinationRouter(mux)
	routes.PremiumRouter(mux)

	return middlewares.RequestLoggerMiddleware(mux)
}
