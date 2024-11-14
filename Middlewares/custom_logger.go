package middlewares

import (
	"fmt"
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

// getPathColor returns a unique color based on the request path
func getPathColor(path string) string {
	switch path {
	case "/api/v2/mail/rules":
		return "\033[36m" // Cyan
	case "/api/v2/user":
		return "\033[35m" // Magenta
	case "/api/v2/auth":
		return "\033[33m" // Yellow
	case "/api/v2/mail/destinations":
		return "\033[34m" // Blue
	default:
		return "\033[37m" // White for other paths
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
		pathColor := getPathColor(r.URL.Path)
		yellow := "\033[33m"
		reset := "\033[0m"

		fmt.Printf("%s%s%s %s%s%s Path: %s%s%s Status: %s%d%s Duration: %s%v%s\n",
			statusColor, r.Method, reset,
			yellow, r.Proto, reset,
			pathColor, r.URL.Path, reset,
			statusColor, crw.statusCode, reset,
			yellow, duration, reset)
	})
}
