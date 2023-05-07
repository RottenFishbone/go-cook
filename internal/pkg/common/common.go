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

// Recursively delete all empty directories from passed root
func CleanupEmptyDir(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	// Clean all subdirectories first, if possible
	if len(entries) > 0 {
		for _, entry := range entries {
			if entry.IsDir() {
				path := filepath.Join(root, entry.Name())
				if err = CleanupEmptyDir(path); err != nil {
					return err
				}
			}
		}
		// Check the dir again after cleanup
		if entries, err = os.ReadDir(root); err != nil {
			return err
		}
	}

	// Delete root if its empty
	if len(entries) == 0 {
		if err = os.Remove(root); err != nil {
			return err
		}
	}

	return nil
}

// Prints an error to Stderr
func ShowError(err error) {
	if err == nil {
		// Prevent null pointer references
		return
	}
	errstr := fmt.Sprintf("ERR: \t%v\n", err.Error())
	os.Stderr.WriteString(errstr)
}
