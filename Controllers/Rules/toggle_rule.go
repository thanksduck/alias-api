package rules

import (
	"fmt"
	db "github.com/thanksduck/alias-api/Database"
	"net/http"
	"strconv"
	"strings"

	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func ToggleRule(w http.ResponseWriter, r *http.Request) {
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
		utils.SendErrorResponse(w, "You are not allowed to toggle this rule", http.StatusForbidden)
		return
	}
	_, err = db.SQL.FindDestinationByEmail(ctx, rule.DestinationEmail)
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

	err = db.SQL.ToggleRuleByID(ctx, ruleID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	rule.IsActive = !rule.IsActive

	utils.CreateSendResponse(w, rule, "Rule Toggled Successfully", http.StatusOK, "rule", user.Username)
}
