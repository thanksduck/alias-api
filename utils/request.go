package utils

import (
	"io"
	"net/http"
)

// ReadRequestBody reads the request body and returns it as a byte slice.
func ReadRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}
