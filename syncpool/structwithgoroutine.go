package main

import (
	"fmt"
	"strconv"
	"sync"
)

var personCreateCount int = 0

// MyModel 테스트용 구조체
type Person struct {
	Name  string
	Age   int16
	State string
}

// NewMyModel 테스트용 구조체 생성
func NewPerson(name string, age int16) *Person {
	return &Person{
		Name: name,
		Age:  age,
	}
}

func NewPersonPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			personCreateCount++
			fmt.Printf("Create new Person, %v\n", personCreateCount)
			return Person{}
		},
	}
}

func personPoolTest(numTests int) {
	fmt.Println("구조체를 사용한 예제 시작")

	// Person 객체를 관리하는 Pool 생성
	p := NewPersonPool()

	var wg sync.WaitGroup

	for i := 0; i < numTests; i++ {
		wg.Add(1)

		go func(seq int) {
			// Person 객체를 생성하거나 Pool에서 가져오기
			o := p.Get().(Person)

			// Person 객체 사용하기
			o.Name = strconv.FormatInt(int64(seq), 10)
			o.Age = int16(seq)

			// Person 객체 정보 출력
			// 객체 주소를 출력해 보니 가끔 동일한 주소로 가져오는 경우도 있는데, 거의 매번 새로운 주소를 가지고 있다.
			// 객체를 새로 만드는 듯한 느낌이다.
			fmt.Printf("객체 주소: %p, 객체 값: %v\n", &o, o)

			// 사용한 Person 객체를 Pool에 반환하기
			p.Put(o)

			wg.Done()
		}(i)

	}

	wg.Wait()
	fmt.Printf("구조체를 사용한 예제 종료, %n", personCreateCount)
}

var personPtrCreateCount int
var personAddrMap sync.Map

func init() {
	personAddrMap = sync.Map{}
}

func NewPersonPtrPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			personPtrCreateCount++
			fmt.Printf("Create new Person, %v\n", personPtrCreateCount)
			return &Person{}
		},
	}
}

func personPtrPoolTest(numTests int) {
	fmt.Println("구조체의 포인터를 반환하는 Pool을 사용한 예제 시작")

	// Person 객체를 관리하는 Pool 생성
	p := NewPersonPtrPool()

	var wg sync.WaitGroup

	for i := 0; i < numTests; i++ {
		wg.Add(1)

		go func(seq int) {
			// Person 객체를 생성하거나 Pool에서 가져오기
			o := p.Get().(*Person)

			// Person 객체 사용하기
			o.Name = strconv.FormatInt(int64(seq), 10)
			o.Age = int16(seq)

			// Person 객체 정보 출력
			// 객체 주소를 출력해 보니 pool.Get()으로 가져온 객체의 주소는 이전에 사용한 객체의 주소와 동일하다
			// Get() 함수 안을 들여다 볼까?해서 봤더니... 다 이해하기엔 쉽지 않을 듯하다.
			// godoc을 다시 들여다보다가 그 안의 예제를 보니.. 예제안에 주석에 아래와 같은 내용이 있다.
			//
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			//
			// New 함수는 포인터 타입을 반환해야 한단다...
			fmt.Printf("객체 주소: %p, 객체 값: %v\n", o, o)

			personAddrMap.Store(fmt.Sprintf("%p", o), "")
			// 사용한 Person 객체를 Pool에 반환하기
			p.Put(o)

			wg.Done()
		}(i)

	}

	wg.Wait()
	fmt.Printf("구조체를 사용한 예제 종료, %v, %v, \n", personPtrCreateCount, personAddrMap)
}
