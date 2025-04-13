package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type PaymentRequiredResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Plan    string `json:"plan"`
	Limit   int    `json:"limit"`
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

func SendPaymentRequiredResponse(w http.ResponseWriter, message string, plan string, limit int) {
	response := PaymentRequiredResponse{
		Status:  "fail",
		Message: message,
		Plan:    plan,
		Limit:   limit,
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPaymentRequired)
	w.Write(buf.Bytes())
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func CreateSendResponse(w http.ResponseWriter, data any, message string, statusCode int, dataName string, username string) {
	response := SuccessResponse{
		Success: true,
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

func SendSuccessResponse(w http.ResponseWriter, message string) {
	response := map[string]string{
		"message": message,
		"status":  "success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
