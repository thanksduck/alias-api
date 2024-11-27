package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func ForgetPassword(w http.ResponseWriter, r *http.Request) {
	// we will first extract the email from the request payload
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
		// Find user by email (if exists)
		user, err := repository.FindUserByUsernameOrEmail("", requestData.Email)
		if err != nil {
			// If user not found, do nothing silently
			return
		}

		// Additional checks
		if !user.EmailVerified {
			return
		}

		err = repository.HasNoActiveResetToken(user.ID)
		if err != nil {
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

		err = repository.SavePasswordResetToken(user.ID, user.Username, hashedToken)
		if err != nil {
			return
		}

		resetURL := os.Getenv("FRONTEND_HOST") + "/reset-password/" + hashedToken
		message := "Dear " + user.Name + "\n\n" +
			"You have requested to reset your password. Please click on the link below to reset your password. This link is valid for 10 minutes.\n\n" +
			resetURL + "\n\n" +
			"Thank you,\n" +
			"One Alias Service Team"

		utils.SendEmail(requestData.Email, "Password Reset Link", message)
	}()

	// Always send the same response, regardless of user existence
	utils.SendSuccessResponse(w, "Reset Link Has Been Sent if it exists. Do check spam folder as well")
}

/*

	auth := smtp.PlainAuth("", "user@example.com", "password", "mail.example.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{"recipient@example.net"}
	msg := []byte("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail("mail.example.com:25", auth, "sender@example.org", to, msg)
	if err != nil {
		log.Fatal(err)
	}
*/
