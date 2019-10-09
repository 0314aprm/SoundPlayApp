package main

import (
	"fmt"

	//"./service"
	"kplay.test/service"
)

func main() {
	fmt.Println("starting server")

	s := service.NewServer()
	s.StartUARTService()
	s.StartGATTServer()

	fmt.Println("stopping server")
}
