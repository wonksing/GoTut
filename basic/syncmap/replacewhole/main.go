package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	m  *sync.Map
	wg sync.WaitGroup
)

func init() {
	m = &sync.Map{}
	wg = sync.WaitGroup{}
}
func main() {
	// map 채우기
	m = createNewMap("A")

	// 키 찾고 있기
	wg.Add(1)
	go func() {
		attempts := 0
		for {
			attempts++
			// if v, ok := m.Load(0); ok {
			// 	fmt.Printf("key: %v, value: %v\n", 0, v)
			// } else {
			// 	fmt.Printf("key %v not found\n", 0)
			// }
			_, ok := m.Load(0)
			if ok == false {
				fmt.Printf("key %v not found\n", 0)
			}
			if attempts > 50000 {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		wg.Done()
	}()

	// 새로운 맵으로 대치
	wg.Add(1)
	go func() {
		attempts := 0
		for {
			attempts++
			m = createNewMap("B" + strconv.Itoa(attempts) + "_")

			if attempts > 50000 {
				break
			}
			time.Sleep(10 * time.Millisecond)
			// time.Sleep(300 * time.Millisecond)
		}
		wg.Done()
	}()
	wg.Wait()
}

func createNewMap(prefix string) *sync.Map {
	newMap := sync.Map{}
	for i := 0; i < 1; i++ {
		newMap.Store(i, prefix+strconv.Itoa(rand.Intn(10000)))
	}
	return &newMap
}
