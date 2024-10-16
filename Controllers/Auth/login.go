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

	if email == "" {
		if !middlewares.ValidBody.IsValidUsername(username) {
			utils.SendErrorResponse(w, "Username cant be processed", http.StatusUnprocessableEntity)
			return
		}
		user, err := repository.FindUserByUsernameOrEmail(username, "")
		if err != nil {
			utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Provider == "google" {
			utils.SendErrorResponse(w, "Please login with Google", http.StatusUnauthorized)
			return
		}

		if user.Provider == "github" {
			utils.SendErrorResponse(w, "Please login with Github", http.StatusUnauthorized)
			return
		}

		// check password
		if !utils.CheckPassword(requestData.Password, user.Password) {
			utils.SendErrorResponse(w, "Username or Password is invalid", http.StatusUnauthorized)
			return
		}

		utils.CreateSendResponse(w, user, "Login Successful", http.StatusOK, "user", user.ID)

	} else if username == "" {

		if !middlewares.ValidBody.IsValidEmail(requestData.Email) {
			utils.SendErrorResponse(w, "Email cant be processed", http.StatusUnprocessableEntity)
			return
		}
		user, err := repository.FindUserByUsernameOrEmail("", email)
		if err != nil {
			utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}

		// check password
		if !utils.CheckPassword(requestData.Password, user.Password) {
			utils.SendErrorResponse(w, "Email or Password is invalid", http.StatusUnauthorized)
			return
		}

		utils.CreateSendResponse(w, user, "Login Successful", http.StatusOK, "user", user.ID)
	} else if username != "" && email != "" {
		if !middlewares.ValidBody.IsValidUsername(username) {
			utils.SendErrorResponse(w, "Username cant be processed", http.StatusUnprocessableEntity)
			return
		}

		if !middlewares.ValidBody.IsValidEmail(email) {
			utils.SendErrorResponse(w, "Email cant be processed", http.StatusUnprocessableEntity)
			return
		}
		user, err := repository.FindUserByUsernameOrEmail(username, email)
		if err != nil {
			utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}

		if !utils.CheckPassword(requestData.Password, user.Password) {
			utils.SendErrorResponse(w, "Email or Password is invalid", http.StatusUnauthorized)
			return
		}

		utils.CreateSendResponse(w, user, "Login Successful", http.StatusOK, "user", user.ID)
	} else {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

}
