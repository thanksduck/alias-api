package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	// extract the token from the url
	var requestData struct {
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}
	// extract the token from the uri
	hashedToken := r.PathValue("token")
	fmt.Println(hashedToken)
	if hashedToken == "" {
		utils.SendErrorResponse(w, "Invalid token", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
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
		utils.SendErrorResponse(w, "Token is not valid", http.StatusBadRequest)
		return
	}
	userID := user.ID

	if user.PasswordResetExpires.Before(time.Now()) {
		err := repository.DeletePasswordResetToken(userID)
		if err != nil {
			utils.SendErrorResponse(w, "Something Went Wrong", http.StatusInternalServerError)
			return
		}
		utils.SendErrorResponse(w, "Your Password Reset Token Has been Expired", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := utils.HashPassword(requestData.Password)
	if err != nil {
		utils.SendErrorResponse(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	err = repository.UpdatePassword(userID, hashedPassword)
	if err != nil {
		utils.SendErrorResponse(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, nil, "Password Reset Successful Please Login", http.StatusOK, "message", 0)

}
