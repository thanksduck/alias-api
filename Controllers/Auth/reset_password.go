package auth

import (
	"encoding/json"
	"net/http"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}

	hashedToken := r.PathValue("token")
	if hashedToken == "" {
		utils.SendErrorResponse(w, "Invalid reset link", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
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

	user, err := repository.FindUserByPasswordResetToken(hashedToken)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid reset link", http.StatusBadRequest)
		return
	}
	userID := user.ID

	err = repository.IsPasswordResetTokenExpired(hashedToken)
	if err != nil {
		err = repository.RemovePasswordResetToken(userID)
		if err != nil {
			utils.SendErrorResponse(w, "An error occurred. Please try again", http.StatusInternalServerError)
			return
		}
		utils.SendErrorResponse(w, "Reset link has expired. Please request a new one", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(requestData.Password)
	if err != nil {
		utils.SendErrorResponse(w, "An error occurred. Please try again", http.StatusInternalServerError)
		return
	}

	err = repository.UpdatePassword(userID, hashedPassword)
	if err != nil {
		utils.SendErrorResponse(w, "An error occurred. Please try again", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, nil, "Password reset successful. You can now login with your new password", http.StatusOK, "message", `0`)
}
