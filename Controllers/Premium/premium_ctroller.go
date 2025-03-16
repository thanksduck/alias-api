package premium

import (
	"encoding/json"
	"net/http"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	models "github.com/thanksduck/alias-api/Models"
	"github.com/thanksduck/alias-api/paymentutils"
	"github.com/thanksduck/alias-api/utils"
)

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	var requestBody models.PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middlewares.ValidatePaymentBody(requestBody); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	url, err := paymentutils.InitialisePaymentAndRedirect(&requestBody, user)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, url, "Payment Initialised Successfully", http.StatusOK, "payment", user.Username)
}
