package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {

}
func main() {
	//   You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("logs/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	defer file.Close()

	// 간단한 예제
	// example()

	// 고루틴을 사용한 예제
	// exampleWithGoroutine()

	// 구조체를 사용한 예제
	// exampleWithStruct()

	// 리스트를 사용한 예제
	// exampleWithList(50)

	// 슬라이스와 구조체를 사용한 예제
	// exampleWithListAndStruct(500)

	// personPoolTest(50)
	personPtrPoolTest(30)
}

func example() {
	fmt.Println("간단한 예제 시작")
	p := sync.Pool{
		// 새로운 슬라이스를 만든다
		New: func() interface{} {
			slice := make([]int, 0)
			fmt.Printf("%p, new slice\n", slice[:0])
			return slice
		},
	}

	for n := 0; n < 50; n++ {
		slice := p.Get().([]int)

		prevLen := len(slice)
		slice = append(slice, 1, 2, 3, 4, 5)

		fmt.Printf("%p, %v, %v\n", slice[:0], prevLen, slice)

		// 사용을 마친 슬라이스는 초기화하여 풀에 반환합니다.
		p.Put(slice[:0])
	}
}

func exampleWithGoroutine() {
	fmt.Println("고루틴을 사용한 예제 시작")
	// 길이가 0이고 capacity 5인 슬라이스 풀
	p := sync.Pool{
		New: func() interface{} {
			return make([]int, 0, 5)
		},
	}

	var wg sync.WaitGroup
	for n := 0; n < 10; n++ {
		wg.Add(1)
		go func() {
			// 꺼내오기
			s := p.Get().([]int)

			// 슬라이스로 작업하기
			s = append(s, 1, 2, 3, 4, 5)
			fmt.Println(s)

			// 반환
			p.Put(s[:0])
			wg.Done()
		}()
	}

	// 모든 루틴이 종료될 때까지 기다리기
	wg.Wait()
}

// MyModel 테스트용 구조체
type MyModel struct {
	Name  string
	Age   int16
	State string
}

// NewMyModel 테스트용 구조체 생성
func NewMyModel(name string, age int16) *MyModel {
	return &MyModel{
		Name: name,
		Age:  age,
	}
}

func NewMyModelPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			fmt.Println("Create new MyModel")
			return MyModel{}
		},
	}
}

func exampleWithStruct() {
	fmt.Println("구조체를 사용한 예제 시작")
	p := NewMyModelPool()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(seq int) {
			m := p.Get().(MyModel)

			m.Name = strconv.FormatInt(int64(seq), 10)
			m.Age = int16(seq)
			fmt.Printf("%p %v\n", &m, m)

			p.Put(m)

			wg.Done()
		}(i)

	}

	wg.Wait()
	fmt.Println("구조체를 사용한 예제 종료")
}

type MyModelList []MyModel

func NewMyModelListPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return make(MyModelList, 0)
		},
	}
}

// 커스텀 구조체 슬라이스
func exampleWithList(numTests int) {
	fmt.Println("구조체 슬라이스를 사용한 예제 시작")
	p := NewMyModelListPool()

	for i := 0; i < numTests; i++ {
		if i == numTests-1 {
			time.Sleep(time.Second * 6)
		}
		l := p.Get().(MyModelList)
		rt := reflect.TypeOf(l)
		switch rt.Kind() {
		case reflect.Slice:
			fmt.Println("is a slice with element type", rt.Elem())
		case reflect.Array:
			fmt.Println("is an array with element type", rt.Elem())
		default:
			fmt.Println("is something else entirely")
		}
		prevLen := len(l)
		if prevLen == 0 {
			l = append(l, *NewMyModel(strconv.FormatInt(int64(i), 10), int16(i)))
		}
		// l = append(l, *NewMyModel(strconv.FormatInt(int64(i), 10), int16(i)))
		// l = append(l, *NewMyModel(strconv.FormatInt(int64(i), 10), int16(i)))
		// l = append(l, *NewMyModel(strconv.FormatInt(int64(i), 10), int16(i)))

		// fmt.Printf("%p %v\n", l[:0], l)
		log.WithFields(log.Fields{
			"addr": fmt.Sprintf("%p", l[:0]),
		}).Info(l, prevLen)

		p.Put(l[:0])
	}
}

// 커스텀 구조체 슬라이스
func exampleWithListWithGoroutine(numTests int) {
	fmt.Println("구조체 슬라이스를 사용한 예제 시작")
	p := NewMyModelListPool()

	var wg sync.WaitGroup
	for i := 0; i < numTests; i++ {
		wg.Add(1)
		go func(seq int) {
			if seq == numTests-1 {
				time.Sleep(time.Second * 6)
			}
			l := p.Get().(MyModelList)
			if len(l) == 0 {
				l = append(l, *NewMyModel(strconv.FormatInt(int64(seq), 10), int16(seq)))
			}
			// l = append(l, *NewMyModel(strconv.FormatInt(int64(seq), 10), int16(seq)))
			// l = append(l, *NewMyModel(strconv.FormatInt(int64(seq), 10), int16(seq)))
			// l = append(l, *NewMyModel(strconv.FormatInt(int64(seq), 10), int16(seq)))

			// fmt.Printf("%p %v\n", l[:0], l)
			log.WithFields(log.Fields{
				"addr": fmt.Sprintf("%p", l[:0]),
			}).Info(l, len(l))

			p.Put(l[:0])
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// 커스텀 구조체 슬라이스
func exampleWithListAndStruct(numTests int) {
	fmt.Println("구조체 슬라이스를 사용한 예제 시작")
	pl := NewMyModelListPool()
	ps := NewMyModelPool()

	var wg sync.WaitGroup
	for i := 0; i < numTests; i++ {
		wg.Add(1)
		go func(seq int) {
			if seq == numTests-1 {
				time.Sleep(time.Second * 6)
			}
			l := pl.Get().(MyModelList)

			for j := 0; j < 1000; j++ {
				m := ps.Get().(MyModel)
				// fmt.Printf("%p new state is %v\n", &m, m.State)
				log.WithFields(log.Fields{
					"addr": fmt.Sprintf("%p", &m),
				}).Infof("new state is %v\n", m.State)
				m.Name = strconv.FormatInt(int64(seq), 10)
				m.Age = int16(j)
				m.State = "Assigned"
				l = append(l, m)
			}
			log.WithFields(log.Fields{
				"addr": fmt.Sprintf("%p", l[:0]),
			}).Infof("%v", l[:4])
			// fmt.Printf("%p %p %v\n", l[:0], &l[0], l)

			for j := len(l) - 1; j >= 0; j-- {
				ps.Put(l[j])
			}
			pl.Put(l[:0])
			wg.Done()
		}(i)
	}
	wg.Wait()
}
