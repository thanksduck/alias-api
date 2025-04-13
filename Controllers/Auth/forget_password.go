package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	db "github.com/thanksduck/alias-api/Database"
	emailtemplate "github.com/thanksduck/alias-api/Email_Template"
	middlewares "github.com/thanksduck/alias-api/Middlewares"
	q "github.com/thanksduck/alias-api/internal/db"
	"github.com/thanksduck/alias-api/utils"
)

func ForgetPassword(w http.ResponseWriter, r *http.Request) {
	// we will first extract the email from the request payload
	ctx := r.Context()
	var requestData struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.SendSuccessResponse(w, "Reset Link Has Been Sent if it exists. Do check spam folder as well")
		return
	}

	requestData.Email = strings.ToLower(requestData.Email)

	// Validate email format first
	if !middlewares.ValidBody.IsValidEmail(requestData.Email) {
		utils.SendSuccessResponse(w, "Reset Link Has Been Sent if it exists. Do check spam folder as well")
		return
	}

	// Use a goroutine to process email reset asynchronously also not send emails dont inform user if email is not found
	go func() {
		user, err := db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{
			Email: requestData.Email,
		})
		if err != nil {
			return
		}

		if !user.IsEmailVerified {
			return
		}

		_, err = db.SQL.HasNoActiveResetToken(ctx,
			user.ID,
		)
		if err == nil {
			// has an active token
			return
		}

		token, err := utils.GeneratePasswordResetToken(user.Username)
		if err != nil {
			return
		}

		hash := sha256.New()
		hash.Write([]byte(token))
		hashedBytes := hash.Sum(nil)
		hashedToken := hex.EncodeToString(hashedBytes)

		err = db.SQL.CreateNewPasswordResetToken(ctx, &q.CreateNewPasswordResetTokenParams{
			UserID:             user.ID,
			Username:           user.Username,
			PasswordResetToken: hashedToken,
		})
		if err != nil {
			return
		}

		resetLink := os.Getenv("FRONTEND_HOST") + "/reset-password/" + hashedToken
		htmlBody, textBody := emailtemplate.ForgetPasswordTemplate(user.Name, resetLink)
		err = utils.SendEmail(requestData.Email, "Password Reset Link | One Alias Service", htmlBody, textBody)
		if err != nil {
			return
		}
	}()

	// Always send the same response, regardless of user existence
	utils.SendSuccessResponse(w, "Reset Link Has Been Sent if it exists. Do check spam folder as well")
}
