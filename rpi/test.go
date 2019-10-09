package main

import (
	"fmt"
)

func main() {
	array := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	for _, v := range array {
		fmt.Printf("%X ", v)
	}
}
