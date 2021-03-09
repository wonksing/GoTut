package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("hello world")
	// Len is actual length left in the reader
	// Size is the original length of the input string
	fmt.Printf("length: %d, size: %d\n", r.Len(), r.Size())

	// buf 에 읽어보자, capacity를 4로 지정한다.
	buf := make([]byte, 4)
	// "hello world" 문자열의 0번째 인덱스부터 4 바이트를 buf에 읽어온다.
	n, err := r.Read(buf)
	if err != nil {
		return
	}
	fmt.Printf("read %d bytes, %s\n", n, buf)

	// 읽어온 후에 Reader "r"의 첫 4 바이트는 소비되어 사라진다. buf로 이동된 것이다.
	fmt.Printf("length: %d, size: %d\n", r.Len(), r.Size())

	// Reader "r"을 "wonk"라는 문자열로 리셋한다.
	// length와 size 역시 리셋된다.
	r.Reset("wonk")
	fmt.Printf("length: %d, size: %d\n", r.Len(), r.Size())

	// 리더의 시작(SeekStart)에서부터 2번째 인덱스까지 오프셋을 이동시킨다.
	// 0 | 1 | 2 | 3
	// 인덱스가 0에서 시작해서 2까지 이동한다.
	sk, err := r.Seek(2, io.SeekStart)
	if err != nil {
		return
	}
	fmt.Printf("%d\n", sk)

	// 현재 오프셋(SeekCurrent)에서부터 2번째 인덱스까지 오프셋을 이동시킨다.
	// 위에서 2까지 이동했고, 2번 더 이동시킨다.
	// 이러면 오프셋이 인덱스 4에 위치하게 되지만, out of index 오류는 없다. 왜일까?
	sk, err = r.Seek(2, io.SeekCurrent)
	if err != nil {
		return
	}
	fmt.Printf("%d\n", sk)

	// 오프셋이 4에 위치해 있고 읽을 것이 없으므로 Len는 0이다.
	// 하지만 Size는 여전히 4이다.
	fmt.Printf("length: %d, size: %d\n", r.Len(), r.Size())

	// 읽어보자, 읽은 바이트 수는 0이다.
	buf2 := make([]byte, 4)
	n, err = r.Read(buf2)
	if err != nil {
		// EOF 오류를 반환한다.
		if err == io.EOF {
			fmt.Printf("EOF이니 봐줘... %v\n", err)
		} else {
			fmt.Printf("error %v\n", err)
			return
		}
	}
	fmt.Printf("read %d bytes, %s, %v\n", n, buf2, buf2)

	// buf2 = nil // len:0, cap:0
	// buf2 = buf2[:0] // len:0, cap:4
	// fmt.Printf("buf2: %v, %v\n", len(buf2), cap(buf2))

	// 그럼 다시
	// 리더안에 있는 문자열보다 더 많이 읽어보자
	r.Reset("abc")
	fmt.Printf("length: %d, size: %d\n", r.Len(), r.Size())

	// 슬라이스를 nil로 만들거나 비운 후에는 읽히지 않는다. 왜일까?
	// Read implements the io.Reader interface.
	//
	// buf2 = nil
	// buf2 = buf2[:0]
	//
	// 하지만, 새로 할당하거나 그대로 사용하면 읽힌다.
	// buf2 = make([]byte, 4)
	//
	// Read는 copy를 사용하고 이 함수는 대상 슬라이스의 길이가 0이면 카피하지 않는다.
	// func (r *Reader) Read(b []byte) (n int, err error) {
	// 	if r.i >= int64(len(r.s)) {
	// 		return 0, io.EOF
	// 	}
	// 	r.prevRune = -1
	// 	n = copy(b, r.s[r.i:])
	// 	r.i += int64(n)
	// 	return
	// }
	//

	n, err = r.Read(buf2)
	if err != nil {
		// EOF 오류를 반환한다.
		if err == io.EOF {
			fmt.Printf("EOF이니 봐줘... %v\n", err)
		} else {
			fmt.Printf("error %v\n", err)
			return
		}
	}
	fmt.Printf("read %d bytes, %s\n", n, buf2)
}
