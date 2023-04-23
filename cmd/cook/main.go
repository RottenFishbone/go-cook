package main

import (
	"fmt"

	"git.sr.ht/~rottenfishbone/cooklang-go/pkg/config"
)

func main() {
	cfg, success := config.LoadConfig("")
	if !success {
		cfg = config.ConfigInit("", "", "")
	}

	fmt.Printf("cfg: %v\n", cfg)
}
