package main

import (
	"fmt"
	"strconv"
	"sync"
)

type Person struct {
	Name string
	ID   string
}
type PersonList []Person

type People struct {
	l sync.RWMutex
	m *sync.Map
}

func NewPeople() *People {
	m := &sync.Map{}
	return &People{
		m: m,
		l: sync.RWMutex{},
	}
}

func (o *People) Init(pl PersonList) {

	m := &sync.Map{}
	for _, p := range pl {
		m.Store(p.ID, &p)
	}

	o.l.Lock()
	defer o.l.Unlock()
	o.m = m
}

func (o *People) Add(p Person) {
	o.l.Lock()
	defer o.l.Unlock()

	o.m.Store(p.ID, p)
}

func (o *People) Find(ID string) *Person {
	o.l.RLock()
	defer o.l.RUnlock()

	v, ok := o.m.Load(ID)
	if ok {
		return v.(*Person)
	}
	return nil
}

var (
	people *People
	wg     sync.WaitGroup
)

func init() {
	people = NewPeople()

	wg = sync.WaitGroup{}
}

func main() {
	// map 채우기
	pl := make(PersonList, 0)
	for i := 0; i < 3; i++ {
		pl = append(pl, Person{
			Name: "W" + strconv.Itoa(i),
			ID:   strconv.Itoa(i),
		})
	}
	people.Init(pl)

	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			pl := make(PersonList, 0)
			for i := 0; i < 3; i++ {
				pl = append(pl, Person{
					Name: "W" + strconv.Itoa(i),
					ID:   strconv.Itoa(i),
				})
			}
			people.Init(pl)
		}
		wg.Done()
	}()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			p := people.Find("1")
			fmt.Println(p.ID, p.Name)
			wg.Done()
		}()
	}

	wg.Wait()
}
