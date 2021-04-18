package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

func main() {
	// 선언, 초기화
	m := sync.Map{}

	// 저장
	m.Store("wonk", "12")

	// 찾기
	v, ok := m.Load("wonk")
	if ok {
		fmt.Println("found", v)
	} else {
		fmt.Println("not found")
	}

	// 찾기+저장
	actual, loaded := m.LoadOrStore("mink", "23")
	if loaded {
		// 찾은 값 반환하지만, 저장은 하지 않는다
		fmt.Println("found", actual)
	} else {
		// 저장하고 저장한 값 반환
		fmt.Println("stored", actual)
	}

	// 삭제
	m.Delete("wonk12")
	m.Delete("wonk")

	// 맵 순환
	m.Range(func(k, v interface{}) bool {
		fmt.Printf("key: %v, value: %v\n", k, v)
		return true // false이면 종료
	})

	// 맵 순환중에 삭제?
	// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
	// Range does not necessarily correspond to any consistent snapshot of the Map's contents: no key will be visited more than once, but if the value for any key is stored or deleted concurrently, Range may reflect any mapping for that key from any point during the Range call.
	// Range may be O(N) with the number of elements in the map even if f returns false after a constant number of calls.
	m.Range(func(k, v interface{}) bool {
		if k.(string) == "mink" {
			m.Delete("mink")
		}
		fmt.Printf("key: %v, value: %v\n", k, v)
		return true // false이면 종료
	})

	test()
}

func test() {
	m := sync.Map{}
	var wg sync.WaitGroup

	wg.Add(1)
	go func(prefix string) {
		for i := 0; i < 1000; i++ {
			m.Store(i, prefix+strconv.Itoa(rand.Intn(10000)))
		}
		wg.Done()
	}("A")

	wg.Add(1)
	go func(prefix string) {
		for i := 0; i < 1000; i++ {
			m.Store(i, prefix+strconv.Itoa(rand.Intn(10000)))
		}
		wg.Done()
	}("B")

	wg.Wait()

	m.Range(func(k, v interface{}) bool {
		fmt.Printf("key: %v, value: %v\n", k, v)
		return true
	})
}
