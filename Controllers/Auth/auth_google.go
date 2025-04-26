package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/thanksduck/alias-api/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig *oauth2.Config
	OauthStateString  = "random"
)

func init() {
	RdUrl := os.Getenv("REDIRECT_HOST") + "/api/v2/auth/google/cb"
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  RdUrl,
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if GoogleOauthConfig.ClientID == "" || GoogleOauthConfig.ClientSecret == "" {
		GoogleOauthConfig.RedirectURL = os.Getenv("REDIRECT_HOST") + "/api/v2/auth/google/cb"
		GoogleOauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
		GoogleOauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}
	Url := GoogleOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, Url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	content, err := GetUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println("Error getting user info:", err)
		utils.SendErrorResponse(w, "Error getting user info", http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(content, &userInfo); err != nil {
		fmt.Println("Error unmarshalling user info:", err)
		utils.SendErrorResponse(w, "Error unmarshalling user info", http.StatusInternalServerError)
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok || email == "" {
		fmt.Println("Error: Email not found or empty")
		utils.SendErrorResponse(w, "Email not found or empty", http.StatusBadRequest)
		return
	}

	// Determine user name
	name, _ := userInfo["name"].(string)
	if name == "" {
		givenName, _ := userInfo["given_name"].(string)
		familyName, _ := userInfo["family_name"].(string)
		if givenName != "" {
			name = givenName
			if familyName != "" {
				name += " " + familyName
			}
		}
	}
	if name == "" {
		name = strings.Split(email, "@")[0]
	}

	// Try to find existing user
	user, err := db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{
		Username: "", // ignored if empty
		Email:    email,
	})

	var isNewUser bool
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			isNewUser = true
		} else {
			fmt.Println("Error finding user:", err)
			utils.SendErrorResponse(w, "Error finding user", http.StatusInternalServerError)
			return
		}
	}

	// Variable to store the user object (new or existing)
	var userObj q.FindUserByUsernameOrEmailRow

	if isNewUser {
		username := strings.Split(email, "@")[0]
		if len(username) > 15 {
			username = username[:10]
		}

		// Check for username conflict
		if _, err := db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{
			Username: username,
			Email:    "",
		}); err == nil {
			username = username + "1"
		}

		// Create new user and capture the returned user
		now := time.Now()
		newUser, err := db.SQL.CreateOrUpdateUser(ctx, &q.CreateOrUpdateUserParams{
			Email:           email,
			Username:        username,
			Name:            name,
			IsEmailVerified: true,
			Provider:        "google",
			Avatar:          getStringValue(userInfo, "picture", ""),
			Password:        " ", // Google users don't use passwords
			CreatedAt:       now,
			UpdatedAt:       now,
		})
		if err != nil {
			fmt.Println("Error creating new user:", err)
			utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		userObj = q.FindUserByUsernameOrEmailRow{
			ID:                newUser.ID,
			Username:          newUser.Username,
			Name:              newUser.Name,
			Email:             newUser.Email,
			AliasCount:        newUser.AliasCount,
			DestinationCount:  newUser.DestinationCount,
			IsPremium:         newUser.IsPremium,
			Provider:          newUser.Provider,
			Avatar:            newUser.Avatar,
			PasswordChangedAt: newUser.PasswordChangedAt,
			IsActive:          newUser.IsActive,
			Password:          newUser.Password,
			IsEmailVerified:   newUser.IsEmailVerified,
			CreatedAt:         newUser.CreatedAt,
			UpdatedAt:         newUser.UpdatedAt,
		}
	} else {
		// Existing user - store in userObj
		userObj = *user
		updated := false

		if !user.IsEmailVerified {
			userObj.IsEmailVerified = true
			updated = true
		}

		if user.Avatar == "" {
			userObj.Avatar = getStringValue(userInfo, "picture", "")
			updated = true
		}

		if user.Provider != "google" {
			userObj.Provider = "google"
			updated = true
		}

		if updated {
			userObj.UpdatedAt = time.Now()
			updatedUser, err := db.SQL.CreateOrUpdateUser(ctx, &q.CreateOrUpdateUserParams{
				Email:           userObj.Email,
				Username:        userObj.Username,
				Name:            userObj.Name,
				IsEmailVerified: userObj.IsEmailVerified,
				Provider:        userObj.Provider,
				Avatar:          userObj.Avatar,
				Password:        userObj.Password,
				CreatedAt:       userObj.CreatedAt,
				UpdatedAt:       userObj.UpdatedAt,
			})
			if err != nil {
				fmt.Println("Error updating user:", err)
				utils.SendErrorResponse(w, "Error updating user", http.StatusInternalServerError)
				return
			}

			// Update userObj with the returned values
			userObj = q.FindUserByUsernameOrEmailRow{
				ID:                updatedUser.ID,
				Username:          updatedUser.Username,
				Name:              updatedUser.Name,
				Email:             updatedUser.Email,
				AliasCount:        updatedUser.AliasCount,
				DestinationCount:  updatedUser.DestinationCount,
				IsPremium:         updatedUser.IsPremium,
				Provider:          updatedUser.Provider,
				Avatar:            updatedUser.Avatar,
				PasswordChangedAt: updatedUser.PasswordChangedAt,
				IsActive:          updatedUser.IsActive,
				Password:          updatedUser.Password,
				IsEmailVerified:   updatedUser.IsEmailVerified,
				CreatedAt:         updatedUser.CreatedAt,
				UpdatedAt:         updatedUser.UpdatedAt,
			}
		}
	}

	// Now use userObj.ID to find or create social profile
	socialProfile, err := db.SQL.FindSocialProfileByUserID(ctx, userObj.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("Error fetching social profile:", err)
		utils.SendErrorResponse(w, "Error fetching social profile", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		// No existing social profile, create a new one
		_, err = db.SQL.CreateOrUpdateSocialProfile(ctx, &q.CreateOrUpdateSocialProfileParams{
			UserID:    userObj.ID,
			Username:  userObj.Username,
			Google:    getStringValue(userInfo, "sub", ""),
			Github:    "NULL",
			Facebook:  "NULL",
			CreatedAt: now,
			UpdatedAt: now,
		})
	} else {
		// Update existing social profile
		_, err = db.SQL.CreateOrUpdateSocialProfile(ctx, &q.CreateOrUpdateSocialProfileParams{
			UserID:    socialProfile.UserID,
			Username:  socialProfile.Username,
			Google:    getStringValue(userInfo, "sub", ""),
			Github:    socialProfile.Github,
			Facebook:  socialProfile.Facebook,
			CreatedAt: socialProfile.CreatedAt,
			UpdatedAt: now,
		})
	}
	if err != nil {
		fmt.Println("Error creating/updating social profile:", err)
		utils.SendErrorResponse(w, "Error creating/updating social profile", http.StatusInternalServerError)
		return
	}

	RedirectToFrontend(w, r, userObj.Username)
}

// Helper function to safely get string values from the map
func getStringValue(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key].(string); ok && val != "" {
		return val
	}
	return defaultValue
}

func GetUserInfo(state string, code string) ([]byte, error) {
	if state != OauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func RedirectToFrontend(w http.ResponseWriter, r *http.Request, username string) {
	// Create a token or session for the user here if needed
	token, err := utils.GenerateTempToken(username)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Encode the URL
	frontendURL := os.Getenv("FRONTEND_URL")
	redirectURL := fmt.Sprintf("%s?token=%s", frontendURL, url.QueryEscape(token))

	// Redirect to frontend with encoded URL
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	// utils.CreateSendResponse(w, user, redirectURL, http.StatusOK, "user", user.Username)
}
