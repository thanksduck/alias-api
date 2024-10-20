package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Status:  "fail",
		Message: message,
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateSendResponse(w http.ResponseWriter, data interface{}, message string, statusCode int, dataName string, username string) {
	response := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := GenerateToken(username)
	if err != nil {
		SendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set headers and cookies
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(1 * time.Hour),
	})

	// Write status code and body
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
