package rules

import (
	"encoding/json"
	"fmt"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"strings"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func CreateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody struct {
		AliasEmail       string `json:"aliasEmail"`
		DestinationEmail string `json:"destinationEmail"`
		Name             string `json:"name"`
		Comment          string `json:"comment"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.AliasEmail == "" || requestBody.DestinationEmail == "" {
		utils.SendErrorResponse(w, "Both alias and destination are required are required", http.StatusBadRequest)
		return
	}

	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if user.DestinationCount == 0 {
		utils.SendErrorResponse(w, "You have not created any destinations yet", http.StatusBadRequest)
		return
	}

	alias := strings.ToLower(requestBody.AliasEmail)
	if !middlewares.ValidBody.IsValidEmail(alias) {
		utils.SendErrorResponse(w, "Alias Cant be Processed", http.StatusUnprocessableEntity)
		return
	}
	destination := strings.ToLower(requestBody.DestinationEmail)
	if !middlewares.ValidBody.IsValidEmail(destination) {
		utils.SendErrorResponse(w, "Destination Cant be Processed", http.StatusUnprocessableEntity)
		return
	}
	domain := strings.Split(alias, "@")[1]

	savedDestination, err := db.SQL.FindDestinationByEmailAndDomainAndUserID(ctx, &q.FindDestinationByEmailAndDomainAndUserIDParams{
		DestinationEmail: destination,
		Domain:           domain,
		UserID:           user.ID,
	})
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Invalid alias Domain or Destination not found", http.StatusNotFound)
		return
	}

	if !savedDestination.IsVerified {
		utils.SendErrorResponse(w, "Destination not verified", http.StatusForbidden)
		return
	}

	_, err = db.SQL.FindRuleByAliasEmail(ctx, alias)
	if err == nil {
		utils.SendErrorResponse(w, "Email is already Taken 🥲", http.StatusConflict)
		return
	}

	newRule := &q.CreateNewRuleParams{
		UserID:           user.ID,
		Username:         user.Username,
		AliasEmail:       alias,
		DestinationEmail: destination,
		Comment:          requestBody.Comment,
		Name:             requestBody.Name,
	}

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to begin transaction: %s", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)
	qtx := q.New(tx)
	err = requests.CreateRuleRequest(`POST`, newRule.AliasEmail, newRule.DestinationEmail, newRule.Username, domain)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, fmt.Sprintf("%s already taken,Please Choose Other Email", newRule.AliasEmail), http.StatusInternalServerError)
		return
	}
	err = qtx.CreateNewRule(ctx, newRule)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	_ = qtx.IncrementUserAliasCount(ctx, user.ID)
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, nil, "Rule Created Successfully", http.StatusCreated, "rule", user.Username)

}
