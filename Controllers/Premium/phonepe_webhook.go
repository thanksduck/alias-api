package premium

import (
	"fmt"
	"net/http"
)

func PhonePeWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PhonePe Webhook")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Thank you so much")
}
