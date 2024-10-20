package destinations

import (
	"net/http"
	"strconv"

	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func ListDestinations(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if user.DestinationCount == 0 {
		utils.SendErrorResponse(w, "You have not created any destinations yet", http.StatusBadRequest)
		return
	}
	destinations, err := repository.FindDestinationsByUserID(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, destinations, "Destinations Retreived Successfully", http.StatusOK, "destinations", user.Username)
}

func GetDestination(w http.ResponseWriter, r *http.Request) {
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
		utils.SendErrorResponse(w, "You are not allowed to view this destination", http.StatusForbidden)
		return
	}
	utils.CreateSendResponse(w, destination, "Destination Retreived Successfully", http.StatusOK, "destination", user.Username)
}
