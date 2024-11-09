package destinations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
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
	if !user.IsPremium && destinationEmail != user.Email {
		utils.SendErrorResponse(w, "Destination Email should be same", http.StatusUnprocessableEntity)
		return
	}

	domain := strings.ToLower(requestBody.Domain)

	// Check if destination already exists
	destination, err := repository.FindDestinationByEmail(destinationEmail)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error finding destination", http.StatusInternalServerError)
		return
	}
	if destination != nil {
		if destination.Username != user.Username {
			utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		utils.SendErrorResponse(w, "Destination already exists, Please Verify if You havnt", http.StatusConflict)
		return
	}

	destinationResponse, err := requests.DestinationRequest(`POST`, domain, destinationEmail, ``)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error creating destination", http.StatusInternalServerError)
		return
	}
	var verificationCheck bool
	if destinationResponse.Result.Verified.IsZero() {
		verificationCheck = false
	} else {
		verificationCheck = true
	}

	newDest := &models.Destination{
		Username:                user.Username,
		UserID:                  user.ID,
		DestinationEmail:        destinationEmail,
		Domain:                  domain,
		CloudflareDestinationID: destinationResponse.Result.ID,
		Verified:                verificationCheck,
	}

	newDestination, err := repository.CreateDestination(newDest)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error creating destination", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, newDestination, "Destination Created Successfully", http.StatusCreated, "destination", user.Username)
}
