package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// 파일 열기
	f, err := os.OpenFile("data/test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := f.Write([]byte(time.Now().String()))
	if err != nil {
		panic(err)
	}
	fmt.Printf("written %d bytes", n)
}
