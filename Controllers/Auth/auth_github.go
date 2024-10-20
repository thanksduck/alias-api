package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
		Scopes:       []string{"email", "profile"},
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
	// resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	var v any

	// Reading the JSON body using JSON decoder
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.CreateSendResponse(w, v, "Github Login Successful", http.StatusOK, "user", `1`)
}
