package user

import (
	"net/http"

	"github.com/thanksduck/alias-api/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}

	utils.CreateSendResponse(w, user, "User Retreived Successfully", http.StatusOK, "user", user.ID)

}
