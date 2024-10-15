package rules

import (
	"fmt"
	"net/http"

	"strconv"

	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
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
	ruleID := uint32(ruleIDInt)
	rule, err := repository.FindRuleByID(ruleID)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Rule not found", http.StatusNotFound)
		return
	}
	if rule.Username != user.Username {
		utils.SendErrorResponse(w, "You are not allowed to delete this rule", http.StatusForbidden)
		return
	}
	err = repository.DeleteRuleByID(ruleID)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, nil, "Rule Deleted Successfully", http.StatusNoContent, "rule", user.ID)
}
