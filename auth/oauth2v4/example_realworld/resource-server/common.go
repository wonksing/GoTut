package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-session/session"
)

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func genAccessTokenJWT(usrID string, tokenSecret string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["usr_id"] = usrID
	claims["exp"] = time.Now().Add(1 * time.Minute).Unix()

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type TokenClaim struct {
	UsrID string  `json:"usr_id"`
	Exp   float64 `json:"exp"`
}

// ValidateAccessToken 액세스토큰의 유효성을 판단하고 리프레시 가능한지 확인
func validateAccessToken(accessToken string, tokenSecret string) (*TokenClaim, bool, error) {
	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(tokenSecret), nil
	})

	var claim *TokenClaim
	if token != nil && token.Claims.Valid() != nil {
		c := token.Claims.(jwt.MapClaims)
		claim = &TokenClaim{
			UsrID: c["usr_id"].(string),
			Exp:   c["exp"].(float64),
		}
	}
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return claim, true, err
		}
		return nil, false, err
	}

	if token.Valid {
		c := token.Claims.(jwt.MapClaims)
		claim = &TokenClaim{
			UsrID: c["usr_id"].(string),
			Exp:   c["exp"].(float64),
		}
		return claim, false, nil
	}

	return nil, false, errors.New("Not Authorized")
}

type myCtx struct {
}

func AuthSessionIDHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")

		if token != "" {
			fmt.Println(token)
			claim, _, err := validateAccessToken(token, "asdfasdf1234")
			if err != nil {
				http.Error(w, "Not Authorized", http.StatusUnauthorized)
				return
			}

			store, err := session.Start(r.Context(), w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// session, err := store.Get(r, sessionName)
			// if err != nil {
			// 	log.WithError(err).Error("bad session")
			// 	http.SetCookie(w, &http.Cookie{Name: sessionName, MaxAge: -1, Path: "/"})
			// 	return
			// }

			// r = r.WithContext(context.WithValue(r.Context(), "session", session))
			// h(w, r)

			store.Set("userID", claim.UsrID)
			store.Save()

			ctx := r.Context()
			ctx = context.WithValue(ctx, myCtx{}, store)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// r.Header.Add("user_id", userID)
		// r.Header.Add("sid", sid)
		// r.AddCookie("Cookie", "go_session_id="+sid)
		// w.Header().Set("Cookie", "go_session_id="+sid)
		// w.Header().Set("user_id", userID)
		// w.Header().Add("redirect_uri", redirectURI)
		// w.Header().Add("code", code)

		// expiration := time.Now().Add(365 * 24 * time.Hour)
		// cookie := http.Cookie{Name: "go_session_id", Value: sid, Expires: expiration}
		// r.AddCookie(&cookie)

		next.ServeHTTP(w, r)
	}
}
