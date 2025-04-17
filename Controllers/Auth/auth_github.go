package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/thanksduck/alias-api/Database"
	q "github.com/thanksduck/alias-api/internal/db"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/thanksduck/alias-api/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func getGithubOAuthConfig() *oauth2.Config {
	client := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectUrl := os.Getenv("REDIRECT_HOST") + "/api/v2/auth/github/cb"

	conf := &oauth2.Config{
		ClientID:     client,
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
	ctx := r.Context()

	conf := getGithubOAuthConfig()
	code := r.URL.Query().Get("code")

	t, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Error exchanging code:", err)
		utils.SendErrorResponse(w, "Error exchanging code", http.StatusBadRequest)
		return
	}

	client := conf.Client(context.Background(), t)

	// Get user profile
	userResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		fmt.Println("Error getting user profile:", err)
		utils.SendErrorResponse(w, "Error getting user profile", http.StatusBadRequest)
		return
	}
	defer userResp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(userResp.Body).Decode(&userData); err != nil {
		fmt.Println("Error decoding user profile:", err)
		utils.SendErrorResponse(w, "Error decoding user profile", http.StatusInternalServerError)
		return
	}

	// Get user emails
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		fmt.Println("Error getting user emails:", err)
		utils.SendErrorResponse(w, "Error getting user emails", http.StatusBadRequest)
		return
	}
	defer emailResp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		fmt.Println("Error decoding user emails:", err)
		utils.SendErrorResponse(w, "Error decoding user emails", http.StatusInternalServerError)
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
		Available GitHub data:
		userData["email"]       // primary email
		userData["login"]       // username
		userData["name"]        // name
		userData["avatar_url"]  // avatar
		userData["id"]          // github id
	*/

	email, ok := userData["email"].(string)
	if !ok || email == "" {
		fmt.Println("Error: Email not found or empty")
		utils.SendErrorResponse(w, "Email not found or empty", http.StatusBadRequest)
		return
	}

	// Get gitHub username
	githubUsername, ok := userData["login"].(string)
	if !ok || githubUsername == "" {
		fmt.Println("Error: Github username not found or empty")
		utils.SendErrorResponse(w, "Github username not found or empty", http.StatusBadRequest)
		return
	}

	// Determine user name
	name, _ := userData["name"].(string)
	if name == "" {
		name = githubUsername
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
		username = githubUsername
		if len(username) > 15 {
			username = username[:10]
		}

		// Check for username conflict
		existingUser, err := db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{
			Username: username,
			Email:    "",
		})
		if err == nil && existingUser != nil {
			username = username + "1"
		}

		// Create new user
		now := time.Now()
		_, err = db.SQL.CreateOrUpdateUser(ctx, &q.CreateOrUpdateUserParams{
			Email:           email,
			Username:        username,
			Name:            name,
			IsEmailVerified: true,
			Provider:        "github",
			Avatar:          getStringValue(userData, "avatar_url", ""),
			Password:        " ", // GitHub users don't use passwords
			CreatedAt:       now,
			UpdatedAt:       now,
		})
		if err != nil {
			fmt.Println("Error creating new user:", err)
			utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Fetch the newly created user to get the ID
		user, err = db.SQL.FindUserByUsernameOrEmail(ctx, &q.FindUserByUsernameOrEmailParams{
			Email: email,
		})
		if err != nil {
			fmt.Println("Error fetching new user:", err)
			utils.SendErrorResponse(w, "Error fetching new user", http.StatusInternalServerError)
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
			user.Avatar = getStringValue(userData, "avatar_url", "")
			updated = true
		}

		if user.Provider != "github" && user.Provider != "google" {
			user.Provider = "github"
			updated = true
		}

		if user.Name != name {
			user.Name = name
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
			Google:    "NULL",
			Github:    githubUsername,
			Facebook:  "NULL",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	} else {
		_, err = db.SQL.CreateOrUpdateSocialProfile(ctx, &q.CreateOrUpdateSocialProfileParams{
			UserID:    socialProfile.UserID,
			Username:  socialProfile.Username,
			Google:    socialProfile.Google,
			Github:    githubUsername,
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
