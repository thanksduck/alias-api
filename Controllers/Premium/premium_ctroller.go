package premium

import (
	"encoding/json"
	"fmt"
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
	fmt.Printf("%v", err)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, url, "Payment Initialised Successfully", http.StatusOK, "payment", user.Username)
}

/*
// ProcessPaymentAndSubscribe handles payment verification and subscription creation
func VerifyPaymentAndSubscribe(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}

	txnID := r.URL.Query().Get("txnid")
	plan := r.URL.Query().Get("plan")
	months := int(r.URL.Query().Get("months"))

	if txnID == "" {
		utils.SendErrorResponse(w, "Transaction ID is required", http.StatusBadRequest)
		return
	}

	if plan == "" {
		utils.SendErrorResponse(w, "Plan is required", http.StatusBadRequest)
		return
	}

	// Implement retry logic as per PhonePe documentation
	paymentStatus, err := verifyWithRetry(txnID)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to verify payment: %s", err), http.StatusInternalServerError)
		return
	}

	// If payment is successful, create subscription
	if paymentStatus == "SUCCESS" {
		payment, err := repository.FindPaymentByTxnID(txnID)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to find payment: %s", err), http.StatusInternalServerError)
			return
		}

		// Calculate expiration based on plan
		var months int
		var planType models.PlanType

		// Create subscription object
		subscription := &models.Subscription{
			UserID:    user.ID,
			Plan:      planType,
			Price:     uint32(payment.Amount),
			ExpiresAt: time.Now().AddDate(0, months, 0),
			Status:    "ACTIVE",
		}

		// Update payment status, credit and create subscription in a transaction
		err = repository.UpdatePaymentStatusCreditAndCreateSubscription(subscription, payment, paymentStatus)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to create subscription: %s", err), http.StatusInternalServerError)
			return
		}

		// Return success response
		utils.SendSuccessResponse(w, map[string]interface{}{
			"message":      "Payment successful and subscription created",
			"subscription": subscription,
		})
		return
	} else if paymentStatus == "PENDING" {
		utils.SendErrorResponse(w, "Payment is still pending, please try again later", http.StatusAccepted)
		return
	} else {
		utils.SendErrorResponse(w, "Payment failed", http.StatusBadRequest)
		return
	}
}
*/
// verifyWithRetry implements the retry logic as per PhonePe documentation
/*
func verifyWithRetry(txnID string) (string, error) {
	// First check after 20-25 seconds (we'll use 20)
	time.Sleep(20 * time.Second)

	status, err := paymentutils.VerifyPhonePePayment(txnID)
	if err != nil {
		return "", err
	}

	if status != "PENDING" {
		return status, nil
	}

	// Retry schedule as per documentation
	intervals := []struct {
		duration time.Duration
		count    int
	}{
		{3 * time.Second, 10}, // Every 3 seconds for 30 seconds
		{6 * time.Second, 10}, // Every 6 seconds for 60 seconds
		{10 * time.Second, 6}, // Every 10 seconds for 60 seconds
		{30 * time.Second, 2}, // Every 30 seconds for 60 seconds
		{1 * time.Minute, 16}, // Every 1 minute until timeout (16 minutes more)
	}

	for _, interval := range intervals {
		for i := 0; i < interval.count; i++ {
			time.Sleep(interval.duration)

			status, err := paymentutils.VerifyPhonePePayment(txnID)
			if err != nil {
				return "", err
			}

			if status != "PENDING" {
				return status, nil
			}
		}
	}

	// If we've reached this point, payment is still pending after 20 minutes
	return "PENDING", nil
}
*/
