package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/go-session/session"
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
	flag.StringVar(&addr, "addr", ":9099", "listening address(eg. :9099)")
	flag.StringVar(&configFileName, "cn", "./configs/server.yml", "config file name")
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
}

func main() {
	flag.Parse()
	if dumpvar {
		log.Println("Dumping requests")
	}

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/oauth/login", loginHandler)

	http.HandleFunc("/protected", protectedHandler)

	http.HandleFunc("/logout", logoutHandler)

	tmpAddr := addr
	if strings.HasPrefix(addr, ":") {
		tmpAddr = "localhost" + addr
	}
	log.Printf("Resource is running at %s. Please open http://%s", addr, tmpAddr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "indexHandler", r) // Ignore the error
	}

	outputHTML(w, r, "static/index.html")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "loginHandler", r) // Ignore the error
	}

	redirectURI := r.FormValue("redirect_uri")

	if redirectURI != "" {
		store, err := session.Start(context.Background(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(redirectURI)

		store.Set("redirect_uri", redirectURI)
		store.Save()
	}

	if r.Method == "GET" {
		outputHTML(w, r, "static/login.html")
		return
	}

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp, ok := store.Get("redirect_uri")
	if !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	userID := r.Form.Get("user_id")
	userPW := r.Form.Get("user_pw")
	if len(userID) <= 0 {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if userID != userPW {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	store.Set("userID", userID)
	store.Save()

	w.Header().Set("Location", tmp.(string))
	w.WriteHeader(http.StatusFound)
	return

}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "protectedHandler", r) // Ignore the error
	}
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, ok := store.Get("userID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/protected.html")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "logoutHandler", r) // Ignore the error
	}

	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = store.Flush()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outputHTML(w, r, "static/index.html")
}
