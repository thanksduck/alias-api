package user

import (
	"encoding/json"
	"net/http"

	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
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

	if !utils.CheckPassword(requestData.Password, user.Password) {
		utils.SendErrorResponse(w, "Password is incorrect", http.StatusBadRequest)
		return
	}

	err = repository.DeleteUser(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"message": "User deleted successfully", "status": "success"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)

}
