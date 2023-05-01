package recipe

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"git.sr.ht/~rottenfishbone/go-cook"
)

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

// Decodes a recipe JSON string into a Recipe struct
//
// Returns empty recipe on failure (prints to stderr)
func DecodeFromJson(obj string) cook.Recipe {
	var r cook.Recipe
	err := json.Unmarshal([]byte(obj), &r)
	if err != nil {
		errstr := fmt.Sprintf("Failed to decode JSON recope:\n%v\n", obj)
		os.Stderr.WriteString(errstr)
		return cook.Recipe{
			Name:        "",
			Metadata:    []cook.Metadata{},
			Ingredients: []cook.Ingredient{},
			Cookware:    []cook.Cookware{},
			Timers:      []cook.Timer{},
			Steps:       []cook.Step{},
		}
	}

	return r
}

// Encodes a recipe into a JSON string
//
// Exits on encoding failure
func EncodeToJson(r *cook.Recipe) string {
	bytes, err := json.Marshal(*r)
	if err != nil {
		errstr := fmt.Sprintf("Failed to JSON encode recipe:\n%v\n", *r)
		os.Stderr.WriteString(errstr)
		os.Exit(1)
	}
	return string(bytes)
}

// Prints a recipe to stdout using nice formatting.
func PrettyPrint(recipe *cook.Recipe) {
	fmt.Printf("========= %v ========\n", recipe.Name)
	wr := new(tabwriter.Writer)
	if len(recipe.Metadata) > 0 {
		fmt.Println("Metadata:")
		for _, meta := range recipe.Metadata {
			fmt.Printf("\t%v: %v\n", meta.Tag, meta.Body)
		}
		fmt.Println("")
	}

	if len(recipe.Ingredients) > 0 {
		fmt.Println("Ingredients:")
		wr.Init(os.Stdout, 0, 4, 4, ' ', tabwriter.TabIndent)
		for _, ingr := range recipe.Ingredients {
			var qtyStr string
			if ingr.QtyVal != math.Inf(-1) {
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
			if cookware.QtyVal != math.Inf(-1) {
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
