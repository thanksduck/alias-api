package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
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
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if GoogleOauthConfig.ClientID == "" || GoogleOauthConfig.ClientSecret == "" {
		GoogleOauthConfig.RedirectURL = os.Getenv("REDIRECT_HOST") + "/api/v2/auth/google/cb"
		GoogleOauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
		GoogleOauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}
	url := GoogleOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := GetUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println("Error getting user info:", err)
		http.Redirect(w, r, "/health", http.StatusTemporaryRedirect)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(content, &userInfo); err != nil {
		fmt.Println("Error unmarshalling user info:", err)
		http.Redirect(w, r, "/health", http.StatusTemporaryRedirect)
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok || email == "" {
		fmt.Println("Error: Email not found or empty")
		http.Redirect(w, r, "/health", http.StatusTemporaryRedirect)
		return
	}

	existingUser, err := repository.FindUserByUsernameOrEmail("", email)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error finding user:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var user *models.User
	if existingUser != nil {
		user = existingUser
		user.EmailVerified = true
		user.Avatar = getStringValue(userInfo, "picture", "")
		if user.Provider != "google" {
			user.Provider = "google"
		}
		_, err = repository.UpdateUser(user.ID, user)
	} else {
		name := getStringValue(userInfo, "name", "Google User")
		if name == "Google User" {
			nameParts := strings.Split(email, "@")
			if len(nameParts) > 0 {
				name = nameParts[0]
			}
		}
		user = &models.User{
			Email:         email,
			Name:          name,
			EmailVerified: true,
			Provider:      "google",
			Avatar:        getStringValue(userInfo, "picture", ""),
		}
		user, err = repository.CreateOrUpdateUser(user)
	}

	if err != nil {
		fmt.Println("Error updating/creating user:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	socialProfile, err := repository.FindSocialProfileByIDOrUsername(user.ID, user.Username)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error finding social profile:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if socialProfile == nil {
		socialProfile = &models.SocialProfile{
			UserID:   user.ID,
			Username: user.Username,
		}
	}
	socialProfile.Google = getStringValue(userInfo, "sub", "")

	_, err = repository.CreateOrUpdateSocialProfile(socialProfile)
	if err != nil {
		fmt.Println("Error updating social profile:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	RedirectToFrontend(w, r, user)
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

func RedirectToFrontend(w http.ResponseWriter, r *http.Request, user *models.User) {
	// Create a token or session for the user here if needed
	// For example:
	// token, err := utils.CreateJWTToken(user)
	// if err != nil {
	//     http.Error(w, "Error creating token", http.StatusInternalServerError)
	//     return
	// }

	// Redirect to frontend with token or user info
	// frontendURL := os.Getenv("FRONTEND_URL")
	// You might want to append user info or token to the URL
	// frontendURL += "?token=" + token

	utils.CreateSendResponse(w, user, "Login Successful", http.StatusOK, "user", user.ID)
	// http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}
