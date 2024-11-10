package destinations

import (
	"encoding/json"
	"fmt"
	"log"
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

	_, err = requests.DestinationRequest(`DELETE`, destination.Domain, destination.Username, destination.CloudflareDestinationID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = repository.DeleteDestinationByID(destinationID, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	domain := destination.Domain
	rules, err := repository.FindActiveRulesByDestinationEmail(destination.DestinationEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.CreateSendResponse(w, nil, "Destination Deleted Successfully", http.StatusNoContent, "destination", user.Username)
			return
		}
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Set up error handling
	errChan := make(chan error, len(rules))
	var wg sync.WaitGroup

	// Create semaphore for limiting concurrent goroutines
	sem := make(chan struct{}, 10)

	// Process rules in parallel
	for _, rule := range rules {
		wg.Add(1)
		go func(rule models.Rule) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			// Update rule
			err := requests.CreateRuleRequest(`PATCH`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain)
			if err != nil {
				errChan <- fmt.Errorf("failed to update rule %d: %w", rule.ID, err)
				return
			}

			// Toggle rule
			err = repository.ToggleRuleByID(rule.ID)
			if err != nil {
				errChan <- fmt.Errorf("failed to toggle rule %d: %w", rule.ID, err)
				return
			}
		}(rule)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	var hasError bool
	for err := range errChan {
		if err != nil {
			hasError = true
			log.Printf("Error processing rule: %v", err)
		}
	}

	if hasError {
		utils.SendErrorResponse(w, "Something went wrong while processing rules", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, nil, "Destination Deleted Successfully", http.StatusNoContent, "destination", user.Username)
}
