package user

import (
	"encoding/json"
	"net/http"
)

func LogoutUser(w http.ResponseWriter, r *http.Request) {

	// Clear the cookie
	cookie := &http.Cookie{
		Name:     "token", // Replace with your actual cookie name
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	response := map[string]string{"message": "Successfully logged out", "status": "success"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
