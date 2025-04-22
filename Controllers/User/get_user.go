package user

import (
	"net/http"

	db "github.com/thanksduck/alias-api/Database"
	models "github.com/thanksduck/alias-api/Models"

	"github.com/thanksduck/alias-api/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	s := &models.SafeUser{
		Username:         user.Username,
		Name:             user.Name,
		Email:            user.Email,
		IsPremium:        user.IsPremium,
		IsEmailVerified:  user.IsEmailVerified,
		AliasCount:       user.AliasCount,
		DestinationCount: user.DestinationCount,
		Avatar:           user.Avatar,
	}
	utils.CreateSendResponse(w, s, "User Retrieved Successfully", http.StatusOK, "user", user.Username)

}

func GetUserWithPlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	s := &models.SafeUser{
		Username:         user.Username,
		Name:             user.Name,
		Email:            user.Email,
		IsPremium:        user.IsPremium,
		IsEmailVerified:  user.IsEmailVerified,
		AliasCount:       user.AliasCount,
		DestinationCount: user.DestinationCount,
		Avatar:           user.Avatar,
	}

	if user.IsPremium {
		plan, _ := db.SQL.GetPlanByUserID(ctx, user.ID)
		s.Plan = models.PlanType(plan)
	}
	utils.CreateSendResponse(w, s, "User Retrieved Successfully", http.StatusOK, "user", user.Username)

}
