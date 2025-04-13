package auth

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

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	// Validate username and email
	if email == "" && username == "" {
		utils.SendErrorResponse(w, "Username or Email is required", http.StatusBadRequest)
		return
	}

	if email != "" && !middlewares.ValidBody.IsValidEmail(email) {
		utils.SendErrorResponse(w, "Email can't be processed", http.StatusUnprocessableEntity)
		return
	}

	if username != "" && !middlewares.ValidBody.IsValidUsername(username) {
		utils.SendErrorResponse(w, "Username can't be processed", http.StatusUnprocessableEntity)
		return
	}

	// Find user by username or email
	user, err := db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{Username: username, Email: email})
	if err != nil {
		utils.SendErrorResponse(w, "Invalid email or password. If you used a social sign-in, Please use that method.", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(requestData.Password, user.Password) {
		utils.SendErrorResponse(w, "Invalid email or password. If you used a social sign-in, Please use that method.", http.StatusUnauthorized)
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

	// if user.IsPremium {
	// 	plan, _ := repository.GetPlanByUserID(user.ID)
	// 	user.Plan = plan
	// }
	utils.CreateSendResponse(w, s, "Login Successful", http.StatusOK, "user", user.Username)

}
