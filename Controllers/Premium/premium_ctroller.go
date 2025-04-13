package premium

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"strings"
	"time"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	models "github.com/thanksduck/alias-api/Models"
	"github.com/thanksduck/alias-api/paymentutils"
	"github.com/thanksduck/alias-api/utils"
)

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
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
	url, err := paymentutils.InitialisePaymentAndRedirect(ctx, &requestBody, user)
	fmt.Printf("%v", err)
	if err != nil {
		utils.SendErrorResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, url, "Payment Initialised Successfully", http.StatusOK, "payment", user.Username)
}

// ProcessPaymentAndSubscribe handles payment verification and subscription creation
func VerifyPaymentAndSubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := utils.GetUserFromContext(ctx)
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

	// Find payment by transaction ID
	payment, err := db.SQL.FindPaymentByTxnID(ctx, txnID)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to find payment: %s", err), http.StatusInternalServerError)
		return
	}

	// If payment already successful, return existing subscription
	if payment.Status == "success" {
		subs, err := db.SQL.GetSubscriptionByUserID(ctx, user.ID)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to get subscription: %s", err), http.StatusInternalServerError)
			return
		}
		utils.CreateSendResponse(w, subs, "Subscription Retrieved Successfully", http.StatusAccepted, "subscription", user.Username)
		return
	}

	// Verify payment with retry
	paymentStatus, err := verifyWithRetry(txnID)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("Failed to verify payment: %s", err), http.StatusInternalServerError)
		return
	}

	// Handle payment based on status
	if paymentStatus == "SUCCESS" {
		// Determine plan type and duration
		var months int
		var planType string

		switch plan {
		case "star":
			planType = "star"
			months = 1
		case "galaxy":
			planType = "galaxy"
			months = 1
		default:
			planType = "star"
			months = 1
		}

		// Begin database transaction
		tx, err := db.DB.Begin(ctx)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to begin transaction: %s", err), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(ctx)

		qtx := q.New(tx)

		// Update payment status
		err = qtx.UpdatePaymentStatus(ctx, &q.UpdatePaymentStatusParams{
			Status: strings.ToLower(paymentStatus),
			ID:     payment.ID,
		})
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to update payment status: %s", err), http.StatusInternalServerError)
			return
		}

		// Check if credit exists for user
		var creditID int64
		credit, err := qtx.FindCreditByUserID(ctx, user.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Create credit if doesn't exist
				creditID, err = qtx.CreateCredit(ctx, &q.CreateCreditParams{
					UserID:  user.ID,
					Balance: 0, // Starting balance
				})
				if err != nil {
					utils.SendErrorResponse(w, fmt.Sprintf("Failed to create credit: %s", err), http.StatusInternalServerError)
					return
				}
			} else {
				utils.SendErrorResponse(w, fmt.Sprintf("Failed to find credit: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			creditID = credit.ID

			// Update credit balance
			err = qtx.UpdateCreditBalance(ctx, &q.UpdateCreditBalanceParams{
				Balance: payment.Amount, // Add payment amount to balance
				ID:      creditID,
			})
			if err != nil {
				utils.SendErrorResponse(w, fmt.Sprintf("Failed to update credit balance: %s", err), http.StatusInternalServerError)
				return
			}
		}

		// Calculate expiration date
		expiresAt := time.Now().AddDate(0, months, 0)

		// Create subscription
		err = qtx.CreateSubscription(ctx, &q.CreateSubscriptionParams{
			UserID:    user.ID,
			CreditID:  creditID,
			Plan:      planType,
			Price:     payment.Amount,
			ExpiresAt: expiresAt,
			Status:    "active",
		})
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to create subscription: %s", err), http.StatusInternalServerError)
			return
		}

		// Update user to premium
		err = qtx.UpdateUserToPremium(ctx, user.ID)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to update user to premium: %s", err), http.StatusInternalServerError)
			return
		}

		// Commit transaction
		err = tx.Commit(ctx)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Failed to commit transaction: %s", err), http.StatusInternalServerError)
			return
		}

		// Get the newly created subscription to return in response
		subscription, err := db.SQL.GetSubscriptionByUserID(ctx, user.ID)
		if err != nil {
			utils.SendErrorResponse(w, fmt.Sprintf("Subscription created but failed to retrieve it: %s", err), http.StatusInternalServerError)
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
