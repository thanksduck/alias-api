package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	routes "github.com/thanksduck/alias-api/Routes"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Init() http.Handler {
	fmt.Println("Making the Application with HTTP/2 support")

	mux := http.NewServeMux()
	mux.HandleFunc(`GET /health`, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `I'm Healthy Boi`)
	})
	routes.AuthRouter(mux)
	routes.UserRouter(mux)
	routes.RuleRouter(mux)
	routes.DestinationRouter(mux)
	// Create an HTTP/2 server
	h2s := &http2.Server{}

	// Wrap the mux with h2c for HTTP/2 support without TLS
	h2cHandler := h2c.NewHandler(RequestLoggerMiddleware(mux), h2s)

	return h2cHandler
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
