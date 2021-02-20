package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	ip   string
	port int
	key  string
	cert string
	addr string

	filePath string
)

func init() {

	flag.StringVar(&ip, "i", "", "ip")
	flag.IntVar(&port, "p", 8080, "port")
	flag.StringVar(&key, "K", "./certs/server.key", "key file path")
	flag.StringVar(&cert, "P", "./certs/server.crt", "cert file path")
	flag.StringVar(&filePath, "filepath", "./download/", "path to download files from")

	flag.Parse()

	addr = fmt.Sprintf("%v:%v", ip, strconv.Itoa(port))
}
func main() {
	tlsConfig := &tls.Config{
		// MinVersion: tls.VersionTLS12,
		// MinVersion: tls.VersionTLS11,
		MinVersion:               tls.VersionTLS10, // weak, only for xp
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}

	// 라우터, gorilla mux를 쓴다
	router := mux.NewRouter()

	// 파일 목록 조회용
	router.Handle("/testapi/file/download", handleFindFileList(filePath)).Methods("GET")

	// 정적 파일 서버 (자세한건 건너뛰자)
	fileServer := http.FileServer(http.Dir(filePath))

	// 정적 파일 서버를 api와 매핑한다
	// http.StripPrefix를 이용해서 서버 디스크의 경로를 제대로 찾아갈 수 있게 한다.
	router.PathPrefix("/testapi/file/download/").
		Handler(http.StripPrefix("/testapi/file/download/", handleFileServe(fileServer)))

	// http 서버 생성
	httpServer := &http.Server{
		Addr:         addr,                                                            // listen 할 주소(ip:port)
		TLSConfig:    tlsConfig,                                                       // tls 설정
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0), // 음..
		WriteTimeout: 30 * time.Second,                                                // 서버 > 클라이언트 응답
		ReadTimeout:  30 * time.Second,                                                // 클라이언트 > 서버 요청
		Handler:      router,                                                          // mux다
	}

	log.Fatal(httpServer.ListenAndServeTLS(cert, key))
}

// FileEntity 파일정보 구조체
// 파일목록 요청시 이 구조체를 json 문자열로 응답한다
type FileEntity struct {
	FileNam string `json:"file_nam"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
}

// FileEntityList FileEntity 구조체의 배열
type FileEntityList []FileEntity

// basePath 경로 안에 있는 모든 폴더, 파일 리스트를 찾아 json 문자열로 변환하여 클라이언트에 응답한다
func handleFindFileList(basePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 다운로드할 파일이 있는 절대경로 가져오기
		// 클라이언트에게 해당 경로에 있는 폴더, 파일명만 전달하기 위해서
		absPath, err := filepath.Abs(basePath)
		if err != nil {
			log.Printf("%v\n", err)
			return
		}
		absPath = absPath + string(os.PathSeparator)

		// 파일정보를 추가할 slice를 만든다
		list := make(FileEntityList, 0)

		// 경로 안의 파일과 폴더를 찾아서 list에 추가한다
		err = filepath.Walk(absPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.Compare(absPath, path) != 0 {
					// log.Println(strings.Replace(path, absPath, "", 1), info.Size(), info.IsDir())
					list = append(list, FileEntity{
						FileNam: strings.Replace(path, absPath, "", 1),
						Size:    info.Size(),
						IsDir:   info.IsDir(),
					})
				}
				return nil
			})
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		// list를 json 문자열로 변환하여 응답본문을 만들어 전달한다
		err = json.NewEncoder(w).Encode(&list)
		if err != nil {
			log.Printf("%v\n", err)
		}
	})
}

// 다운로드 요청한 파일을 전달한다.
// 파일 다운로드는 한개씩만 가능하다.
func handleFileServe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
