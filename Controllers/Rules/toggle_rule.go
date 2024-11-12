package rules

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	repository "github.com/thanksduck/alias-api/Repository"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func ToggleRule(w http.ResponseWriter, r *http.Request) {
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
		utils.SendErrorResponse(w, "You are not allowed to toggle this rule", http.StatusForbidden)
		return
	}

	_, err = repository.FindDestinationByEmail(rule.DestinationEmail)
	if err != nil {
		message := fmt.Sprintf("The destination email %s has been removed by You. You can update the rule name and comment, but the rule must remain inactive until same destination is not added", rule.DestinationEmail)
		utils.SendErrorResponse(w, message, http.StatusNotFound)
		return
	}

	domain := strings.Split(rule.AliasEmail, "@")[1]
	err = requests.CreateRuleRequest(`PATCH`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = repository.ToggleRuleByID(ruleID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	rule.Active = !rule.Active

	utils.CreateSendResponse(w, rule, "Rule Toggled Successfully", http.StatusOK, "rule", user.Username)
}
