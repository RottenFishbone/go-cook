package cook

import (
	"regexp"
	"strconv"
	"strings"

	y "github.com/prataprc/goparsec"
)

// TODO Error reporting
// TODO Implement "servings" system a la cooklang roadmap
// TODO Create extended lang? or an extension system of sorts.
// TODO Parse image tags, This may be the responsibility of the renderer, though

// Used as an `ASTNodify` callback to ensure a node is named and added
func forceNamed(name string, s y.Scanner, node y.Queryable) y.Queryable {
	return &y.NonTerminal{Name: name, Children: []y.Queryable{node}}
}

// Builds a `goparsec` `Parser` using the provided `AST`.
// The parser will populate the AST according the the cooklang spec.
// Some points to note:
//   - comments are expected to have been stripped prior to parsing,
//
// [Cooklang spec](https://github.com/cooklang/spec/blob/fa9bc51515b3317da434cb2b5a4a6ac12257e60b/EBNF.md)
func buildCookY(ast *y.AST) y.Parser {
	//----------------
	// Terminal Parsers
	//----------------
	// Atoms
	tilde := y.AtomExact("~", "TILDE")
	at := y.AtomExact("@", "AT")
	hash := y.AtomExact("#", "HASH")
	ocurl := y.AtomExact("{", "OCURL")
	ccurl := y.AtomExact("}", "CCURL")
	percent := y.AtomExact("%", "PERCENT")
	meta := y.AtomExact(">>", "META")
	colon := y.AtomExact(":", "COLON")

	// Newline
	crlf := y.AtomExact("\r\n", "CRLF")
	lf := y.AtomExact("\n", "LF")
	cr := y.AtomExact("\r", "CR")
	uniNl := y.TokenExact(`[\x{000A}-\x{000D}\x{0085}\x{0028}\x{0029}]`, "UNICODE_NL")

	// Tokens
	specifierRegex := `[~@#]`
	specifier := y.TokenExact(specifierRegex, "SPEC")
	whitespace := y.TokenExact(`[\p{Zs}\x{0009}]`, "WHITESPACE")
	punctuation := y.TokenExact(`\pP`, "PUNCT")
	char := y.TokenExact(`.`, "CHAR")
	rawText := y.TokenExact(`.+`, "RAW")
	// Combinators
	nl := y.Many(nil, y.OrdChoice(nil, crlf, lf, cr, uniNl))
	//---------------

	//-------------
	// Text
	//-------------
	// all chacters EXCEPT `specifier`
	text := ast.ManyUntil("text", nil, char, nil, specifier)
	// consume a `specifier` then act as `text` normally does
	specText := ast.And("text", nil, specifier, text)
	// non-punctuation, non-white space text
	word := ast.ManyUntil("word", nil, char, nil, y.OrdChoice(nil, punctuation, whitespace))

	//-------------
	// Amount
	//-------------
	quantity := ast.ManyUntil("quantity", nil, char, nil,
		ast.OrdChoice("", nil, percent, ccurl))
	unit := ast.ManyUntil("unit", nil, char, nil, ccurl)
	quantityWithUnit := ast.And("quantity_with_unit", nil, quantity, percent, unit)
	amount := ast.OrdChoice("amount", nil, quantityWithUnit, quantity)
	optAmount := ast.Maybe("amount", nil, amount)

	//-------------
	// Components
	//-------------
	amountField := ast.And("amount_field", nil, ocurl, optAmount, ccurl)
	// No-word component
	nwComponent := amountField
	// One-Word component
	optAmountField := ast.Maybe("", nil, amountField)
	owComponent := ast.And("one_word_component", nil, word, optAmountField)
	// Multi-word component
	mwComponentText := ast.ManyUntil("words", nil, char, nil, ocurl)
	mwComponent := ast.And("multiword_component",
		func(name string, s y.Scanner, node y.Queryable) y.Queryable {
			// Lookahead to prevent a `specifier` within `words`
			words := node.GetChildren()[1].GetValue()
			match, _ := regexp.MatchString(`(.*`+specifierRegex+`.)`, words)
			if match {
				return nil
			}
			return node
		},
		word, mwComponentText, amountField)

	//-------------
	// Ingredients
	//-------------
	ingredientTypes := ast.OrdChoice("", nil, mwComponent, owComponent)
	ingredient := ast.And("ingredient", nil, at, ingredientTypes)

	//-------------
	// Cookware
	//-------------
	cookwareTypes := ast.OrdChoice("", nil, mwComponent, owComponent)
	cookware := ast.And("cookware", nil, hash, cookwareTypes)

	//-------------
	// Timers
	//-------------
	timerTypes := ast.OrdChoice("", nil, mwComponent, owComponent, nwComponent)
	timer := ast.And("timer", nil, tilde, timerTypes)

	//------------
	// Metadata
	//------------
	metaHeader := ast.ManyUntil("meta_header", nil, char, nil, colon)
	metadata := ast.And("metadata", nil, meta, metaHeader, colon, rawText)

	//------------
	// Step
	//------------
	// NOTE: chunk uses nodify callback `forceNamed` to ensure a named node is created
	chunk := ast.OrdChoice("chunk", forceNamed, ingredient, cookware, timer, text, specText)
	step := ast.Kleene("step", nil, chunk)

	// Either metadata or step
	recipeElem := ast.OrdChoice("element", nil, metadata, step)

	// Parse each line until EOF
	return ast.ManyUntil("steps", nil, recipeElem, nl, ast.End("EOF"))
}

