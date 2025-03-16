package paymentutils

import models "github.com/thanksduck/alias-api/Models"

// returns the monthly price for the given plan and months
// if the months are not 1, 3, 6, or 12, it returns 0
func GetMonthlyPrice(plan models.PlanType, months int) int {
	var basePrice int

	switch plan {
	case "star":
		basePrice = 49
	case "galaxy":
		basePrice = 79
	default:
		return 0
	}

	switch months {
	case 1:
		return basePrice * 2
	case 3:
		return basePrice + 30
	case 6:
		return basePrice + 20
	case 12:
		fallthrough
	default:
		if months >= 12 {
			return basePrice
		}
		return 0
	}
}
