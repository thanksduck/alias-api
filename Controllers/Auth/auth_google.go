package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  "/api/v2/auth/google/cb",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {

	if GoogleOauthConfig.ClientID == "" || GoogleOauthConfig.ClientSecret == "" {
		GoogleOauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
		GoogleOauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}
	url := GoogleOauthConfig.AuthCodeURL(OauthStateString)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := GetUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/health", http.StatusTemporaryRedirect)
		return
	}

	var userInfo map[string]interface{}
	json.Unmarshal(content, &userInfo)

	// Create or update user
	user := models.User{
		Email:         userInfo["email"].(string),
		Name:          userInfo["name"].(string),
		EmailVerified: true,
		Provider:      "google",
		Avatar:        userInfo["picture"].(string),
	}

	// Assuming you have a function to create or update the user
	createdUser, err := repository.CreateOrUpdateUser(&user)
	if err != nil {
		fmt.Println("Error creating/updating user:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Create or update social profile
	socialProfile := models.SocialProfile{
		UserID:   createdUser.ID,
		Username: createdUser.Username,
		Google:   userInfo["sub"].(string), // Using 'sub' as unique Google identifier
	}

	// Assuming you have a function to create or update the social profile
	_, err = repository.CreateOrUpdateSocialProfile(&socialProfile)
	if err != nil {
		fmt.Println("Error creating/updating social profile:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Redirect to the final destination
	// http.Redirect(w, r, "/send/to", http.StatusTemporaryRedirect)
	utils.CreateSendResponse(w, createdUser, "Login Successful", http.StatusOK, "user", createdUser.ID)
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
