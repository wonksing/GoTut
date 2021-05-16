package commonutil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v3/generates"
)

func OutputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func DumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}

func VerifyJWT(secret string, tokenStr string) (string, error) {
	// Parse and verify jwt access token
	token, err := jwt.ParseWithClaims(tokenStr, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse error")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*generates.JWTAccessClaims)
	if !ok || !token.Valid {
		// panic("invalid token")
		return "", errors.New("invalid token")
	}

	fmt.Println("claims:", claims.Audience, claims.Id, claims.Subject)
	return claims.Audience, nil
}

func SetCookie(w http.ResponseWriter, name, value string, validHours time.Duration) {

	ck := http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		Expires:  time.Now().Add(validHours * time.Hour),
	}
	http.SetCookie(w, &ck)
}
