package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/spf13/viper"
	"github.com/wonksing/gotut/auth/oauth2v4/example_realworld/oauth2-server/handler"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var (
	dumpvar bool
	// idvar     string
	// secretvar string
	// domainvar string
	// portvar   int

	addr           string
	configFileName string
)

func init() {
	flag.StringVar(&addr, "addr", ":9096", "listening address(eg. :9096)")
	flag.StringVar(&configFileName, "cn", "./configs/server.yml", "config file name")
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	// flag.StringVar(&idvar, "i", "12345", "The client id being passed in")
	// flag.StringVar(&secretvar, "s", "12345678", "The client secret being passed in")
	// flag.StringVar(&domainvar, "r", "http://localhost:9094", "The domain of the redirect url")
	// flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {
	flag.Parse()
	if dumpvar {
		log.Println("Dumping requests")
	}

	conf := viper.New()
	conf.SetConfigFile(configFileName)
	err := conf.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	jwtAccessToken := conf.GetBool("token_config.jwt_access_token")
	jwtSecret := conf.GetString("token_config.jwt_secret")
	authCodeAccessTokenExp := conf.GetInt("token_config.auth_code.access_token_exp")
	authCodeRefreshTokenExp := conf.GetInt("token_config.auth_code.refresh_token_exp")
	authCodeGenerateRefresh := conf.GetBool("token_config.auth_code.generate_refresh")
	clientCredentialsAccessTokenExp := conf.GetInt("token_config.client_credential.access_token_exp")
	clientCredentialsRefreshTokenExp := conf.GetInt("token_config.client_credential.refresh_token_exp")
	clientCredentialsGenerateRefresh := conf.GetBool("token_config.client_credential.generate_refresh")

	tokenStoreFilePath := conf.GetString("token_store.file.path")

	clientStore := store.NewClientStore()
	cc := conf.Sub("client_credentials")
	ccSettings := cc.AllSettings()
	for key, val := range ccSettings {
		v := val.(map[string]interface{})
		id := v["id"].(string)
		secret := v["secret"].(string)
		domain := v["domain"].(string)
		fmt.Println(key)
		clientStore.Set(id, &models.Client{
			ID:     id,
			Secret: secret,
			Domain: domain,
		})
	}

	manager := manage.NewDefaultManager()
	// manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * time.Duration(authCodeAccessTokenExp),
		RefreshTokenExp:   time.Hour * time.Duration(authCodeRefreshTokenExp),
		IsGenerateRefresh: authCodeGenerateRefresh,
	})
	manager.SetClientTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * time.Duration(clientCredentialsAccessTokenExp),
		RefreshTokenExp:   time.Hour * time.Duration(clientCredentialsRefreshTokenExp),
		IsGenerateRefresh: clientCredentialsGenerateRefresh,
	})

	// token store
	if tokenStoreFilePath != "" {
		manager.MustTokenStorage(store.NewFileTokenStore(tokenStoreFilePath))
	} else {
		// memory
		manager.MustTokenStorage(store.NewMemoryTokenStore())
	}

	// generate jwt access token
	if jwtAccessToken {
		manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(jwtSecret), jwt.SigningMethodHS512))
	} else {
		manager.MapAccessGenerate(generates.NewAccessGenerate())
	}

	manager.MapClientStorage(clientStore)

	// srv := server.NewServer(server.NewConfig(), manager)
	srv := server.NewDefaultServer(manager)
	srv.Config.AllowGetAccessRequest = true

	// Password credentials
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	// Authorization Code Grant
	srv.SetUserAuthorizationHandler(handler.UserAuthorizeHandler)
	// srv.SetUserAuthorizationHandler(userAuthorizeHandler2)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	h := handler.ServerHandler{
		Srv:         srv,
		JwtSecret:   jwtSecret,
		ClientStore: clientStore,
	}

	http.HandleFunc("/", handler.AuthJWTHandler(h.HelloHandler, jwtSecret, "/login"))
	http.HandleFunc("/hello", handler.AuthJWTHandler(h.HelloHandler, jwtSecret, "/login"))
	http.HandleFunc("/login", h.LoginHandler)

	http.HandleFunc("/oauth/login", h.OAuthLoginHandler)
	http.HandleFunc("/auth", handler.AuthJWTHandler(h.AuthHandler, jwtSecret, "/oauth/login"))

	http.HandleFunc("/oauth/authorize", handler.AuthJWTHandler(h.OAuthAuthHandler, jwtSecret, "/oauth/login"))

	// token request for all types of grant
	// Client Credentials Grant comes here directly
	// http.HandleFunc("/oauth/token", handler.AuthJWTHandler(h.OAuthTokenHandler, jwtSecret))
	http.HandleFunc("/oauth/token", h.OAuthTokenHandler)

	// validate access token
	http.HandleFunc("/test", h.OAuthTestHandler)

	// client credential 저장
	http.HandleFunc("/credentials", h.CredentialHandler)

	log.Printf("Server is running at %v.\n", addr)
	log.Printf("Point your OAuth client Auth endpoint to %s%s", "http://"+addr, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s%s", "http://"+addr, "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v", addr), nil))
}
