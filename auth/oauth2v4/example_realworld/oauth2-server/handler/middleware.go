package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/wonksing/gotut/auth/oauth2v4/example_realworld/oauth2-server/commonutil"
)

func handleJWTAuth(w http.ResponseWriter, r *http.Request, secretKey string) (*commonutil.TokenClaim, error) {

	ck, err := r.Cookie("access_token")
	if err != nil || ck.Value == "" {
		if r.Form == nil {
			r.ParseForm()
		}

		commonutil.SetCookie(w, "oauth_return_uri", r.Form.Encode(), time.Duration(24*365))

		w.Header().Set("Location", "/oauth/login")
		w.WriteHeader(http.StatusFound)
		return nil, errors.New("no valid access_token")
		// http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	token := ck.Value
	claim, _, err := commonutil.ValidateAccessToken(token, secretKey)

	if err != nil || claim == nil {
		commonutil.SetCookie(w, "access_token", "", time.Duration(24*365))
		w.Header().Set("Location", "/oauth/login")
		w.WriteHeader(http.StatusFound)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))

		// if expired {
		// 	// refresh logic...
		// }
		// http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// return errors.New(http.StatusText(http.StatusUnauthorized))
	}

	// oauth authorize 핸들러에서
	// commonutil.SetCookie(w, "oauth_return_uri", r.Form.Encode(), time.Duration(24*365))

	// w.Header().Set("usrID", claim.UsrID)
	return claim, nil
}

// AuthJWTHandler to verify the request
func AuthJWTHandler(next http.HandlerFunc, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		claim, err := handleJWTAuth(w, r, secret)
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), commonutil.TokenClaim{}, claim)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
