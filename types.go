package cook

import (
	"encoding/json"
)

// Chunks are the building blocks of recipe steps.
// They are a union of Text, Ingredient, Timer and Cookware.
type Chunk interface {
	isChunk()
	ToString() string
}

func (Text) isChunk()       {}
func (Ingredient) isChunk() {}
func (Cookware) isChunk()   {}
func (Timer) isChunk()      {}

// Text is a string wrapper (to allow for safe inclusion to Chunk interface)
type Text string	
type Ingredient component
type Cookware component
type Timer component

// Unwraps Text chunk into a string
func (x Text) ToString() string {
	return string(x)
}

// Converts an ingredient to a string, with its name followed by qty and units if they
// exist.
func (x Ingredient) ToString() string {
	return x.Name
}

// Converts cookware to a string, with its name followed by qty and units if they
// exist.
func (x Cookware) ToString() string {
	return x.Name
}

// Converts a timer to a string, with its name followed by qty and units if they
// exist.
func (x Timer) ToString() string {
	return x.Name
}

// A Step is one part of a recipe, consisting of a set of chunks which can be
// read in order to build a human readable recipe.
//
// Ingredients, Cookware and Timer are kept as structs to allow for post processing
// (such as text formatting).
type Step []Chunk

// Custom JSON encoding wraps each of `Step`'s chunk into a struct that stores the 
// type to allow for unambiguous decoding.
//
// e.g. an ingredient chunk will be wrapped as encoded as:
// `{'tag': 'ingredient', 'data': {...}}`
func (s *Step) MarshalJSON() ([]byte, error){
	type wrapper struct {
		Tag 	string	`json:"tag"`	// The underlying type of a Chunk
		Data 	Chunk	`json:"data"`	// The actual chunk data
	}
	
	// Construct a new list of wrapped chunks
	wrappedSteps := make([]wrapper, len(*s))
	for i, chunk := range *s {
		var tag string
		// Can't type switch into a hole (e.g. switch _ = chunk.(type))
		// because Go Devs are supreme beings
		var fixYourDamnCompilerWarnings Chunk
		switch opinionatedLanguageDevs := chunk.(type) {
		case Text:
			tag = "text"
		case Ingredient:
			tag = "ingredient"
		case Cookware:
			tag = "cookware"
		case Timer:
			tag = "timer"
		default:
			fixYourDamnCompilerWarnings = opinionatedLanguageDevs
			panic("Tried to encode unhandled Chunk type")
		}
		var _ = fixYourDamnCompilerWarnings
		
		wrapped := wrapper{
			Tag: 	tag,
			Data:   chunk,
		}

		wrappedSteps[i] = wrapped
	}
	
	// Encode the wrapped chunks as JSON
	return json.Marshal(wrappedSteps)
}

// Decodes the custom wrapped chunks created by the encoder back into an array of `Chunk`s.
func (s *Step) UnmarshalJSON(data []byte) error {
	// Unwrap the list
	var out []interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	// Build an output list to populate
	step := make(Step, len(out))
	
	// Iterate over each wrapper chunk in the JSON
	for i, chunkMapRaw := range out {
		// chunkWrap is of form: {'tag':..., 'data':...}
		chunkWrap := chunkMapRaw.(map[string]interface{})
		tag := chunkWrap["tag"].(string)
		// Handle text separate, as it requires no extra parsing
		if tag == "text" {
			step[i] = Text(chunkWrap["data"].(string))
		} else {
			// Re-encode chunks to json to unmarshal as a component
			dataMap := chunkWrap["data"].(map[string]interface{})
			data, _ := json.Marshal(dataMap)
			var comp component
			_ = json.Unmarshal(data, &comp)

			// Convert component to relevant type
			switch chunkWrap["tag"].(string) {
			case "ingredient":
				step[i] = comp.toIngredient()
			case "cookware":
				step[i] = comp.toCookware()
			case "timer":
				step[i] = comp.toTimer()
			default:
				panic("Encountered unhandled tag on JSON decode of `step`")
			}
		}	
	}
	// Push array to *s and return that there was no error
	*s = step
	return nil
}


// Metadata is arbitrary information about a recipe, consisting
// simply of a tag and a body.
//
// While some implementations only allow metadata at the start, the spec
// defines it in equal precedence to a `Step`.
//
// It may be prudent to further parse metadata before displaying.
type Metadata struct {
	Tag  string		`json:"tag"`
	Body string		`json:"body"`
}

// Recipes consist primarily of Metadata and Steps. Steps are stored sequentially
// and offer a continuous construction of the parsed `.cook` file.
//
// Additionally, the Ingredients, Cookware and Timer members provide a manifest
// for each of the respective item classes.
//
// Recipes can be easily parsed from a string using the function `ParseRecipe`.
type Recipe struct {
	Name        string			`json:"name"`
	Metadata    []Metadata		`json:"metadata"`
	Ingredients []Ingredient	`json:"ingredients"`
	Cookware    []Cookware		`json:"cookware"`
	Timers      []Timer			`json:"timers"`
	Steps       []Step			`json:"steps"`
}

// Represents a generic `component`, used in cooklang to define
// ingredients, cookware and timers.
type component struct {
	Name   string	`json:"name"`
	Qty    string	`json:"qty"`
	QtyVal float64	`json:"qtyVal"`
	Unit   string	`json:"unit"`
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
