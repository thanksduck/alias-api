package user

import (
	"encoding/json"
	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"strings"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	"github.com/thanksduck/alias-api/utils"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	var requestData struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// at least one field is required
	if requestData.Username == "" && requestData.Name == "" && requestData.Email == "" && requestData.Avatar == "" {
		utils.SendErrorResponse(w, "At least one field is required", http.StatusBadRequest)
		return
	}
	if requestData.Username != "" {

		requestData.Username = strings.ToLower(requestData.Username)

		if !middlewares.ValidBody.IsValidUsername(requestData.Username) {
			utils.SendErrorResponse(w, "Username must be at least 4 characters long, start with an alphabet, and can contain numbers, dots, or underscores", http.StatusBadRequest)
			return
		}
		// check if the user already exists
		_, err = db.SQL.FindUserByUsername(ctx, requestData.Username)
		if err == nil {
			utils.SendErrorResponse(w, "Username Is Taken", http.StatusConflict)
			return
		}
		user.Username = requestData.Username
	}
	if requestData.Name != "" {
		if !middlewares.ValidBody.IsValidName(requestData.Name) {
			utils.SendErrorResponse(w, "Name must be more than 4 characters long and can only contain alphabets and spaces", http.StatusBadRequest)
			return
		}
		user.Name = requestData.Name
	}
	if requestData.Email != "" {
		if !middlewares.ValidBody.IsValidEmail(requestData.Email) {
			utils.SendErrorResponse(w, "Email must contain only letters, numbers, dots, or hyphens", http.StatusBadRequest)
			return
		}
		user.Email = requestData.Email
	}
	if requestData.Avatar != "" {
		user.Avatar = requestData.Avatar
	}
	err = db.SQL.UpdateUser(ctx, &q.UpdateUserParams{
		Username: user.Username,
		Name:     user.Name,
		Avatar:   user.Avatar,
		Email:    user.Email,
		ID:       user.ID,
	})
	if err != nil {
		utils.SendErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	s := &models.SafeUser{
		Username:         user.Username,
		Name:             user.Name,
		Email:            user.Email,
		IsPremium:        user.IsPremium,
		IsEmailVerified:  user.IsEmailVerified,
		AliasCount:       user.AliasCount,
		DestinationCount: user.DestinationCount,
		Avatar:           user.Avatar,
	}
	utils.CreateSendResponse(w, s, "User updated successfully", http.StatusOK, "user", user.Username)
}
