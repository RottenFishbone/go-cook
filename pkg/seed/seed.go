package seed

import (
	"embed"
	"os"
	"path/filepath"

	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
)

// `.cook` files in this directory
//
//go:embed *.cook
var recipes embed.FS

func SeedToDir(path string) {
	// Ensure seed directory exists
	if !common.FileExists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	files, _ := recipes.ReadDir(".")
	for _, entry := range files {
		filename := entry.Name()
		data, _ := recipes.ReadFile(filename)
		newfile, err := os.Create(filepath.Join(path, filename))
		defer newfile.Close()
		if err != nil {
			panic(err)
		}

		newfile.WriteString(string(data))
	}
}
