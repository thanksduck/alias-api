package destinations

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("Request body decode error:", err)
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		fmt.Println("User context error: user not found in context")
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(requestBody.Password, user.Password) {
		fmt.Println("Password verification failed for user:", user.Username)
		utils.SendErrorResponse(w, "Password is required to delete a destination address", http.StatusUnauthorized)
		return
	}

	destinationIDStr := r.PathValue("id")
	destinationIDInt, err := strconv.Atoi(destinationIDStr)
	if err != nil {
		fmt.Println("Destination ID conversion error:", err)
		utils.SendErrorResponse(w, "Invalid destination ID", http.StatusBadRequest)
		return
	}
	destinationID := uint32(destinationIDInt)
	destination, err := repository.FindDestinationByID(destinationID)
	if err != nil {
		fmt.Println("Find destination error:", err)
		utils.SendErrorResponse(w, "Destination not found", http.StatusNotFound)
		return
	}
	if destination.Username != user.Username {
		fmt.Println("Unauthorized deletion attempt by", user.Username, "for destination owned by", destination.Username)
		utils.SendErrorResponse(w, "You are not allowed to delete this destination", http.StatusForbidden)
		return
	}

	_, err = requests.DestinationRequest(`DELETE`, destination.DestinationEmail, destination.Username, destination.CloudflareDestinationID)
	if err != nil {
		fmt.Println("Destination DELETE request error:", err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = repository.DeleteDestinationByID(destinationID, user.ID)
	if err != nil {
		fmt.Println("Delete destination from DB error:", err)
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
		fmt.Println("Find active rules error:", err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	sem := make(chan struct{}, 10)
	errChan := make(chan error, len(rules))
	var wg sync.WaitGroup

	for _, rule := range rules {
		wg.Add(1)
		go func(rule models.Rule) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			err := requests.CreateRuleRequest(`PATCH`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain)
			if err != nil {
				fmt.Println("Rule PATCH request error for rule ID", rule.ID, ":", err)
				errChan <- err
				return
			}

			err = repository.ToggleRuleByID(rule.ID)
			if err != nil {
				fmt.Println("Toggle rule error for rule ID", rule.ID, ":", err)
				errChan <- err
				return
			}
		}(rule)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error in goroutine:", err)
			utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	utils.CreateSendResponse(w, nil, "Destination Deleted Successfully", http.StatusNoContent, "destination", user.Username)
}
