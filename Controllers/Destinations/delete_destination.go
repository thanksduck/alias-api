package destinations

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func DeleteDestination(w http.ResponseWriter, r *http.Request) {
	// check for the password in the request body
	var requestBody struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(requestBody.Password, user.Password) {
		utils.SendErrorResponse(w, "Password is required to delete a destination address", http.StatusUnauthorized)
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
		utils.SendErrorResponse(w, "You are not allowed to delete this destination", http.StatusForbidden)
		return
	}

	_, err = requests.DestinationRequest(`DELETE`, destination.DestinationEmail, destination.Username, destination.CloudflareDestinationID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = repository.DeleteDestinationByID(destinationID, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// we will toggle all the rules to inactive and also toggle them in the rules table
	domain := destination.Domain

	rules, err := repository.FindActiveRulesByDestinationEmail(destination.DestinationEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No rules found, we can proceed
			utils.CreateSendResponse(w, nil, "Destination Deleted Successfully", http.StatusNoContent, "destination", user.Username)
			return
		}
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Create a semaphore with buffer of 10 to limit concurrent operations
	sem := make(chan struct{}, 10)
	errChan := make(chan error, len(rules))
	var wg sync.WaitGroup

	for _, rule := range rules {
		wg.Add(1)
		go func(rule models.Rule) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			// Make the API request
			err := requests.CreateRuleRequest(`PATCH`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain)
			if err != nil {
				errChan <- err
				return
			}

			// Toggle the rule in database
			err = repository.ToggleRuleByID(rule.ID)
			if err != nil {
				errChan <- err
				return
			}
		}(rule)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	utils.CreateSendResponse(w, nil, "Destination Deleted Successfully", http.StatusNoContent, "destination", user.Username)
}
