package premium

import (
	"fmt"
	"net/http"

	"github.com/thanksduck/alias-api/utils"
)

func SubscribeToStar(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not Found", http.StatusUnauthorized)
	}
	fmt.Println(user)
}
