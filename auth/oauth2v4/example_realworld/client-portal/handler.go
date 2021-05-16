package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	u := config.AuthCodeURL("xyz",
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	log.Println(u)
	http.Redirect(w, r, u, http.StatusFound)
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.Form.Get("state")
	if state != "xyz" {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}
	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}
	token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	globalToken = token

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	if globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	globalToken.Expiry = time.Now()
	token, err := config.TokenSource(context.Background(), globalToken).Token()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	globalToken = token
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
	log.Println(token)
}

func tryHandler(w http.ResponseWriter, r *http.Request) {
	if globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, globalToken.AccessToken))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func pwdHandler(w http.ResponseWriter, r *http.Request) {
	token, err := config.PasswordCredentialsToken(context.Background(), "test", "test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	globalToken = token
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
	cfg := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     config.Endpoint.TokenURL,
	}

	token, err := cfg.Token(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(token)
}
