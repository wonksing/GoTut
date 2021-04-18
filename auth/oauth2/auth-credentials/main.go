package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/oauth2.v3/models"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()

	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	http.HandleFunc("/credentials", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("login_id")
		pw := r.FormValue("login_pw")

		if login(id, pw) == false {
			http.Error(w, "invalid", http.StatusInternalServerError)
			return
		}
		// clientId := uuid.New().String()[:8]
		// clientSecret := uuid.New().String()[:8]
		clientId := uuid.New().String()
		clientSecret := uuid.New().String()
		err := clientStore.Set(clientId, &models.Client{
			ID:     clientId,
			Secret: clientSecret,
			Domain: "http://localhost:9094",
			UserID: id,
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"CLIENT_ID": clientId, "CLIENT_SECRET": clientSecret})
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		// typ, gr, err := srv.ValidationTokenRequest(r)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// fmt.Println(typ, gr)

		ti, err := srv.ValidationBearerToken(r)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"CLIENT_ID": ti.GetClientID(), "ACCESS_TOKEN": ti.GetAccess()})
	})

	log.Fatal(http.ListenAndServe(":9096", nil))
}

func login(id, pw string) bool {
	if id == "admin" && pw == "admin" {
		return true
	}
	return false
}
