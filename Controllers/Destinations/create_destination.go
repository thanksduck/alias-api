package destinations

import (
	"encoding/json"
	"fmt"
	db "github.com/thanksduck/alias-api/Database"
	requests "github.com/thanksduck/alias-api/Requests"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"strings"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	"github.com/thanksduck/alias-api/utils"
)

func CreateDestination(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if !user.IsEmailVerified {
		utils.SendErrorResponse(w, "Please Verify Your Email to add Destination", http.StatusForbidden)
		return
	}
	if user.DestinationCount >= 1 && !user.IsPremium {
		utils.SendPaymentRequiredResponse(w, "For more Destinations go premium", "Free", 1)
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
		utils.SendErrorResponse(w, "Destination Cant be Processed", http.StatusUnprocessableEntity)
		return
	}
	if !user.IsPremium && destinationEmail != user.Email {
		utils.SendErrorResponse(w, "Destination Email should be same", http.StatusUnprocessableEntity)
		return
	}

	domain := strings.ToLower(requestBody.Domain)

	// Check if destination already exists
	destination, err := db.SQL.FindDestinationByEmailAndDomain(ctx, &q.FindDestinationByEmailAndDomainParams{Domain: domain,
		DestinationEmail: destinationEmail})
	if err == nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Destination already exist", http.StatusConflict)
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

	newDest := &q.CreateDestinationParams{
		Username:                user.Username,
		UserID:                  user.ID,
		DestinationEmail:        destinationEmail,
		Domain:                  domain,
		CloudflareDestinationID: destinationResponse.Result.ID,
		IsVerified:              verificationCheck,
	}
	// Begin database transaction
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to begin transaction: %s", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)
	qtx := q.New(tx)

	err = qtx.CreateDestination(ctx, newDest)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error creating destination", http.StatusInternalServerError)
		return
	}
	err = qtx.IncrementUserDestinationCount(ctx, user.ID)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Error incrementing user destination count", http.StatusInternalServerError)
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to commit transaction"), http.StatusInternalServerError)
	}
	utils.CreateSendResponse(w, nil, "Destination Created Successfully", http.StatusCreated, "destination", user.Username)
}
