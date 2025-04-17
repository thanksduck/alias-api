package auth

import (
	"encoding/json"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"time"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	"github.com/thanksduck/alias-api/utils"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	userID, err := db.SQL.FindUserByValidResetToken(ctx, hashedToken)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid reset link", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(requestData.Password)
	if err != nil {
		utils.SendErrorResponse(w, "An error occurred. Please try again", http.StatusInternalServerError)
		return
	}
	err = db.SQL.UpdatePasswordUser(ctx, &q.UpdatePasswordUserParams{
		Password:          hashedPassword,
		PasswordChangedAt: time.Now(),
		ID:                userID,
	})
	if err != nil {
		utils.SendErrorResponse(w, "An error occurred. Please try again", http.StatusInternalServerError)
		return
	}
	err = db.SQL.UpdatePasswordAuth(ctx, userID)
	if err != nil {
		utils.SendErrorResponse(w, "Password updated but failed to clean reset token", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, nil, "Password reset successful. You can now login with your new password", http.StatusOK, "message", `0`)
}
