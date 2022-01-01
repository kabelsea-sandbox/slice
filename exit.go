package slice

import (
	"fmt"
	"os"
)

// exitError var for testing purposes
var exitError = defaultExitError

// defaultExitError prints error and exist with status
func defaultExitError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
