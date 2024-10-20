package user

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
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
	// atleast one field is required
	if requestData.Username == "" && requestData.Name == "" && requestData.Email == "" && requestData.Avatar == "" {
		utils.SendErrorResponse(w, "Atleast one field is required", http.StatusBadRequest)
		return
	}
	if requestData.Username != "" {

		if !middlewares.ValidBody.IsValidUsername(requestData.Username) {
			utils.SendErrorResponse(w, "Username must be at least 4 characters long, start with an alphabet, and can contain numbers, dots, or underscores", http.StatusBadRequest)
			return
		}
		// check if the user already exists
		_, err = repository.FindUserByUsernameOrEmail(requestData.Username, ``)
		if err == nil {
			utils.SendErrorResponse(w, "Username Is Taken", http.StatusConflict)
			return
		} else if err != pgx.ErrNoRows {
			utils.SendErrorResponse(w, "Error checking user existence", http.StatusInternalServerError)
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
	updatedUser, err := repository.UpdateUser(user.ID, user)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, updatedUser, "User updated successfully", http.StatusOK, "user", updatedUser.Username)
}
