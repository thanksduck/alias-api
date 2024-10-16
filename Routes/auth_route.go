package routes

import (
	"net/http"

	auth "github.com/thanksduck/alias-api/Controllers/Auth"
)

func AuthRouter(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/v2/auth/signup", auth.Signup)
	mux.HandleFunc("POST /api/v2/auth/register", auth.Signup)
	mux.HandleFunc("POST /api/v2/auth/login", auth.Login)
	mux.HandleFunc("GET /api/v2/auth/google", auth.HandleGoogleLogin)
	mux.HandleFunc("GET /api/v2/auth/google/cb", auth.HandleGoogleCallback)
	mux.HandleFunc("POST /api/v2/auth/forget-password", auth.ForgetPassword)
	mux.HandleFunc("POST /api/v2/auth/reset-password/{token}", auth.ResetPassword)

}
