package recipe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"git.sr.ht/~rottenfishbone/go-cook"
	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
)

// Reads a recipe file and parses it into a Recipe struct.
//
// nil returned on failure
func LoadFromFile(path string) *cook.Recipe {
	if !common.FileExists(path) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		errstr := fmt.Sprintf("Error: %v\n", err.Error())
		os.Stderr.WriteString(errstr)
	}

	r := cook.ParseRecipe(FilepathToName(path), &data)
	return &r
}

// Name a recipe using its filepath
func FilepathToName(path string) string {
	path = filepath.Base(path)
	ext := filepath.Ext(path)

	path = strings.TrimRight(path, ext)
	path = strings.Map(func(r rune) rune {
		if r == '_' || r == '-' {
			return ' '
		}
		return r
	}, path)

	return path
}

// Prints a recipe to stdout using nice formatting.
func PrettyPrint(recipe *cook.Recipe) {
	fmt.Printf("========= %v ========\n", recipe.Name)
	wr := new(tabwriter.Writer)
	if len(recipe.Metadata) > 0 {
		fmt.Println("Metadata:")
		for k, v  := range recipe.Metadata {
			fmt.Printf("\t%v: %v\n", k, v)
		}
		fmt.Println("")
	}

	if len(recipe.Ingredients) > 0 {
		fmt.Println("Ingredients:")
		wr.Init(os.Stdout, 0, 4, 4, ' ', tabwriter.TabIndent)
		for _, ingr := range recipe.Ingredients {
			var qtyStr string
			if ingr.QtyVal != cook.NoQty {
				qtyStr = fmt.Sprintf("%v %v", ingr.QtyVal, ingr.Unit)
			} else {
				qtyStr = ingr.Qty + " " + ingr.Unit
			}

			fmt.Fprintf(wr, "\t%v\t%v\n", ingr.Name, qtyStr)
		}
		wr.Flush()
		fmt.Println("")
	}

	if len(recipe.Cookware) > 0 {
		fmt.Println("Cookware:")
		wr.Init(os.Stdout, 0, 4, 4, ' ', tabwriter.TabIndent)
		for _, cookware := range recipe.Cookware {
			var qtyStr string
			if cookware.QtyVal != cook.NoQty {
				qtyStr = fmt.Sprintf("%v %v", cookware.QtyVal, cookware.Unit)
			} else {
				qtyStr = cookware.Qty + " " + cookware.Unit
			}

			fmt.Fprintf(wr, "\t%v\t%v\n", cookware.Name, qtyStr)
		}
		wr.Flush()
		fmt.Println("")
	}

	if len(recipe.Steps) > 0 {
		fmt.Println("Steps:")
		var builder strings.Builder
		for n, step := range recipe.Steps {
			fmt.Printf("\t%v. ", n+1)
			for _, chunk := range step {
				builder.WriteString(chunk.ToString())
			}
			fmt.Println(builder.String())
			builder.Reset()
		}
		fmt.Println("")
	}
}
