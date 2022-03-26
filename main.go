package main

import (
	"fmt"
)

func main() {
	var x int = 2
	s := 8
	p := &s
	fmt.Println(&x, p, *p)
}
