package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// var utf8Val encoding.UTF8Validator

func IsUTF8(str string) bool {
	dst := make([]byte, len(str))

	nDst, nSrc, err := encoding.UTF8Validator.Transform(dst, []byte(str), true)
	if err != nil {
		return false
	}
	log.Println(str, nDst, nSrc)
	return true
}

func IsUtf8(str string) bool {
	return utf8.ValidString(str)
}

func EuckrToUtf8(str string) string {
	if IsUTF8(str) {
		return str
	}
	r := transform.NewReader(strings.NewReader(str), korean.EUCKR.NewDecoder())
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return string(str)
	}
	return string(b)
}

func main() {

	tmp := "BEC8B3E7" + "BEC8B3E7"
	bs, err := hex.DecodeString(tmp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs), len(bs), len(string(bs)))

	str := "안녕" + string(bs)
	log.Println(IsUTF8(str))
	log.Println(IsUtf8(str))
	log.Println(IsUtf8("고마워"))
}
