package user

import (
	models "github.com/thanksduck/alias-api/Models"
	"net/http"

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
