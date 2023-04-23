package parser

// Chunks are the building blocks of recipe steps.
// They are a union of Text, Ingredient, Timer and Cookware.
type Chunk interface{ isChunk() }

func (Text) isChunk()       {}
func (Ingredient) isChunk() {}
func (Cookware) isChunk()   {}
func (Timer) isChunk()      {}

// Text is a string wrapper (to allow for safe inclusion to Chunk interface)
type Text string
type Ingredient component
type Cookware component
type Timer component

// A Step is one part of a recipe, consisting of a set of chunks which can be
// read in order to build a human readable recipe.
//
// Ingredients, Cookware and Timer are kept as structs to allow for post processing
// (such as text formatting).
type Step []Chunk

// Metadata is arbitrary information about a recipe, consisting
// simply of a tag and a body.
//
// While some implementations only allow metadata at the start, the spec
// defines it in equal precedence to a `Step`.
//
// It may be prudent to further parse metadata before displaying.
type Metadata struct {
	Tag  string
	Body string
}

// Recipes consist primarily of Metadata and Steps. Steps are stored sequentially
// and offer a continuous construction of the parsed `.cook` file.
//
// Additionally, the Ingredients, Cookware and Timer members provide a manifest
// for each of the respective item classes.
//
// Recipes can be easily parsed from a string using the function `ParseRecipe`.
type Recipe struct {
	Name        string
	Metadata    []Metadata
	Ingredients []Ingredient
	Cookware    []Cookware
	Timer       []Timer
	Steps       []Step
}

// Represents a generic `component`, used in cooklang to define
// ingredients, cookware and timers.
type component struct {
	Name   string
	Qty    string
	QtyVal float64
	Unit   string
}

// Build an `Ingredient` from a `component`
func (node *component) toIngredient() Ingredient {
	return Ingredient{
		Name:   node.Name,
		Qty:    node.Qty,
		QtyVal: node.QtyVal,
		Unit:   node.Unit,
	}
}

// Build a `Cookware` from a `component`
func (node *component) toCookware() Cookware {
	return Cookware{
		Name:   node.Name,
		Qty:    node.Qty,
		QtyVal: node.QtyVal,
		Unit:   node.Unit,
	}
}

// Build a `Timer` from a `component`
func (node *component) toTimer() Timer {
	return Timer{
		Name:   node.Name,
		Qty:    node.Qty,
		QtyVal: node.QtyVal,
		Unit:   node.Unit,
	}
}
