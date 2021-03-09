package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type HandlerParams struct {
	// Key           string
	HmacHeaderKey string
	SecretKey     string
}

func main() {
	hp := HandlerParams{
		HmacHeaderKey: "hmackey",
		SecretKey:     "asdf",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, HandlerParams{}, hp)

	buf := bytes.NewBuffer([]byte(strings.Repeat("a", 30)))
	readBuf := buf.Bytes()
	fmt.Printf("read buffer %s\n", readBuf)

	readBuf2 := buf.Bytes()
	fmt.Printf("read buffer %s\n", readBuf2)

}