// Strips comment blocks of the form `--<example>\n` and
// (possibly multiline) block comments bounded by `[-` and `-]` from
// a byte array and returns the result.
func stripComments(data *[]byte) []byte {
	regex, _ := regexp.Compile(`((--.*((\r\n)|(\n)|(\r)|$))|(\[-(.|\s)*-\]))`)
	return regex.ReplaceAll(*data, []byte("\n"))
}

// Attempt to parse a quantity fraction string into a float representation.
// Additionally, parses decimals/whole values. e.g. "1.5", "5"
// cook.NoQty is returned on failure
func TryParseQty(qty string) float64 {
	if qty == "" {
		return NoQty
	}

	// Match non-negative, non-leading 0 numbers (and "0.X" or ".X")
	numMatch, _ := regexp.MatchString(
		`^((0?\.[0-9]+)|([1-9][0-9]*(\.?[0-9]+)?))$`, qty)
	if numMatch {
		val, err := strconv.ParseFloat(qty, 64)
		if err != nil {
			panic("Failed to parse number string into float during quantity parsing")
		}

		return val
	}

	// Match "a/b" where a and b are numbers without leading 0
	r := regexp.MustCompile(`^([1-9][0-9]*)\s?\/\s?([1-9][0-9]*)$`)
	matches := r.FindAllStringSubmatch(qty, -1)

	if matches == nil {
		return NoQty
	}

	a, aErr := strconv.ParseFloat(matches[0][1], 64)
	b, bErr := strconv.ParseFloat(matches[0][2], 64)
	if aErr != nil || bErr != nil {
		panic("Failed to parse value into float during quantity fraction parsing.")
	}
	return a / b
}

// Parses an "amount" node into `(qty, qtyVal, unit)`.
func parseAmountNode(node y.Queryable) (string, float64, string) {
	qty := ""
	unit := ""

	if node.GetName() != "missing" {
		quantityNode := node.GetChildren()[1]
		if quantityNode.GetName() != "missing" {
			switch quantityNode.GetName() {
			case "quantity_with_unit":
				qtyChildren := quantityNode.GetChildren()
				qty = strings.TrimSpace(qtyChildren[0].GetValue())
				unit = strings.TrimSpace(qtyChildren[2].GetValue())
			case "quantity":
				qty = strings.TrimSpace(quantityNode.GetValue())
			default:
				panic(`Unhandled node within "amount" node.`)
			}
		}
	}

	qtyVal := TryParseQty(qty)

	return qty, qtyVal, unit
}

