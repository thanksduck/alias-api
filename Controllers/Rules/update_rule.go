package rules

import (
	"encoding/json"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"strconv"
	"strings"

	models "github.com/thanksduck/alias-api/Models"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func UpdateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}

	ruleIDStr := r.PathValue("id")
	ruleIDInt, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}
	ruleID := int64(ruleIDInt)
	rule, err := db.SQL.FindRuleByID(ctx, ruleID)
	if err != nil {
		utils.SendErrorResponse(w, "Rule not found", http.StatusNotFound)
		return
	}

	if rule.Username != user.Username {
		utils.SendErrorResponse(w, "You are not allowed to update this rule", http.StatusForbidden)
		return
	}

	var ruleData models.Rule
	err = json.NewDecoder(r.Body).Decode(&ruleData)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if ruleData.AliasEmail == "" || ruleData.DestinationEmail == "" {
		utils.SendErrorResponse(w, "Alias and Destination email cannot be empty", http.StatusBadRequest)
		return
	}

	// Only make backend requests if alias or destination email changes
	if rule.AliasEmail != ruleData.AliasEmail || rule.DestinationEmail != ruleData.DestinationEmail {
		domain := strings.Split(rule.AliasEmail, "@")[1]
		err = requests.CreateRuleRequest(`DELETE`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain)
		if err != nil {
			utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		err = requests.CreateRuleRequest(`POST`, ruleData.AliasEmail, ruleData.DestinationEmail, rule.Username, strings.Split(ruleData.AliasEmail, "@")[1])
		if err != nil {
			utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	rule.AliasEmail = ruleData.AliasEmail
	rule.DestinationEmail = ruleData.DestinationEmail
	rule.Comment = ruleData.Comment
	rule.Name = ruleData.Name

	updatedRule := &q.UpdateRuleByIDParams{
		Name:             rule.Name,
		DestinationEmail: rule.DestinationEmail,
		AliasEmail:       rule.AliasEmail,
		Comment:          rule.Comment,
		ID:               ruleID,
	}
	err = db.SQL.UpdateRuleByID(ctx, updatedRule)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, updatedRule, "Rule Updated Successfully", http.StatusOK, "rule", user.Username)
}
