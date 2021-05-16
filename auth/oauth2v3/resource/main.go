package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	config = oauth2.Config{
		ClientID:     "12341234",
		ClientSecret: "12341234",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9096/credentials",
			TokenURL: "http://localhost:9096/token",
		},
	}
)

func main() {

	http.HandleFunc("/protected", validateToken(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm protected"))
	}))

	log.Fatal(http.ListenAndServe(":9097", nil))
}

func validateAccessToken(accessToken string) error {
	// GET 호출
	resp, err := http.Get("http://localhost:9096/validate?grant_type=client_credentials&access_token=" + accessToken)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var dat map[string]interface{}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	if dat["ACCESS_TOKEN"].(string) == accessToken {
		return nil
	}
	return errors.New("invalid token")
}

func validateToken(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.FormValue("access_token")
		err := validateAccessToken(val)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f.ServeHTTP(w, r)
	})
}
