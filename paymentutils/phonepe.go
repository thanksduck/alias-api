package paymentutils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
	"resty.dev/v3"
)

type PhonePePaymentRequest struct {
	MerchantID            string `json:"merchantId"`
	MerchantTransactionID string `json:"merchantTransactionId"`
	MerchantUserID        string `json:"merchantUserId"`
	Amount                int    `json:"amount"`
	RedirectURL           string `json:"redirectUrl"`
	RedirectMode          string `json:"redirectMode"`
	CallbackURL           string `json:"callbackUrl"`
	// MobileNumber          string `json:"mobileNumber"`
	PaymentInstrument struct {
		Type string `json:"type"`
	} `json:"paymentInstrument"`
}

type PhonePePaymentResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MerchantID            string `json:"merchantId"`
		MerchantTransactionID string `json:"merchantTransactionId"`
		InstrumentResponse    struct {
			Type         string `json:"type"`
			RedirectInfo struct {
				URL    string `json:"url"`
				Method string `json:"method"`
			} `json:"redirectInfo"`
		} `json:"instrumentResponse"`
	} `json:"data"`
}
type PhonePeStatusResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MerchantID            string `json:"merchantId"`
		MerchantTransactionID string `json:"merchantTransactionId"`
		TransactionID         string `json:"transactionId"`
		Amount                int64  `json:"amount"`
		State                 string `json:"state"`
		ResponseCode          string `json:"responseCode"`
		PaymentInstrument     struct {
			Type string `json:"type"`
			// Additional fields based on payment type
		} `json:"paymentInstrument"`
	} `json:"data"`
}

func InitialisePaymentAndRedirect(requestBody *models.PaymentRequest, user *models.User) (url string, err error) {
	merchantID := os.Getenv("PHONEPE_MERCHENT_ID")
	txnID := GenerateTransactionID(user.Username, requestBody.Plan)

	amount := GetMonthlyPrice(requestBody.Plan, requestBody.Months) * 100 * requestBody.Months

	phonePeRequest := PhonePePaymentRequest{
		MerchantID:            merchantID,
		MerchantTransactionID: txnID,
		MerchantUserID:        user.Username,
		Amount:                amount,
		RedirectURL:           fmt.Sprintf("%s/pay/cb?txnid=%s&plan=%s", os.Getenv("FRONTEND_HOST"), txnID, requestBody.Plan),
		CallbackURL:           fmt.Sprintf("%s/webhook/phonepay", os.Getenv("REDIRECT_HOST")),
		RedirectMode:          "REDIRECT",
	}
	phonePeRequest.PaymentInstrument.Type = "PAY_PAGE"

	endpoint := "/pg/v1/pay"
	reqUrl := GetPhonePeBaseURL()

	jsonData, err := json.Marshal(phonePeRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

	base64Payload := base64.StdEncoding.EncodeToString(jsonData)
	xverify := GenerateXVerifyHeader(base64Payload, endpoint)

	reqBody, err := json.Marshal(map[string]string{
		"request": base64Payload,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", reqUrl+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-VERIFY", xverify)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var phonePeResp PhonePePaymentResponse
	if err := json.Unmarshal(respBody, &phonePeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !phonePeResp.Success {
		return "", fmt.Errorf("payment initialization failed: %s - %s", phonePeResp.Code, phonePeResp.Message)
	}

	redirectURL := phonePeResp.Data.InstrumentResponse.RedirectInfo.URL
	if redirectURL == "" {
		return "", fmt.Errorf("empty redirect URL in response")
	}
	// now we have the transaction id we will create the payment in the database behind the scenes in go routine
	newPayment := models.Payment{
		UserID:  user.ID,
		Type:    "credit",
		Gateway: "phonepe",
		TxnID:   txnID,
		Amount:  int64(amount / 100),
		Status:  "pending",
	}

	err = repository.InitialisePayment(&newPayment)
	if err != nil {
		return "", fmt.Errorf("failed to create payment record: %w", err)
	}

	return redirectURL, nil
}

func VerifyPhonePePayment(txnID string) (string, error) {
	client := resty.New()
	defer client.Close()
	merchantID := os.Getenv("PHONEPE_MERCHENT_ID")
	endpoint := fmt.Sprintf("/pg/v1/status/%s/%s", merchantID, txnID)

	// Generate X-VERIFY header
	xverify := hashString(endpoint+os.Getenv("PHONEPE_SALT")) + "###" + os.Getenv("PHONEPE_SALT_INDEX")
	baseURL := GetPhonePeBaseURL()

	// Make the request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-MERCHANT-ID", merchantID).
		SetHeader("X-VERIFY", xverify).
		SetResult(&PhonePeStatusResponse{}).
		Get(baseURL + endpoint)

	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("received error response: %s", resp.String())
	}

	statusResp, ok := resp.Result().(*PhonePeStatusResponse)
	if !ok || statusResp == nil {
		return "", fmt.Errorf("failed to parse response")
	}

	// Determine payment status based on response
	var paymentStatus string

	switch statusResp.Code {
	case "PAYMENT_SUCCESS":
		paymentStatus = "SUCCESS"
	case "PAYMENT_ERROR", "PAYMENT_DECLINED", "TIMED_OUT":
		paymentStatus = "FAILED"
	case "PAYMENT_PENDING":
		paymentStatus = "PENDING"
	case "INTERNAL_SERVER_ERROR", "BAD_REQUEST", "AUTHORIZATION_FAILED", "TRANSACTION_NOT_FOUND":
		return "", fmt.Errorf("payment verification error: %s - %s", statusResp.Code, statusResp.Message)
	default:
		return "", fmt.Errorf("unknown response code: %s", statusResp.Code)
	}

	return paymentStatus, nil
}
