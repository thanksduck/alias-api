package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var requestData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if requestData.Password == "" {
		utils.SendErrorResponse(w, "Password is required", http.StatusBadRequest)
		return
	}

	username := strings.ToLower(requestData.Username)
	email := requestData.Email

	// Validate username and email
	if email == "" && username == "" {
		utils.SendErrorResponse(w, "Username or Email is required", http.StatusBadRequest)
		return
	}

	if email != "" && !middlewares.ValidBody.IsValidEmail(email) {
		utils.SendErrorResponse(w, "Email can't be processed", http.StatusUnprocessableEntity)
		return
	}

	if username != "" && !middlewares.ValidBody.IsValidUsername(username) {
		utils.SendErrorResponse(w, "Username can't be processed", http.StatusUnprocessableEntity)
		return
	}

	// Find user by username or email
	user, err := repository.FindUserByUsernameOrEmail(username, email)
	if err != nil {
		utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	if !utils.CheckPassword(requestData.Password, user.Password) {
		utils.SendErrorResponse(w, "Invalid username or password. If you used a social sign-in, please use that method.", http.StatusUnauthorized)
		return
	}

	utils.CreateSendResponse(w, user, "Login Successful", http.StatusOK, "user", user.ID)

}
