package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	models "github.com/thanksduck/alias-api/Models"
	repository "github.com/thanksduck/alias-api/Repository"
	"github.com/thanksduck/alias-api/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func getGithubOAuthConfig() *oauth2.Config {
	clientid := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectUrl := os.Getenv("REDIRECT_HOST") + "/api/v2/auth/github/cb"

	conf := &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	return conf
}

func HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	conf := getGithubOAuthConfig()
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	conf := getGithubOAuthConfig()
	code := r.URL.Query().Get("code")

	t, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Error exchanging code:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := conf.Client(context.Background(), t)

	// Get user profile
	userResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer userResp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(userResp.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user emails
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer emailResp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add primary email to user data
	for _, email := range emails {
		if email["primary"].(bool) {
			userData["email"] = email["email"]
			break
		}
	}

	/*

		so what i will be getting for it ,
		userData["email"]
		userData["login"] // username
		userData["name"] // name
		userData["avatar_url"] // avatar
		userData["id"] // github id

	*/
	email, ok := userData["email"].(string)
	if !ok || email == "" {
		utils.SendErrorResponse(w, "Email not found or empty", http.StatusBadRequest)
		return
	}

	// Get github username
	githubUsername, ok := userData["login"].(string)
	if !ok || githubUsername == "" {
		utils.SendErrorResponse(w, "Github username not found or empty", http.StatusBadRequest)
		return
	}

	name, _ := userData["name"].(string)
	if name == "" {
		name, _ = userData["login"].(string)
	}
	if name == "" {
		nameParts := strings.Split(email, "@")
		if len(nameParts) > 0 {
			name = nameParts[0]
		}
	}

	// Check if user exists
	existingUser, err := repository.FindUserByUsernameOrEmail("", email)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error finding user:", err)
		utils.SendErrorResponse(w, "Error finding user", http.StatusInternalServerError)
		return
	}

	var user *models.User
	if existingUser != nil {
		user = existingUser
		updated := false

		if user.Provider != "google" && user.Provider != "github" {
			if !user.EmailVerified {
				user.EmailVerified = true
				updated = true
			}

			if user.Avatar == "" {
				user.Avatar = getStringValue(userData, "avatar_url", "")
				updated = true
			}

			user.Provider = "github"
			updated = true
		}

		if user.Name != name {
			user.Name = name
			updated = true
		}

		if updated {
			_, err = repository.UpdateUser(user.ID, user)
		}
	} else {
		username := githubUsername
		usernameExists, err := repository.FindUserByUsernameOrEmail(username, "")
		if err != nil && err != pgx.ErrNoRows {
			fmt.Println("Error finding user by username:", err)
			utils.SendErrorResponse(w, "Error finding user by username", http.StatusInternalServerError)
			return
		}
		if usernameExists != nil {
			username = username + "1"
		}

		user = &models.User{
			Email:         email,
			EmailVerified: true,
			Username:      username,
			Name:          name,
			Provider:      "github",
			Avatar:        getStringValue(userData, "avatar_url", ""),
			Password:      " ",
		}
		user, err = repository.CreateOrUpdateUser(user)
		if err != nil {
			fmt.Println("Error creating user:", err)
			utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
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
			Username: githubUsername,
			Github:   githubUsername,
		}
		_, err = repository.CreateOrUpdateSocialProfile(socialProfile)
		if err != nil {
			fmt.Println("Error creating social profile:", err)
			utils.SendErrorResponse(w, "Error creating social profile", http.StatusInternalServerError)
			return
		}
	}
	RedirectToFrontend(w, r, user)

}
