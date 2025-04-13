package rules

import (
	"fmt"
	"net/http"
	"strings"

	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"

	"strconv"

	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/utils"
)

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	ruleIDStr := r.PathValue("id")
	ruleIDInt, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}
	ruleID := int64(ruleIDInt)
	rule, err := db.SQL.FindRuleByID(ctx, ruleID)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Rule not found", http.StatusNotFound)
		return
	}
	if rule.Username != user.Username {
		utils.SendErrorResponse(w, "You are not allowed to delete this rule", http.StatusForbidden)
		return
	}

	domain := strings.Split(rule.AliasEmail, "@")[1]
	err = requests.CreateRuleRequest(`DELETE`, rule.AliasEmail, rule.DestinationEmail, rule.Username, domain)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to begin transaction: %s", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)
	qtx := q.New(tx)

	err = qtx.DeleteRuleByID(ctx, ruleID)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	_ = qtx.DecrementUserAliasCount(ctx, user.ID)
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, nil, "Rule Deleted Successfully", http.StatusNoContent, "rule", user.Username)
}
