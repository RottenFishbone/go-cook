package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Returns whether a file (or directory) exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Converts a relative path to an absolute path, disallowing paths that travel
// outside of the root directory.
func SanitizeRelPath(root string, path string) (string, error) {
	outPath, err := filepath.Abs(filepath.Join(root, path))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(outPath, root) {
		return "", errors.New("Relative path escapes root directory.")
	}

	return outPath, nil
}

func ShowError(err error) {
	errstr := fmt.Sprintf("Error: %v\n", err.Error())
	os.Stderr.WriteString(errstr)
}
