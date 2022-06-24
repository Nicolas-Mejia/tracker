package main

import (
	"fmt"
	"os"
)

func getTrackingNational(id string) {
	l := len(id)
	prefix := id[:2]
	code := id[2 : l-2]
	var uri string = os.Getenv("correoUri")

	fmt.Println(prefix, code, uri)
}
