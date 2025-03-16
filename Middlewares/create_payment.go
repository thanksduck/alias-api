package middlewares

import (
	"fmt"
	"strings"

	models "github.com/thanksduck/alias-api/Models"
)

func ValidatePaymentBody(requestBody models.PaymentRequest) error {
	var missingFields []string
	if requestBody.Plan == "" {
		missingFields = append(missingFields, "plan")
	}
	if requestBody.Months == 0 {
		missingFields = append(missingFields, "months")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	if requestBody.Plan != "star" && requestBody.Plan != "galaxy" {
		return fmt.Errorf("invalid plan")
	}

	if requestBody.Months != 1 && requestBody.Months != 3 && requestBody.Months != 6 && requestBody.Months != 12 {
		return fmt.Errorf("invalid months")
	}

	return nil
}
