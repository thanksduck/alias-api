package destinations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func CreateDestination(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if !user.EmailVerified {
		utils.SendErrorResponse(w, "Please Verify Your Email to add Destination", http.StatusUnauthorized)
		return
	}
	if user.DestinationCount == 1 && !user.IsPremium {
		utils.SendErrorResponse(w, "You have Reached the Destination Limit", http.StatusPaymentRequired)
		return
	}

	var requestBody struct {
		DestinationEmail string `json:"destinationEmail"`
		Domain           string `json:"domain"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.DestinationEmail == "" || requestBody.Domain == "" {
		utils.SendErrorResponse(w, "Both destination and domain are required", http.StatusBadRequest)
		return
	}
	destinationEmail := strings.ToLower(requestBody.DestinationEmail)
	if !middlewares.ValidBody.IsValidEmail(destinationEmail) {
		utils.SendErrorResponse(w, "Destination Cant be Proccessed", http.StatusUnprocessableEntity)
		return
	}
	domain := strings.ToLower(requestBody.Domain)

	destinationResponse, err := requests.DestinationRequest(`POST`, domain, destinationEmail, ``)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error creating destination", http.StatusInternalServerError)
		return
	}

	newDest := &models.Destination{
		Username:                user.Username,
		UserID:                  user.ID,
		DestinationEmail:        destinationEmail,
		Domain:                  domain,
		CloudflareDestinationID: destinationResponse.Result.ID,
		Verified:                false,
	}

	newDestination, err := repository.CreateDestination(newDest)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error creating destination", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, newDestination, "Destination Created Successfully", http.StatusCreated, "destination", user.ID)
}
