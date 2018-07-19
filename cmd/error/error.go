package error

import (
	"fmt"
	"os"
)

// ExitWithError prints out an error code and an error string to stderr.
func ExitWithError(code int, err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(code)
}
