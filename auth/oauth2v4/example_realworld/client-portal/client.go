package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	authServerURL = "http://localhost:9096"
)

var (
	config = oauth2.Config{
		ClientID:     "12345",
		ClientSecret: "12345678",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func requestAuthorize() {
	// GET 호출
	resp, err := http.Get("http://localhost:9096/token?grant_type=client_credentials&client_id=12341234&client_secret=12341234&scope=all")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(data))
}

func main() {

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/oauth2", oauthHandler)

	http.HandleFunc("/refresh", refreshHandler)

	http.HandleFunc("/try", tryHandler)

	http.HandleFunc("/pwd", pwdHandler)

	http.HandleFunc("/client", clientHandler)

	log.Println("Client is running at 9094 port.Please open http://localhost:9094")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
