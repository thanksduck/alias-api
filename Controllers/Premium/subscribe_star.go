package premium

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	models "github.com/thanksduck/alias-api/Models"
	"github.com/thanksduck/alias-api/utils"
)

/**JS
JavaScript implementation:

function getOASHash(subscriptionId, mobile, gateway) {
  const salt = process.env.OAS_SALT;
  const verifyStr = `${subscriptionId}-${mobile}-${gateway}`;
  const hash = crypto.createHash('sha256')
    .update(verifyStr + salt)
    .digest('hex');
  return hash;
}

function verifyOASSignature(subscriptionId, mobile, gateway, oasVerify) {
  const expectedHash = getOASHash(subscriptionId, mobile, gateway);
  return oasVerify.toLowerCase() === expectedHash.toLowerCase();
}
*/

func getOASHash(subscriptionID, mobile, gateway string) string {
	salt := os.Getenv("OAS_SALT")
	verifyStr := fmt.Sprintf("%s-%s-%s", subscriptionID, mobile, gateway)
	h := sha256.New()
	h.Write([]byte(verifyStr + salt))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyOASSignature(subscriptionID, mobile, gateway, oasVerify string) bool {
	expectedHash := getOASHash(subscriptionID, mobile, gateway)
	return strings.EqualFold(oasVerify, expectedHash)
}

func SubscribeToStar(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	fmt.Print(user)
	if !ok {
		utils.SendErrorResponse(w, "User not Found", http.StatusUnauthorized)
		return
	}

	var requestBody struct {
		SubscriptionID string          `json:"suid"`
		Mobile         string          `json:"mobile"`
		Gateway        string          `json:"gateway"`
		TransactionID  string          `json:"txnid"`
		Plan           models.PlanType `json:"plan"`
		MerchentUserID string          `json:"muid"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err)
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check required fields except MerchentUserID which is optional
	if requestBody.SubscriptionID == "" || requestBody.Mobile == "" ||
		requestBody.Gateway == "" || requestBody.TransactionID == "" {
		utils.SendErrorResponse(w, "Not enough information provided in request body", http.StatusBadRequest)
		return
	}

	// Verify OAS signature
	oasVerify := r.Header.Get("oas-verify")
	if oasVerify == "" {
		utils.SendErrorResponse(w, "Missing OAS verification header", http.StatusBadRequest)
		return
	}

	if !verifyOASSignature(requestBody.SubscriptionID, requestBody.Mobile, requestBody.Gateway, oasVerify) {
		utils.SendErrorResponse(w, "Invalid OAS verification", http.StatusUnauthorized)
		return
	}
}
