package rules

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func UpdateRule(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
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
	ruleID := uint32(ruleIDInt)
	rule, err := repository.FindRuleByID(ruleID)
	if err != nil {
		utils.SendErrorResponse(w, "Rule not found", http.StatusNotFound)
		return
	}

	if rule.Username != user.Username {
		utils.SendErrorResponse(w, "You are not allowed to update this rule", http.StatusForbidden)
		return
	}

	domain := strings.Split(rule.AliasEmail, "@")[1]

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
	err = requests.CreateRuleRequest(`DELETE`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	rule.AliasEmail = ruleData.AliasEmail
	rule.DestinationEmail = ruleData.DestinationEmail
	rule.Comment = ruleData.Comment
	err = requests.CreateRuleRequest(`POST`, rule.AliasEmail, rule.DestinationEmail, rule.Username, strings.Split(rule.AliasEmail, "@")[1])
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	updatedRule, err := repository.UpdateRuleByID(ruleID, rule)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, updatedRule, "Rule Updated Successfully", http.StatusOK, "rule", user.Username)
}
