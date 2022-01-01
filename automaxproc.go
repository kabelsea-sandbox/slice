package slice

import (
	"log"

	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	nop := func(s string, i ...interface{}) {}

	if _, err := maxprocs.Set(maxprocs.Logger(nop)); err != nil {
		log.Printf("automaxprocs set failed, %v", err)
	}
}
