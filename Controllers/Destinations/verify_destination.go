package destinations

import (
	"net/http"
	"strconv"

	repository "github.com/thanksduck/alias-api/Repository"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func VerifyDestination(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
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
	destinationID := uint32(destinationIDInt)
	destination, err := repository.FindDestinationByID(destinationID)
	if err != nil {
		utils.SendErrorResponse(w, "Destination not found", http.StatusNotFound)
		return
	}
	if destination.Username != user.Username {
		utils.SendErrorResponse(w, "You are not allowed to verify this destination", http.StatusForbidden)
		return
	}

	if destination.Verified {
		utils.CreateSendResponse(w, destination, "Destination already verified", http.StatusOK, "destination", user.Username)
		return
	}

	destResponse, err := requests.DestinationRequest(`PATCH`, destination.DestinationEmail, destination.Username, destination.CloudflareDestinationID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if !destResponse.Success {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if destResponse.Result.Verified.IsZero() {
		utils.SendErrorResponse(w, "Please check your mail or spam folder and click on verify", http.StatusBadRequest)
		return
	}
	destination.Verified = true

	err = repository.VerifyDestinationByID(destinationID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, destination, "Destination Verified Successfully", http.StatusOK, "destination", user.Username)

}
