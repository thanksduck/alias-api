package rules

import (
	db "github.com/thanksduck/alias-api/Database"
	"net/http"
	"strconv"

	"github.com/thanksduck/alias-api/utils"
)

func ListRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if user.AliasCount == 0 {
		utils.SendErrorResponse(w, "You have not created any aliases yet", http.StatusBadRequest)
		return
	}
	rules, err := db.SQL.FindRulesByUserID(ctx, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	utils.CreateSendResponse(w, rules, "Rules Retrieved Successfully", http.StatusOK, "rules", user.Username)

}

func GetRule(w http.ResponseWriter, r *http.Request) {
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
		utils.SendErrorResponse(w, "You are not allowed to view this rule", http.StatusForbidden)
		return
	}
	utils.CreateSendResponse(w, rule, "Rule Retreived Successfully", http.StatusOK, "rule", user.Username)
}
