package rules

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

func CreateRule(w http.ResponseWriter, r *http.Request) {
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

	user, ok := utils.GetUserFromContext(r.Context())
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
		utils.SendErrorResponse(w, "Alias Cant be Proccessed", http.StatusUnprocessableEntity)
		return
	}
	destination := strings.ToLower(requestBody.DestinationEmail)
	if !middlewares.ValidBody.IsValidEmail(destination) {
		utils.SendErrorResponse(w, "Destination Cant be Proccessed", http.StatusUnprocessableEntity)
		return
	}
	domain := strings.Split(alias, "@")[1]

	savedDestination, err := repository.FindDestinationByEmailAndDomain(destination, domain)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Destination not found", http.StatusNotFound)
		return
	}

	if savedDestination.Username != user.Username {
		utils.SendErrorResponse(w, "Destination not found", http.StatusNotFound)
		return
	}

	if !savedDestination.Verified {
		utils.SendErrorResponse(w, "Destination not verified", http.StatusForbidden)
		return
	}

	existingRule, err := repository.FindRuleByAliasEmail(alias)
	if err != nil && err != pgx.ErrNoRows {
		utils.SendErrorResponse(w, "Error checking Alias Existence", http.StatusInternalServerError)
		return
	}
	if existingRule != nil {
		utils.SendErrorResponse(w, "Email is already Taken 🥲", http.StatusConflict)
		return
	}

	newRule := &models.Rule{
		UserID:           user.ID,
		Username:         user.Username,
		AliasEmail:       alias,
		DestinationEmail: destination,
		Active:           true,
		Comment:          requestBody.Comment,
		Name:             requestBody.Name,
	}
	err = requests.CreateRuleRequest(`POST`, newRule.AliasEmail, newRule.DestinationEmail, newRule.Username, domain)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	_, err = repository.CreateNewRule(newRule)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, newRule, "Rule Created Successfully", http.StatusCreated, "rule", user.Username)

}
