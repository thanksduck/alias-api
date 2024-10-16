package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	routes "github.com/thanksduck/alias-api/Routes"
)

func Init() http.Handler {
	fmt.Println("Making the Application")
	mux := http.NewServeMux()
	mux.HandleFunc(`GET /health`, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `I'm Healthy Boi`)
	})
	routes.AuthRouter(mux)
	routes.UserRouter(mux)
	routes.RuleRouter(mux)
	routes.DestinationRouter(mux)

	return RequestLoggerMiddleware(mux)
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Create a channel to calculate the response time asynchronously
		done := make(chan struct{})
		go func() {
			defer close(done)
			next.ServeHTTP(w, r)
		}()
		<-done
		duration := time.Since(start)
		log.Printf("Method: %s, Path: %s, Protocol: %s, Response time: %s",
			r.Method, r.URL.Path, r.Proto, duration)
	})
}
