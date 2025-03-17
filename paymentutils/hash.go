package paymentutils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	models "github.com/thanksduck/alias-api/Models"
)

const (
	PhonePeSandboxBaseURL = "https://api-preprod.phonepe.com/apis/hermes"
	PhonePeProdBaseURL    = "https://api.phonepe.com/apis/hermes"
)

// hashString creates a SHA256 hash of the input string
func hashString(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateXVerifyHeader generates the X-VERIFY header for PhonePe API requests
func GenerateXVerifyHeader(base64Payload string, endpoint string) string {
	saltKey := os.Getenv("PHONEPE_SALT")
	saltIndex := os.Getenv("PHONEPE_SALT_INDEX")
	hash := hashString(base64Payload + endpoint + saltKey)
	return fmt.Sprintf("%s###%s", hash, saltIndex)
}

// VerifyPhonePeSignature verifies the X-VERIFY header from PhonePe
func VerifyPhonePeSignature(payload []byte, xVerifyHeader string) bool {
	parts := strings.Split(xVerifyHeader, "###")
	if len(parts) != 2 {
		return false
	}
	receivedChecksum, saltIndex := parts[0], parts[1]
	saltKey := os.Getenv(fmt.Sprintf("PHONEPE_SALT_%s", saltIndex))
	if saltKey == "" {
		saltKey = os.Getenv("PHONEPE_SALT")
	}

	expectedChecksum := hashString(string(payload) + saltKey)
	return receivedChecksum == expectedChecksum
}

func GenerateTransactionID(username string, planType models.PlanType) string {
	return fmt.Sprintf("t_%s_%s", username, uuid.New().String()[:8])
}

// GetPhonePeBaseURL returns the appropriate base URL based on environment
func GetPhonePeBaseURL() string {
	if os.Getenv("GO_ENV") == "production" {
		return PhonePeProdBaseURL
	}
	return PhonePeProdBaseURL
}

// GenerateOASHash generates a hash for OAS verification
func GenerateOASHash(subscriptionID, mobile, gateway string) string {
	verifyStr := fmt.Sprintf("%s-%s-%s", subscriptionID, mobile, gateway)
	return hashString(verifyStr + os.Getenv("OAS_SALT"))
}

// VerifyOASSignature verifies an OAS signature
func VerifyOASSignature(subscriptionID, mobile, gateway, oasVerify string) bool {
	return oasVerify == GenerateOASHash(subscriptionID, mobile, gateway)
}