// Parses an `AST` "*_component" node into a `component` struct.
// These are part of the cooklang spec and are used to define
// ingredients, cookware and timers
func parseComponentNode(node y.Queryable) Component {
	var text string
	var amountNode y.Queryable

	children := node.GetChildren()
	switch node.GetName() {
	case "amount_field": // no_name_component
		text = ""
		amountNode = node
	case "one_word_component":
		text = node.GetChildren()[0].GetValue()
		amountNode = children[1]
	case "multiword_component":
		text = children[0].GetValue() + children[1].GetValue()
		amountNode = children[2]
	default:
		panic("Unknown node found while parsing component.")
	}

	qty, qtyVal, unit := parseAmountNode(amountNode)
	return Component{text, qty, qtyVal, unit}
}

// Parses an `AST` "chunk" node. These are the building blocks of recipes.
// While not explicitly defined in the cooklang spec, they are the union
// all the specified components of a `step`.
//
// As such, they contain either a Text, Ingredient, Cookware or Timer subnode
// which we can parse into a `Chunk` interface.
func parseChunkNode(node y.Queryable) Chunk {
	if node.GetName() != "chunk" {
		panic("Cannot parse non-chunk nodes.")
	}
	// Try basic text parsing
	subNode := node.GetChildren()[0]
	if subNode.GetName() == "text" {
		return Text(subNode.GetValue())
	}

	// Parse component-based chunks
	var chunk Chunk
	compNode := subNode.GetChildren()[1] // e.g. one_word_component
	component := parseComponentNode(compNode)
	switch subNode.GetName() {
	case "ingredient":
		chunk = component.toIngredient()
	case "cookware":
		chunk = component.toCookware()
	case "timer":
		chunk = component.toTimer()
	default:
		panic("Chunk node contained unexpected sub-node")
	}

	return chunk
}

func ParseRecipeString(name string, data string) Recipe {
	bytes := []byte(data)
	return ParseRecipe(name, &bytes)
}

// Parses a byte array containing a recipe following the cooklang specifications
// and returns as a `Recipe` struct
func ParseRecipe(name string, data *[]byte) Recipe {
	r := Recipe{
		Name:        name,
		Metadata:    map[string]string{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timers:      []Timer{},
		Steps:       []Step{},
	}

	// Don't parse empty recipe
	if len(*data) == 0 {
		return r
	}

	// Strip comments before parsing
	stripped := stripComments(data)
	s := y.NewScanner(stripped)

	// Parse into AST
	ast := y.NewAST("recipe", 1024)
	parser := buildCookY(ast)
	root, _ := ast.Parsewith(parser, s)
	if root == nil {
		return r
	}

	// Collect and iterate over each important node to build recipe
	ch := make(chan y.Queryable, 1024)
	go ast.Query("metadata,step", ch)
	// Hanlde each node returned from the query
	for node := range ch {
		// Split into metadata and step
		switch node.GetName() {
		case "metadata":
			// Metadata is super simple, just push to recipe
			children := node.GetChildren()
			tag := strings.TrimSpace(children[1].GetValue())
			val := strings.TrimSpace(children[3].GetValue())
			r.Metadata[tag] = val
		case "step":
			// Steps are built from chunks, we need to parse those
			step := make(Step, 0)
			stepSubNodes := node.GetChildren()
			for i, chunkNode := range stepSubNodes {
				chunk := parseChunkNode(chunkNode)

				switch chunk := chunk.(type) {
				case Ingredient:
					r.Ingredients = append(r.Ingredients, chunk)
				case Cookware:
					r.Cookware = append(r.Cookware, chunk)
				case Timer:
					r.Timers = append(r.Timers, chunk)
				case Text:
					// Join consecutive text blocks together
					if len(step) > 0 {
						switch lastStep := step[i-1].(type) {
						case Text:
							step[i-1] = Text(lastStep.ToString() + chunk.ToString())
							continue
						}
					}
				default:
					panic("Unhandled Chunk type.")
				}

				step = append(step, chunk)
			}
			// Push newly built step into the recipe
			if len(step) > 0 {
				r.Steps = append(r.Steps, step)
			}
		default:
			panic("Unhandled node returned from query.")
		}
	}

	return r
}
