package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	// get the username , name , email , password , passwordConfirm from the request payloas

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

	if !middlewares.ValidBody.IsValidUsername(username) {
		utils.SendErrorResponse(w, "Username must be at least 4 characters long, start with an alphabet, and can contain numbers, dots, or underscores", http.StatusBadRequest)
		return
	}

	if !middlewares.ValidBody.IsValidName(requestData.Name) {
		utils.SendErrorResponse(w, "Name must be more than 4 characters long and can only contain alphabets and spaces", http.StatusBadRequest)
		return
	}

	if !middlewares.ValidBody.IsValidEmail(requestData.Email) {
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
	// create a new user object
	user := &models.User{
		Username: username,
		Name:     requestData.Name,
		Email:    requestData.Email,
		Password: hashedPassword,
		Provider: "local",
		Avatar:   requestData.Avatar,
	}

	// check if the user already exists
	_, err = repository.FindUserByUsernameOrEmail(user.Username, user.Email)
	if err == nil {
		utils.SendErrorResponse(w, "User already exists", http.StatusConflict)
		return
	} else if err != pgx.ErrNoRows {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	newUser, err := repository.CreateUser(user)
	if err != nil {
		utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, newUser, `User Created Successfully`, http.StatusCreated, `user`, newUser.Username)

}
