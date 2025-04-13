package destinations

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"

	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func VerifyDestination(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	destinationIDStr := r.PathValue("id")
	destinationIDInt, err := strconv.Atoi(destinationIDStr)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid destination ID", http.StatusBadRequest)
		return
	}
	destinationID := int64(destinationIDInt)
	destination, err := db.SQL.GetCloudflareDestinationID(ctx, &q.GetCloudflareDestinationIDParams{
		ID:     destinationID,
		UserID: user.ID,
	})
	if err != nil {
		utils.SendErrorResponse(w, "Destination not found", http.StatusNotFound)
		return
	}
	if destination.IsVerified {
		utils.CreateSendResponse(w, destination, "Destination already verified", http.StatusOK, "destination", user.Username)
		return
	}

	destResponse, err := requests.DestinationRequest(`GET`, destination.Domain, user.Username, destination.CloudflareDestinationID)
	if err != nil {
		fmt.Printf("Error: Destination request failed - %v\n", err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if !destResponse.Success {
		fmt.Printf("Error: Destination request unsuccessful\n")
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if destResponse.Result.Verified.IsZero() {
		utils.SendErrorResponse(w, "Please check your mail or spam folder and click on verify", http.StatusBadRequest)
		return
	}
	destination.IsVerified = true

	err = db.SQL.VerifyDestinationByID(ctx, destinationID)
	if err != nil {
		fmt.Printf("Error: Failed to verify destination by ID - %v\n", err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Info: Destination verified successfully\n")
	utils.CreateSendResponse(w, destination, "Destination Verified Successfully", http.StatusOK, "destination", user.Username)
}
