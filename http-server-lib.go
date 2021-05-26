package main

import (
	"fmt"
	"strings"
)

func concat(a, b string) string {
	return strings.Join([]string{a, b}, " ")
}

func main() {
	var a = "hello"
	var b = "world"
	a = concat(a, b)
	fmt.Println(a)
}
