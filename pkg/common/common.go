package common

import (
	"fmt"
	"os"
)

// Returns whether a file (or directory) exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ShowError(err error) {
	errstr := fmt.Sprintf("Error: %v\n", err.Error())
	os.Stderr.WriteString(errstr)
}
