# Gotut - 포인터 주소 가져오기

fmt 패키지의 Sprintf 함수를 이용하여 포인터 타입의 주소를 가져올 수 있다.

포인터 주소 가져오기 샘플 코드
```go
// 포인터 주소 출력하기
func printAddrOfPtr() {
	fmt.Println("포인터 주소 출력하기 시작")

	val := 2
	addr := fmt.Sprintf("%p", &val)
	fmt.Printf("주소: %v, 값: %v\n", addr, val)

	fmt.Println("포인터 주소 출력하기 끝")
}

```

출력 결과
```text
포인터 주소 출력하기 시작
주소: 0x140000140b0, 값: 2
포인터 주소 출력하기 끝
```