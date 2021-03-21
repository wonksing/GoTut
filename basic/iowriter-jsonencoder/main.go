package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	withBytesBuffer()
}

func withBytesBuffer() {
	buf := &bytes.Buffer{}
	u := Person{
		Name: "wonkdsing",
		Age:  1000,
	}
	err := json.NewEncoder(buf).Encode(u)
	if err != nil {
		panic(err)
	}
	fmt.Print(buf.String())
}
