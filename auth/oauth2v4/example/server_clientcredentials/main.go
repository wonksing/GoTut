package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}
func createAndInitClientStore() *store.ClientStore {
	clientStore := store.NewClientStore()

	clientID := "b16bb655"
	clientSecret := "b16bb655-9568-4faa-82c0-4d152ed33035"
	domain := "http://localhost:9094"
	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: domain,
	})

	clientID = "222222"
	clientSecret = "22222222"
	domain = "http://localhost:9094"
	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: domain,
	})
	return clientStore
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	// 로그인

	// 권한 허용

	userID = "wonk"
	return
}

// ClientFormHandler get client data from form
func clientFormHandler(r *http.Request) (string, string, error) {
	clientID := r.Form.Get("client_id")
	if clientID == "" {
		return "", "", errors.ErrInvalidClient
	}
	clientSecret := r.Form.Get("client_secret")
	return clientID, clientSecret, nil
}

func main() {
	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetClientTokenCfg(&manage.Config{
		AccessTokenExp:    time.Second * 30,
		RefreshTokenExp:   time.Hour * 120,
		IsGenerateRefresh: true,
	})
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	// manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
	// 	IsGenerateRefresh:  true,
	// 	IsRemoveAccess:     true,
	// 	IsRemoveRefreshing: true,
	// })

	// token memory store
	// manager.MustTokenStorage(store.NewMemoryTokenStore())

	// token file store
	manager.MustTokenStorage(store.NewFileTokenStore("token.store"))

	// client memory store
	clientStore := createAndInitClientStore()
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	// srv.SetClientInfoHandler(clientFormHandler) // 이걸 사용하지 않으면 헤더의 "Authorization: Basic" 키를 이용한다(base64 of client_id:client_secret)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		// http://localhost:9096/oauth/authorize?client_id=222222&code_challenge=Qn3Kywp0OiU4NK_AFzGPlmrcYJDJ13Abj_jdL08Ahg8%3D&code_challenge_method=S256&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2&response_type=code&scope=all&state=xyz

		_ = dumpRequest(os.Stdout, "authorize", r) // Ignore the error
		// r.Form.Set("client_id", r.FormValue("client_id"))
		// r.Form.Set("code_challenge", r.FormValue("code_challenge"))
		// r.Form.Set("code_challenge_method", r.FormValue("code_challenge_method"))
		// r.Form.Set("redirect_uri", r.FormValue("redirect_uri"))
		// r.Form.Set("response_type", r.FormValue("response_type"))
		// r.Form.Set("scope", r.FormValue("scope"))
		// r.Form.Set("state", r.FormValue("state"))
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	// access token을 요청하거나 access_token을 refresh
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		// authorization code
		// code=MZZKNDE5NDYTYTMWNY0ZZTEXLWE3YJUTMZE5YJRJOGZLODEX&code_verifier=s256example&grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2

		_ = dumpRequest(os.Stdout, "token", r) // Ignore the error
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	// access token을 요청하거나 access_token을 refresh
	http.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		// authorization code
		// code=MZZKNDE5NDYTYTMWNY0ZZTEXLWE3YJUTMZE5YJRJOGZLODEX&code_verifier=s256example&grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2

		_ = dumpRequest(os.Stdout, "token", r) // Ignore the error

		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		ti, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID := ""
		domain := ""
		ci, err := clientStore.GetByID(r.Context(), ti.GetClientID())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			userID = ci.GetUserID()
			domain = ci.GetDomain()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"client_id":     ti.GetClientID(),
			"access_token":  ti.GetAccess(),
			"refresh_token": ti.GetRefresh(),
			"user_id":       userID,
			"redirect_uri":  ti.GetRedirectURI(),
			"scope":         ti.GetScope(),
			"domain":        domain,
		})
	})

	log.Fatal(http.ListenAndServe(":9096", nil))
}
