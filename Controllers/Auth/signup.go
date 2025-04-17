package auth

import (
	"encoding/json"
	models "github.com/thanksduck/alias-api/Models"
	"net/http"
	"strings"

	db "github.com/thanksduck/alias-api/Database"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
	q "github.com/thanksduck/alias-api/internal/db"
	"github.com/thanksduck/alias-api/utils"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestData struct {
		Username        string `json:"username"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
		Avatar          string `json:"avatar"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)

	if err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	username := strings.ToLower(requestData.Username)
	email := strings.ToLower(requestData.Email)
	if !middlewares.ValidBody.IsValidUsername(username) {
		utils.SendErrorResponse(w, "Username must be at least 4 characters long, start with an alphabet, and can contain numbers, dots, or underscores", http.StatusBadRequest)
		return
	}

	if !middlewares.ValidBody.IsValidName(requestData.Name) {
		utils.SendErrorResponse(w, "Name must be more than 4 characters long and can only contain alphabets and spaces", http.StatusBadRequest)
		return
	}

	if !middlewares.ValidBody.IsValidEmail(email) {
		utils.SendErrorResponse(w, "Email must contain only letters, numbers, dots, or hyphens", http.StatusBadRequest)
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
		utils.SendErrorResponse(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	// check if the user already exists
	_, err = db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{Username: username, Email: email})
	if err == nil {
		utils.SendErrorResponse(w, "User already exists", http.StatusConflict)
		return
	}
	err = db.SQL.CreateUser(ctx, &q.CreateUserParams{Username: username, Name: requestData.Name, Email: email, Password: hashedPassword})
	if err != nil {
		utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user, err := db.SQL.FindUserByUsername(ctx, username)
	if err != nil {
		utils.SendErrorResponse(w, "Error getting user", http.StatusInternalServerError)
		return
	}
	s := &models.SafeUser{
		Username:         user.Username,
		Name:             user.Name,
		Email:            user.Email,
		IsPremium:        user.IsPremium,
		AliasCount:       user.AliasCount,
		DestinationCount: user.DestinationCount,
		Avatar:           user.Avatar,
	}
	utils.CreateSendResponse(w, s, `User Created Successfully`, http.StatusCreated, `user`, user.Username)

}
