package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	accessKey := os.Getenv("MY_TELEGRAM_ACCESSKEY")
	chatID := os.Getenv("MY_TELEGRAM_CHATID")
	requestTele("https://api.telegram.org", accessKey, chatID, "sendMessage", "json test")
	requestTeleFormData("https://api.telegram.org", accessKey, chatID, "sendMessage", "form test")
}

func requestTele(url, accessKey, chatID, path, text string) {
	fmt.Println("testing http client")

	maxIdleConn := 10
	idleConnTimeoutSec := 20
	disableCompression := true
	config := &tls.Config{InsecureSkipVerify: true}
	tran := CreateTransport(maxIdleConn, idleConnTimeoutSec, disableCompression, config)

	clientTimeoutSec := 20
	httpClient := NewHttpClient(url, tran, clientTimeoutSec)

	// header
	reqHeaders := make(map[string]string)
	reqHeaders["Content-Type"] = "application/json"

	// body
	reqBodyMap := make(map[string]interface{})
	reqBodyMap["chat_id"] = chatID
	reqBodyMap["text"] = text
	req, err := json.Marshal(reqBodyMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(req))

	reqQueries := make(map[string]string)
	reqQueries["chat_id"] = chatID
	reqQueries["text"] = text

	body, headers, err := httpClient.Request("/"+accessKey+"/"+path, "POST", reqHeaders, nil, bytes.NewBuffer(req))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v %v", headers, string(body))
}

func requestTeleFormData(URL, accessKey, chatID, path, text string) {
	fmt.Println("testing http client")

	maxIdleConn := 10
	idleConnTimeoutSec := 20
	disableCompression := true
	config := &tls.Config{InsecureSkipVerify: true}
	tran := CreateTransport(maxIdleConn, idleConnTimeoutSec, disableCompression, config)

	clientTimeoutSec := 20
	httpClient := NewHttpClient(URL, tran, clientTimeoutSec)

	// header
	reqHeaders := make(map[string]string)
	reqHeaders["Content-Type"] = "application/x-www-form-urlencoded"

	// body
	reqBody := url.Values{
		"chat_id": {chatID},
		"text":    {text},
	}

	body, headers, err := httpClient.Request("/"+accessKey+"/"+path, "POST", reqHeaders, nil, strings.NewReader(reqBody.Encode()))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v %v", headers, string(body))
}

type HttpClient struct {
	Url    string
	Header map[string]string
	client *http.Client
}

func NewHttpClient(url string, tran *http.Transport, clientTimeoutSec int) *HttpClient {
	var hc HttpClient
	hc.Url = url

	// hc.client = &http.Client{
	// 	Timeout:   time.Duration(clientTimeoutSec) * time.Second,
	// 	Transport: tran,
	// }
	hc.client = &http.Client{}
	hc.client.Timeout = time.Duration(clientTimeoutSec) * time.Second
	if tran != nil {
		hc.client.Transport = tran
	}
	return &hc
}

func CreateTransport(maxIdleConn int, idleConnTimeoutSec int, disableCompression bool, tlsClientConfig *tls.Config) *http.Transport {
	tr := &http.Transport{
		MaxIdleConns:       maxIdleConn,
		IdleConnTimeout:    time.Duration(idleConnTimeoutSec) * time.Second,
		DisableCompression: disableCompression,
		TLSClientConfig:    tlsClientConfig,
	}
	return tr
}
func (hc *HttpClient) Request(path string, method string, reqHeader map[string]string, reqQueries map[string]string, reqBody io.Reader) ([]byte, map[string][]string, error) {
	var err error

	// var reqBodyBuff *bytes.Buffer
	// if reqBody != nil {
	// 	reqBodyBuff = bytes.NewBuffer(reqBody)
	// } else {
	// 	reqBodyBuff = nil
	// }
	// request, err := http.NewRequest(method, hc.Url+path, bytes.NewBuffer(reqBody))
	request, err := http.NewRequest(method, hc.Url+path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	if len(reqQueries) > 0 {
		q := request.URL.Query()
		for k, v := range reqQueries {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}

	for k, v := range reqHeader {
		request.Header.Set(k, v)
	}

	resp, err := hc.client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// var headers map[string][]string
	headers := make(map[string][]string)
	for k, v := range resp.Header {
		headers[k] = v
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return res, headers, err
}
