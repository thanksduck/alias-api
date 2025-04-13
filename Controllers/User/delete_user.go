package user

import (
	"encoding/json"
	db "github.com/thanksduck/alias-api/Database"
	"net/http"

	"github.com/thanksduck/alias-api/utils"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	var requestData struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if requestData.Password == "" {
		utils.SendErrorResponse(w, "Password is required", http.StatusBadRequest)
		return
	}
	savedPassword, err := db.SQL.FindPasswordById(ctx, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if !utils.CheckPassword(requestData.Password, savedPassword) {
		utils.SendErrorResponse(w, "Password is incorrect", http.StatusBadRequest)
		return
	}
	err = db.SQL.DeleteUser(ctx, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"message": "User deleted successfully", "status": "success"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
