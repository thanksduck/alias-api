package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	GithubOauthConfig *oauth2.Config
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

	existingUser, err := repository.FindUserByUsernameOrEmail("", email)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error finding user:", err)
		utils.SendErrorResponse(w, "Error finding user", http.StatusInternalServerError)
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
		username := strings.Split(email, "@")[0]
		usernameExists, err := repository.FindUserByUsernameOrEmail(username, "")
		if err != nil && err != pgx.ErrNoRows {
			fmt.Println("Error finding user by username:", err)
			utils.SendErrorResponse(w, "Error finding user by username", http.StatusInternalServerError)
			return
		}
		if usernameExists != nil {
			username = username + "1"
		}
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
			Username:      username,
			EmailVerified: true,
			Provider:      "google",
			Avatar:        getStringValue(userInfo, "picture", ""),
			Password:      " ",
		}
		user, err = repository.CreateOrUpdateUser(user)
		if err != nil && err != pgx.ErrNoRows {
			fmt.Println("Error creating/updating user:", err)
			utils.SendErrorResponse(w, "Error creating/updating user", http.StatusInternalServerError)
			return
		}
	}

	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error updating/creating user:", err)
		utils.SendErrorResponse(w, "Error updating/creating user", http.StatusInternalServerError)
		return
	}

	socialProfile, err := repository.FindSocialProfileByIDOrUsername(user.ID, user.Username)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error finding social profile:", err)
		utils.SendErrorResponse(w, "Error finding social profile", http.StatusInternalServerError)
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
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error updating social profile:", err)
		utils.SendErrorResponse(w, "Error updating social profile", http.StatusInternalServerError)
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
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Encode the URL
	frontendURL := os.Getenv("FRONTEND_URL")
	redirectURL := fmt.Sprintf("%s?token=%s", frontendURL, url.QueryEscape(token))

	// Redirect to frontend with encoded URL
	// http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	utils.CreateSendResponse(w, user, redirectURL, http.StatusOK, "user", user.ID)
}
