package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"git.sr.ht/~rottenfishbone/go-cook"
	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
)

//go:embed template.txt
var template string

const TESTS_SRC = "https://raw.githubusercontent.com/cooklang/spec/main/tests/canonical.yaml"

type Spec struct {
	Version int             `yaml:"version"`
	Tests   map[string]Test `yaml:"tests"`
}

type Test struct {
	Source string `yaml:"source"`
	Result Result `yaml:"result"`
}

type Result struct {
	Steps    [][]Chunk         `yaml:"steps"`
	Metadata map[string]string `yaml:"metadata"`
}

type Chunk struct {
	Type     string `yaml:"type"`
	Value    string `yaml:"value,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Quantity string `yaml:"quantity,omitempty"`
	Units    string `yaml:"units,omitempty"`
}

func main() {
	var err error
	args := os.Args[1:]

	// Grab the path from args
	var path string
	if len(args) == 1 {
		path = args[0]
	} else if len(args) == 0 {
		path = "./"
	} else {
		os.Stderr.WriteString("Error: too many arguments.\n")
		os.Exit(1)
	}

	// Make path absolute
	if path, err = filepath.Abs(path); err != nil {
		common.ShowError(err)
		os.Exit(1)
	}

	// Pull the test spec yaml from the `spec` repo
	test_spec_path := filepath.Join(path, "canonical.yaml")
	if err = downloadSource(test_spec_path); err != nil {
		common.ShowError(err)
		os.Exit(1)
	}
	fmt.Printf("Downloaded test spec into: %s\n", test_spec_path)

	// Build the tests `.go` file
	test_path := filepath.Join(path, "canonical_test.go")
	if err = generateTestFile(test_spec_path, test_path); err != nil {
		common.ShowError(err)
		os.Exit(1)
	}
	fmt.Printf("Generated unit tests into: %s\n", test_path)
}

// Fetches TEST_SRC and places it into `outpath`
func downloadSource(outPath string) error {
	var err error

	var out *os.File
	if out, err = os.Create(outPath); err != nil {
		return err
	}
	defer out.Close()

	var resp *http.Response
	if resp, err = http.Get(TESTS_SRC); err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

// Takes a specification yaml and generates a go unit test file to outpath
func generateTestFile(specPath string, outPath string) error {
	var err error

	// Load the test spec file into memory
	var specData []byte
	if specData, err = os.ReadFile(specPath); err != nil {
		return err
	}

	// Parse the tests into `spec`
	var spec Spec
	if err = yaml.Unmarshal(specData, &spec); err != nil {
		return err
	}

	// Generate the output string
	var sb strings.Builder
	sb.WriteString(template)
	testNames := make([]string, 0, len(spec.Tests))
	for k := range spec.Tests {
		testNames = append(testNames, k)
	}
	sort.Sort(sort.StringSlice(testNames))

	for _, testName := range testNames {
		testData := spec.Tests[testName]
		sb.WriteString(generateTestFunc(testName, &testData))
		sb.WriteString("\n")
	}

	// Create the specified file
	var out *os.File
	if out, err = os.Create(outPath); err != nil {
		return err
	}
	defer out.Close()

	// Write to new file
	if _, err = out.WriteString(sb.String()); err != nil {
		return err
	}

	return nil
}

// Builds a string definition of a test as Go code
func generateTestFunc(testName string, testData *Test) string {
	var sb strings.Builder
	r := testToRecipe(testName, testData)

	// Function signature
	sb.WriteString(fmt.Sprintf(
		"func Test%s(t *testing.T) {\n",
		strings.TrimPrefix(testName, "test")))

	// Call parser
	sb.WriteString(fmt.Sprintf(
		`	got := ParseRecipeString("", %s%s%s)
`, "`", testData.Source, "`"))

	// Define truth value
	sb.WriteString(fmt.Sprintf(
		`	want := %s
`, recipeToStrDef(r)))

	// Insert assertion and close function
	sb.WriteString(fmt.Sprintf(
		`	assertCanonicalRecipe(t, &got, &want)
}`))

	return sb.String()
}

// Converts `Test` to a `cook.Recipe`. `name` is used for error reporting.
func testToRecipe(name string, test *Test) cook.Recipe {
	r := cook.Recipe{
		Name:        "",
		Metadata:    map[string]string{},
		Ingredients: []cook.Ingredient{},
		Cookware:    []cook.Cookware{},
		Timers:      []cook.Timer{},
		Steps:       []cook.Step{},
	}

	r.Metadata = test.Result.Metadata

	for i, tStep := range test.Result.Steps {
		r.Steps = append(r.Steps, cook.Step{})
		for _, tChunk := range tStep {
			switch tChunk.Type {
			case "text":
				r.Steps[i] = append(r.Steps[i], cook.Text(tChunk.Value))
			case "ingredient":
				ingr := cook.Ingredient{
					Name:   tChunk.Name,
					Qty:    tChunk.Quantity,
					QtyVal: cook.TryParseQty(tChunk.Quantity),
					Unit:   tChunk.Units,
				}
				r.Ingredients = append(r.Ingredients, ingr)
				r.Steps[i] = append(r.Steps[i], ingr)
			case "cookware":
				cookware := cook.Cookware{
					Name:   tChunk.Name,
					Qty:    tChunk.Quantity,
					QtyVal: cook.TryParseQty(tChunk.Quantity),
					Unit:   tChunk.Units,
				}
				r.Cookware = append(r.Cookware, cookware)
				r.Steps[i] = append(r.Steps[i], cookware)
			case "timer":
				timer := cook.Timer{
					Name:   tChunk.Name,
					Qty:    tChunk.Quantity,
					QtyVal: cook.TryParseQty(tChunk.Quantity),
					Unit:   tChunk.Units,
				}
				r.Timers = append(r.Timers, timer)
				r.Steps[i] = append(r.Steps[i], timer)
			default:
				fmt.Fprintf(os.Stderr,
					"Invalid chunk in test %s: %+v\nAborting.",
					name, tChunk)
				os.Exit(1)
			}
		}
	}
	return r
}

// Creates a string of a Recipe being defined in Go code.
func recipeToStrDef(r cook.Recipe) string {
	var sb strings.Builder
	sb.WriteString(`Recipe{Name:"",Metadata:map[string]string{`)
	// Populate Metadata
	for k, v := range r.Metadata {
		str := fmt.Sprintf(`"%v":"%v",`, k, v)
		sb.WriteString(str)
	}
	sb.WriteString(`},Ingredients:[]Ingredient{`)
	//Populate Ingredients
	for _, ingr := range r.Ingredients {
		sb.WriteString(componentToStrDef(cook.Component(ingr)) + ",")
	}
	sb.WriteString(`},Cookware:[]Cookware{`)
	//Populate Cookware
	for _, ware := range r.Cookware {
		sb.WriteString(componentToStrDef(cook.Component(ware)) + ",")
	}
	sb.WriteString(`},Timers:[]Timer{`)
	//Populate Timers
	for _, timer := range r.Timers {
		sb.WriteString(componentToStrDef(cook.Component(timer)) + ",")
	}
	sb.WriteString(`},Steps:[]Step{`)
	for _, step := range r.Steps {
		sb.WriteString(`{`)
		for _, step_chunk := range step {
			var str string
			switch chunk := step_chunk.(type) {
			case cook.Text:
				str = fmt.Sprintf(`Text("%s"),`, chunk.ToString())
			case cook.Ingredient:
				str = fmt.Sprintf(
					`Ingredient%s,`,
					componentToStrDef(cook.Component(chunk)))
			case cook.Cookware:
				str = fmt.Sprintf(
					`Cookware%s,`,
					componentToStrDef(cook.Component(chunk)))
			case cook.Timer:
				str = fmt.Sprintf(
					`Timer%s,`,
					componentToStrDef(cook.Component(chunk)))
			default:
				fmt.Fprintf(
					os.Stderr,
					"Unhandled chunk type in steps during recipeToStrDef.\n")
				os.Exit(1)
			}
			sb.WriteString(str)
		}
		sb.WriteString("},")
	}
	sb.WriteString("},}")
	return sb.String()
}

// Builds the string definition of a component in Go code.
func componentToStrDef(c cook.Component) string {
	qtyVal := "NoQty"
	if c.QtyVal != cook.NoQty {
		qtyVal = fmt.Sprintf("%v", c.QtyVal)
	}

	return fmt.Sprintf(
		`{Name:"%s",Qty:"%s",QtyVal:%s,Unit:"%s"}`,
		c.Name, c.Qty, qtyVal, c.Unit)
}
