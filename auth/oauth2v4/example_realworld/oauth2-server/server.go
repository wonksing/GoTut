package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/wonksing/gotut/auth/oauth2v4/example_realworld/oauth2-server/commonutil"
	"github.com/wonksing/gotut/auth/oauth2v4/example_realworld/oauth2-server/handler"

	"github.com/go-oauth2/oauth2/v4/models"
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
	flag.StringVar(&configFileName, "conf", "./configs/server.yml", "config file name")
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

	cc := conf.Sub("client_credentials")
	ccSettings := cc.AllSettings()

	oauthServer := commonutil.NewOAuthServer(authCodeAccessTokenExp, authCodeRefreshTokenExp, authCodeGenerateRefresh,
		clientCredentialsAccessTokenExp, clientCredentialsRefreshTokenExp, clientCredentialsGenerateRefresh,
		tokenStoreFilePath, jwtAccessToken, jwtSecret)
	for _, val := range ccSettings {
		v := val.(map[string]interface{})
		id := v["id"].(string)
		secret := v["secret"].(string)
		domain := v["domain"].(string)
		oauthServer.ClientStore.Set(id, &models.Client{
			ID:     id,
			Secret: secret,
			Domain: domain,
		})
	}

	// Password credentials
	oauthServer.Srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	// Authorization Code Grant
	oauthServer.Srv.SetUserAuthorizationHandler(handler.UserAuthorizeHandler)

	h := handler.ServerHandler{
		Srv:         oauthServer.Srv,
		JwtSecret:   jwtSecret,
		ClientStore: oauthServer.ClientStore,
	}

	// 테스트용 API
	http.HandleFunc(handler.API_INDEX, handler.AuthJWTHandler(h.HelloHandler, jwtSecret, handler.API_LOGIN))
	http.HandleFunc(handler.API_HELLO, handler.AuthJWTHandler(h.HelloHandler, jwtSecret, handler.API_LOGIN))
	http.HandleFunc(handler.API_LOGIN, h.LoginHandler)

	// OAuth2 API
	// 리소스 서버에 인증
	http.HandleFunc(handler.API_OAUTH_LOGIN, h.OAuthLoginHandler)
	// 리소스 서버의 정보 인가
	http.HandleFunc(handler.API_OAUTH_ALLOW, handler.AuthJWTHandler(h.OAuthAllowAuthorizationHandler, jwtSecret, handler.API_OAUTH_LOGIN))
	// Authorization Code Grant Type
	http.HandleFunc(handler.API_OAUTH_AUTHORIZE, handler.AuthJWTHandler(h.OAuthAuthorizeHandler, jwtSecret, handler.API_OAUTH_LOGIN))

	// token request for all types of grant
	// Client Credentials Grant comes here directly
	// Client Server용 API
	http.HandleFunc(handler.API_OAUTH_TOKEN, h.OAuthTokenHandler)

	// validate access token
	http.HandleFunc(handler.API_OAUTH_TOKEN_VALIDATE, h.OAuthValidateTokenHandler)

	// client credential 저장
	http.HandleFunc(handler.API_OAUTH_CREDENTIALS, h.CredentialHandler)

	log.Printf("Server is running at %v.\n", addr)
	log.Printf("Point your OAuth client Auth endpoint to %s%s", "http://"+addr, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s%s", "http://"+addr, "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v", addr), nil))
}
