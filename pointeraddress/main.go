package main

import "fmt"

func main() {
	printAddrOfPtr()
}

// 포인터 주소 출력하기
func printAddrOfPtr() {
	fmt.Println("포인터 주소 출력하기 시작")

	val := 2
	addr := fmt.Sprintf("%p", &val)
	fmt.Printf("주소: %v, 값: %v\n", addr, val)

	fmt.Println("포인터 주소 출력하기 끝")
}
