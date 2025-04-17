package destinations

import (
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"

	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func DeleteDestination(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}

	savedPassword, _ := db.SQL.FindPasswordById(ctx, user.ID)
	if !utils.CheckPassword(requestBody.Password, savedPassword) {
		utils.SendErrorResponse(w, "Password is required to delete a destination address", http.StatusUnauthorized)
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

	_, err = requests.DestinationRequest(`DELETE`, destination.Domain, user.Username, destination.CloudflareDestinationID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = db.SQL.DeleteDestinationByID(ctx, destinationID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	domain := destination.Domain
	rules, err := db.SQL.FindActiveRulesByDestinationEmail(ctx, destination.DestinationEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.CreateSendResponse(w, nil, "Destination Deleted Successfully", http.StatusNoContent, "destination", user.Username)
			return
		}
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Set up error handling
	errChan := make(chan error, len(rules)) // Buffer sized to match maximum possible errors
	var wg sync.WaitGroup

	// Create semaphore for limiting concurrent goroutines
	sem := make(chan struct{}, 10)

	// Process rules in parallel
	for _, rule := range rules {
		wg.Add(1)
		// Capture rule in the closure correctly
		go func(rule *q.FindActiveRulesByDestinationEmailRow) {
			defer wg.Done()
			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			// Update rule
			if err := requests.CreateRuleRequest(`PATCH`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain); err != nil {
				errChan <- fmt.Errorf("failed to update rule %d: %w", rule.ID, err)
				return
			}

			// Toggle rule
		}(rule)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Close error channel after all goroutines are done
	close(errChan)

	// Check for any errors
	var hasError bool
	for err := range errChan {
		hasError = true
		log.Printf("Error processing rule: %v", err)
	}

	if hasError {
		utils.SendErrorResponse(w, "Something went wrong while processing rules", http.StatusInternalServerError)
		return
	}

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to begin transaction: %s", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)
	qtx := q.New(tx)
	err = qtx.MakeAllRuleInactiveByDestinationEmail(ctx, destination.DestinationEmail)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	_ = qtx.DecrementUserDestinationCount(ctx, user.ID)
	err = tx.Commit(ctx)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, nil, "Destination Deleted Successfully", http.StatusNoContent, "destination", user.Username)
}
