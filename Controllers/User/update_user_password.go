package user

import (
	"encoding/json"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	"github.com/thanksduck/alias-api/utils"
)

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
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

	if requestData.Password == "" || requestData.PasswordConfirm == "" || requestData.CurrentPassword == "" {
		utils.SendErrorResponse(w, "All fields are required", http.StatusBadRequest)
		return
	}

	savedPassword, _ := db.SQL.FindPasswordById(ctx, user.ID)
	if !utils.CheckPassword(requestData.CurrentPassword, savedPassword) {
		utils.SendErrorResponse(w, "Current password is incorrect", http.StatusBadRequest)
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

	err = db.SQL.UpdatePasswordUser(ctx, &q.UpdatePasswordUserParams{ID: user.ID, Password: hashedPassword})
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, nil, "Password updated successfully", http.StatusOK, "user", user.Username)

}
