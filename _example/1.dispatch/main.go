package main

import (
	"fmt"

	"github.com/kabelsea-sanbox/slice"
)

func SayHello() {
	fmt.Println("Hello!!")
}

func main() {
	slice.Run(
		slice.SetName("invoke-example"),
		slice.SetDispatcher(SayHello),
	)
}
