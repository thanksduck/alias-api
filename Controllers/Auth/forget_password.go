package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

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
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if !middlewares.ValidBody.IsValidEmail(requestData.Email) {
		utils.SendErrorResponse(w, "Email can't be processed", http.StatusUnprocessableEntity)
		return
	}

	user, err := repository.FindUserByUsernameOrEmail("", requestData.Email)
	if err != nil {
		utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}
	if !user.EmailVerified {
		utils.SendErrorResponse(w, "Email not verified", http.StatusUnauthorized)
		return
	}
	token, err := utils.GeneratePasswordResetToken(user.Username)
	if err != nil {
		utils.SendErrorResponse(w, "Error generating Reset Link", http.StatusInternalServerError)
		return
	}
	hash := sha256.New()

	// Write the token to the hash
	hash.Write([]byte(token))

	// Compute the SHA-256 checksum and get the hashed bytes
	hashedBytes := hash.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string
	hashedToken := hex.EncodeToString(hashedBytes)

	err = repository.SavePasswordResetToken(user.ID, hashedToken)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error Processing Reset Link", http.StatusInternalServerError)
		return
	}
	resetURL := r.Referer() + "auth/reset-password/" + hashedToken
	message := "Dear " + user.Name + "\n\n" +
		"You have requested to reset your password. Please click on the link below to reset your password. This link is valid for 10 minutes.\n\n" +
		resetURL + "\n\n" +
		"Thank you,\n" +
		"One Alias Service Team"

	err = utils.SendEmail(requestData.Email, "Password Reset Link", message)
	if err != nil {
		utils.SendErrorResponse(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	// send ok response with message that email has been sent and close the connection and it must be in json
	response := map[string]string{"message": "Password reset link has been sent to your email", "status": "success"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
