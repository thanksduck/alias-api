package middlewares

import (
	"errors"
	db "github.com/thanksduck/alias-api/Database"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/thanksduck/alias-api/utils"
)

func Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var token string

		// Check for token in Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Check for token in cookies
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			utils.SendErrorResponse(w, "You are not logged in! Please login to get access.", http.StatusUnauthorized)
			return
		}

		// Verify the token
		claims, err := utils.VerifyToken(token)
		if err != nil {
			utils.SendErrorResponse(w, "Your Login Session has expired. Please login again.", http.StatusUnauthorized)
			return
		}

		// Get the user from the database
		username := claims.Username
		user, err := db.SQL.FindUserByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				utils.SendErrorResponse(w, "User not found. Please login again.", http.StatusUnauthorized)
			} else {
				utils.SendErrorResponse(w, "An error occurred. Please try again.", http.StatusInternalServerError)
			}
			return
		}

		passwordChangedAtUnix := user.PasswordChangedAt.Unix()
		if passwordChangedAtUnix > claims.IssuedAt.Unix() {
			utils.SendErrorResponse(w, "User recently changed password! Please login again.", http.StatusUnauthorized)
			return
		}

		// Set the user in the request context
		ctx = utils.SetUserInContext(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
