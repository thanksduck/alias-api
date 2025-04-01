package premium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
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

// ProcessPaymentAndSubscribe handles payment verification and subscription creation
func VerifyPaymentAndSubscribe(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	var requestBody struct {
		TxnID string `json:"txnId"`
		Plan  string `json:"plan"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	txnID := requestBody.TxnID
	plan := requestBody.Plan

	if txnID == "" {
		utils.SendErrorResponse(w, "Transaction ID is required", http.StatusBadRequest)
		return
	}

	if plan == "" {
		utils.SendErrorResponse(w, "Plan is required", http.StatusBadRequest)
		return
	}
	payment, err := repository.FindPaymentByTxnID(txnID)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to find payment: %s", err), http.StatusInternalServerError)
		return
	}
	if payment.Status == "success" {
		subs, err := repository.GetSubscriptionByUserID(user.ID)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed: %s", err), http.StatusInternalServerError)
		}
		utils.CreateSendResponse(w, subs, "Subscription Reterieved Successfully", http.StatusAccepted, "subscription", user.Username)
		return
	}
	paymentStatus, err := verifyWithRetry(txnID)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to verify payment: %s", err), http.StatusInternalServerError)
		return
	}

	// If payment is successful, create subscription
	if paymentStatus == "SUCCESS" {
		// payment, err := repository.FindPaymentByTxnID(txnID)
		// if err != nil {
		// 	utils.SendErrorResponse(w, fmt.Sprintf("Failed to find payment: %s", err), http.StatusInternalServerError)
		// 	return
		// }
		// Calculate expiration based on plan
		var months int
		var planType models.PlanType

		// Determine plan type and months based on the plan parameter
		switch plan {
		case "star":
			planType = models.StarPlan
			months = 1 // Assuming star plan is for 1 month
		case "galaxy":
			planType = models.GalaxyPlan
			months = 1 // Assuming galaxy plan is for 1 month
		default:
			planType = models.StarPlan // Default to star plan as requested
			months = 1
		}

		// Create subscription object
		subscription := &models.Subscription{
			UserID:    user.ID,
			Plan:      planType,
			Price:     uint32(payment.Amount),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().AddDate(0, months, 0),
			Status:    "active",
		}

		// Update payment status, credit and create subscription in a transaction
		err = repository.UpdatePaymentStatusCreditAndCreateSubscription(subscription, payment, strings.ToLower(paymentStatus))
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to create subscription: %s", err), http.StatusInternalServerError)
			return
		}

		// Return success response
		utils.CreateSendResponse(w, subscription, "Payment successful and subscription created", http.StatusAccepted, "subscription", user.Username)
		return
	} else if paymentStatus == "PENDING" {
		utils.SendErrorResponse(w, "Payment is still pending, please try again later", http.StatusAccepted)
		return
	} else {
		utils.SendErrorResponse(w, "Payment failed", http.StatusBadRequest)
		return
	}
}

// verifyWithRetry implements the retry logic as per PhonePe documentation

func verifyWithRetry(txnID string) (string, error) {
	// Just verify the payment once without retries
	status, err := paymentutils.VerifyPhonePePayment(txnID)
	if err != nil {
		return "", err
	}

	return status, nil
}
