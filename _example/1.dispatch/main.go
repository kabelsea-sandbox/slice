package main

import (
	"fmt"

	"github.com/kabelsea-sandbox/slice"
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
