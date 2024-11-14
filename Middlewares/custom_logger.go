package middlewares

import (
	"log"
	"net/http"
	"time"
)

// CustomResponseWriter wraps http.ResponseWriter to capture the status code
type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it
func (crw *CustomResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// getStatusColor returns the appropriate color based on status code
func getStatusColor(code int) string {
	switch {
	case code >= 500:
		return "\033[31m" // Red
	case code >= 400:
		return "\033[35m" // Purple
	case code >= 200 && code < 300:
		return "\033[32m" // Green
	default:
		return "\033[37m" // White
	}
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create custom response writer to capture status code
		crw := &CustomResponseWriter{
			ResponseWriter: w,
			statusCode:     200, // Default to 200 if WriteHeader is never called
		}

		// Create a channel to calculate the response time asynchronously
		done := make(chan struct{})
		go func() {
			defer close(done)
			next.ServeHTTP(crw, r)
		}()
		<-done

		duration := time.Since(start)

		// ANSI color codes for terminal
		statusColor := getStatusColor(crw.statusCode)
		yellow := "\033[33m"
		blue := "\033[34m"
		reset := "\033[0m"

		log.Printf("%s%s%s %s%s%s Path: %s%s%s Status: %s%d%s Duration: %s%v%s\n",
			statusColor, r.Method, reset,
			yellow, r.Proto, reset,
			blue, r.URL.Path, reset,
			statusColor, crw.statusCode, reset,
			yellow, duration, reset)
	})
}
