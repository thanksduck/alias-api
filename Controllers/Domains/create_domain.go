package domains

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	middlewares "github.com/thanksduck/alias-api/Middlewares"
	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
	requests "github.com/thanksduck/alias-api/Requests"
	"github.com/thanksduck/alias-api/cfconfig"
	"github.com/thanksduck/alias-api/utils"
)

func CreateCustomDomain(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	}
	if !user.EmailVerified {
		utils.SendErrorResponse(w, "Please Verify Your Email to add Domain", http.StatusUnauthorized)
		return
	}
	if !user.IsPremium {
		utils.SendErrorResponse(w, "You have to be a premium user to add a custom domain", http.StatusPaymentRequired)
		return
	}

	var requestBody struct {
		Domain    string `json:"domain"`
		Subdomain string `json:"subdomain"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Domain == "" || requestBody.Subdomain == "" {
		utils.SendErrorResponse(w, "Domain is required", http.StatusBadRequest)
		return
	}
	domain := strings.ToLower(requestBody.Domain)
	subdomain := strings.ToLower(requestBody.Subdomain)
	if !middlewares.ValidBody.IsAllowedDomain(domain) || !middlewares.ValidBody.IsValidDomain(subdomain) {
		utils.SendErrorResponse(w, "Domain Cant be Proccessed", http.StatusUnprocessableEntity)
		return
	}

	config, ok := cfconfig.SelectDomain(domain)
	if !ok {
		utils.SendErrorResponse(w, "No configuration found for domain", http.StatusNotFound)
		return
	}

	records, err := requests.CreateNewDomain(subdomain, domain, config)
	if err != nil {
		utils.SendErrorResponse(w, "Error creating domain", http.StatusInternalServerError)
		return
	}
	customDomain := models.CustomDomain{
		UserID:   user.ID,
		Domain:   subdomain + "." + domain,
		Username: user.Username,
	}
	domainId, err := repository.CreateNewSubDomain(&customDomain)
	if err != nil {
		utils.SendErrorResponse(w, "Error creating domain", http.StatusInternalServerError)
		return
	}
	for _, record := range records {
		dnsRecord := models.CustomDomainDNSRecord{
			CustomDomainID: domainId,
			CloudflareID:   record.ID,
			Type:           record.Type,
			Name:           record.Name,
			Content:        record.Content,
			TTL:            record.TTL,
			Priority:       record.Priority,
		}
		err = repository.CreateNewDNSRecord(&dnsRecord)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Error creating DNS record for domain %s yes\n", customDomain.Domain)
			utils.SendErrorResponse(w, "Error creating domain", http.StatusInternalServerError)
			return
		}
	}
	utils.CreateSendResponse(w, customDomain, "Domain Created Successfully", http.StatusCreated, "domain", user.ID)

}
