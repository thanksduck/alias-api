package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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

	var username string
	if isNewUser {
		username = strings.Split(email, "@")[0]
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

		// Create new user using `CreateOrUpdateUser`
		now := time.Now()
		_, err = db.SQL.CreateOrUpdateUser(ctx, &q.CreateOrUpdateUserParams{
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
	} else {
		// Existing user update (if needed)
		updated := false

		if !user.IsEmailVerified {
			user.IsEmailVerified = true
			updated = true
		}

		if user.Avatar == "" {
			user.Avatar = getStringValue(userInfo, "picture", "")
			updated = true
		}

		if user.Provider != "google" {
			user.Provider = "google"
			updated = true
		}

		if updated {
			user.UpdatedAt = time.Now()
			_, err = db.SQL.CreateOrUpdateUser(ctx, &q.CreateOrUpdateUserParams{
				Email:           user.Email,
				Username:        user.Username,
				Name:            user.Name,
				IsEmailVerified: user.IsEmailVerified,
				Provider:        user.Provider,
				Avatar:          user.Avatar,
				Password:        user.Password,
				CreatedAt:       user.CreatedAt,
				UpdatedAt:       user.UpdatedAt,
			})
			if err != nil {
				fmt.Println("Error updating user:", err)
				utils.SendErrorResponse(w, "Error updating user", http.StatusInternalServerError)
				return
			}
		}
	}

	// Upsert Social Profile
	socialProfile, err := db.SQL.FindSocialProfileByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("Error fetching social profile:", err)
		utils.SendErrorResponse(w, "Error fetching social profile", http.StatusInternalServerError)
		return
	}

	if socialProfile == nil {
		_, err = db.SQL.CreateOrUpdateSocialProfile(ctx, &q.CreateOrUpdateSocialProfileParams{
			UserID:    user.ID,
			Username:  user.Username,
			Google:    getStringValue(userInfo, "sub", ""),
			Github:    "NULL",
			Facebook:  "NULL",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	} else {
		_, err = db.SQL.CreateOrUpdateSocialProfile(ctx, &q.CreateOrUpdateSocialProfileParams{
			UserID:    socialProfile.UserID,
			Username:  socialProfile.Username,
			Google:    getStringValue(userInfo, "sub", ""),
			Github:    socialProfile.Github,
			Facebook:  socialProfile.Facebook,
			CreatedAt: socialProfile.CreatedAt,
			UpdatedAt: time.Now(),
		})
	}
	if err != nil {
		fmt.Println("Error creating/updating social profile:", err)
		utils.SendErrorResponse(w, "Error creating/updating social profile", http.StatusInternalServerError)
		return
	}

	RedirectToFrontend(w, r, user.Username)
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
