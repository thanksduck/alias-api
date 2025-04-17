package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"os"
	"strings"

	auth "github.com/thanksduck/alias-api/Controllers/Auth"
	emailtemplate "github.com/thanksduck/alias-api/Email_Template"
	"github.com/thanksduck/alias-api/utils"
)

func GenerateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func GenerateVerifyUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if user.IsEmailVerified || user.Provider == `localVerified` || user.Provider == `github` || user.Provider == `google` {
		utils.CreateSendResponse(w, user, "User already verified", http.StatusOK, "user", user.Username)
		return
	}
	if user.Provider != `local` {
		utils.CreateSendResponse(w, user, "Verification email already sent. If you didn't receive your email, check spam or contact support", http.StatusOK, "user", user.Username)
		return
	}

	token, err := GenerateToken()
	if err != nil {
		utils.SendErrorResponse(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	err = db.SQL.UpdateProviderByID(ctx, &q.UpdateProviderByIDParams{Provider: token, ID: user.ID})
	if err != nil {
		utils.SendErrorResponse(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	magicLink := fmt.Sprintf("%s/api/v2/user/%s/verify/%s", os.Getenv("REDIRECT_HOST"), user.Username, token)

	htmlBody, textBody := emailtemplate.VerifyEmailTemplate(user.Name, magicLink)

	err = utils.SendEmail(user.Email, "Verify Your Email", htmlBody, textBody)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, user, "Magic verification link sent, please Check Your Spam as well", http.StatusOK, "user", user.Username)
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := strings.ToLower(r.PathValue("username"))
	token := r.PathValue("token")

	if username == "" || token == "" {
		utils.SendErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, err := db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{Username: username})
	if err != nil {
		utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Provider != token {
		utils.SendErrorResponse(w, "Invalid token", http.StatusBadRequest)
		return
	}
	err = db.SQL.VerifyEmailByID(ctx, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Error verifying email", http.StatusInternalServerError)
		return
	}
	err = db.SQL.UpdateProviderByID(ctx, &q.UpdateProviderByIDParams{Provider: token, ID: user.ID})
	if err != nil {
		utils.SendErrorResponse(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	user.IsEmailVerified = true
	auth.RedirectToFrontend(w, r, user.Username)
}
