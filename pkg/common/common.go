package common

import (
	"os"
)

// Returns whether a file (or directory) exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
