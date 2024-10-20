package user

import (
	"encoding/json"
	"net/http"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	var requestData struct {
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
		CurrentPassword string `json:"currentPassword"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.CheckPassword(requestData.CurrentPassword, user.Password) {
		utils.SendErrorResponse(w, "Current password is incorrect", http.StatusBadRequest)
		return
	}

	if requestData.Password == "" || requestData.PasswordConfirm == "" || requestData.CurrentPassword == "" {
		utils.SendErrorResponse(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !middlewares.ValidBody.IsValidPassword(requestData.Password) {
		utils.SendErrorResponse(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	if requestData.Password != requestData.PasswordConfirm {
		utils.SendErrorResponse(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(requestData.Password)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword

	err = repository.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, user, "Password updated successfully", http.StatusOK, "user", user.Username)

}
