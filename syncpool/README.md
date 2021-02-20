# Gotut - sync.Pool

## 배경
Go 역시 힙영역에 할당된 메모리를 회수하는 Garbage Collecting 언어이다. 새로 생성한 객체는 힙영역에 할당되고 이 객체가 더이상 참조되지 않는다고 여겨질 때 메모리를 회수하기 위해 GC를 수행한다. 이 때 오버헤드가 생기며 어플리케이션이 잠시 멈추거나? 느려지거나? 하는데 이를 최소화 하기 위해 sync.Pool을 고려해 보았다. 고려해 보기만 하고 써먹진 않고 있다. 아직 Go App에서 이런 증상을 만나보질 못했다... 좋은 언어이다.

## sync.Pool 테스트
객체 주소를 출력해 보니 가끔 동일한 주소로 가져오는 경우도 있는데, 거의 매번 새로운 주소를 가지고 있다.
객체를 새로 만드는 듯한 느낌이다. 그래서 아래와 같이 sync.Pool의 New 함수의 리턴 타입을 포인터로 바꾸니..
예상했던 결과를 봤다.

```go
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
```
콘솔에 출력된 내용
```text
구조체의 포인터를 반환하는 Pool을 사용한 예제 시작
Create new Person, 1
Create new Person, 3
Create new Person, 4
객체 주소: 0x14000188000, 객체 값: &{0 0 }
Create new Person, 7
Create new Person, 5
객체 주소: 0x1400006c2a0, 객체 값: &{15 15 }
객체 주소: 0x140001880c0, 객체 값: &{13 13 }
Create new Person, 6
객체 주소: 0x1400006c2a0, 객체 값: &{16 16 }
객체 주소: 0x14000188150, 객체 값: &{14 14 }
객체 주소: 0x14000188000, 객체 값: &{29 29 }
객체 주소: 0x1400006c2a0, 객체 값: &{17 17 }
Create new Person, 2
객체 주소: 0x140001880c0, 객체 값: &{18 18 }
객체 주소: 0x140001881e0, 객체 값: &{11 11 }
객체 주소: 0x14000188000, 객체 값: &{2 2 }
객체 주소: 0x1400009e000, 객체 값: &{10 10 }
객체 주소: 0x14000188000, 객체 값: &{1 1 }
객체 주소: 0x14000188150, 객체 값: &{20 20 }
객체 주소: 0x14000188000, 객체 값: &{24 24 }
Create new Person, 8
객체 주소: 0x1400006c2a0, 객체 값: &{25 25 }
객체 주소: 0x140001882d0, 객체 값: &{21 21 }
객체 주소: 0x1400006c2a0, 객체 값: &{26 26 }
객체 주소: 0x140001882d0, 객체 값: &{22 22 }
객체 주소: 0x1400006c2a0, 객체 값: &{27 27 }
객체 주소: 0x140001882d0, 객체 값: &{23 23 }
객체 주소: 0x1400006c2a0, 객체 값: &{28 28 }
객체 주소: 0x140001882d0, 객체 값: &{6 6 }
객체 주소: 0x1400006c240, 객체 값: &{12 12 }
객체 주소: 0x1400006c2a0, 객체 값: &{7 7 }
객체 주소: 0x140001880c0, 객체 값: &{19 19 }
객체 주소: 0x140001881e0, 객체 값: &{5 5 }
객체 주소: 0x140001882d0, 객체 값: &{4 4 }
객체 주소: 0x1400006c2a0, 객체 값: &{3 3 }
객체 주소: 0x1400009e000, 객체 값: &{9 9 }
객체 주소: 0x14000188150, 객체 값: &{8 8 }
구조체를 사용한 예제 종료, 8, {{0 0} {{map[] true}} map[0x1400006c240:0x14000190020 0x1400006c2a0:0x1400000e040 0x1400009e000:0x140000a8000 0x14000188000:0x14000190000 0x140001880c0:0x14000190008 0x14000188150:0x14000190010 0x140001881e0:0x1400000e048 0x140001882d0:0x14000190018] 0}, 
```
## 결론
Pool에서 객체를 go routine 안에서 총 50회 시도했다. 새로운 객체를 생성하는 것은 총 8회이고, 나머지 42회는 캐시에서 가져온 것으로 보이며 객체의 주소 역시 생성했던 8회의 주소와 동일했다. Godoc의 예제 안에 있던 주석 내용대로 포인터를 반환하는 것이 맞는 것으로 보인다.