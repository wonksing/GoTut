package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	requestCredential()
	requestToken()
}

func requestCredential() {
	// GET 호출
	resp, err := http.Get("http://localhost:9096/credentials")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(data))
}

func requestToken() {
	// GET 호출
	resp, err := http.Get("http://localhost:9096/token?grant_type=client_credentials&client_id=12341234&client_secret=12341234&scope=all")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(data))
}
