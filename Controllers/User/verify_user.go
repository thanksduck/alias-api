package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	auth "github.com/thanksduck/alias-api/Controllers/Auth"
	emailtemplate "github.com/thanksduck/alias-api/Email_Template"
	repository "github.com/thanksduck/alias-api/Repository"
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
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if user.EmailVerified || user.Provider == `localVerified` || user.Provider == `github` || user.Provider == `google` {
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

	err = repository.UpdateProviderByID(user.ID, token)
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

	username := r.PathValue("username")
	token := r.PathValue("token")

	if username == "" || token == "" {
		utils.SendErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := repository.FindUserByUsernameOrEmail(username, ``)
	if err != nil {
		utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Provider != token {
		utils.SendErrorResponse(w, "Invalid token", http.StatusBadRequest)
		return
	}

	err = repository.VerifyEmailByID(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	err = repository.UpdateProviderByID(user.ID, `localVerified`)
	if err != nil {
		utils.SendErrorResponse(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	user.EmailVerified = true
	auth.RedirectToFrontend(w, r, user)
}
