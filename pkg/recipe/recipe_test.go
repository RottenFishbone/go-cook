package recipe

import (
	"fmt"
	"testing"

	"git.sr.ht/~rottenfishbone/go-cook"
)

// --------------------------------------------------------------
// Unit Tests
// --------------------------------------------------------------
func TestEncodeToJsonUninitdRecipe(t *testing.T) {
	r := cook.Recipe{}
	got := EncodeToJson(&r)
	want := `{"name":"","metadata":null,"ingredients":null,"cookware":null,"timers":null,"steps":null}`
	if got != want {
		t.Fatalf("Failed to encode: \"%v\"\ngot: %v\nwant: %v.", r, got, want)
	}
}

func TestEncodeToJsonEmptyRecipe(t *testing.T) {
	r := cook.Recipe{
		Name:        "",
		Metadata:    []cook.Metadata{},
		Ingredients: []cook.Ingredient{},
		Cookware:    []cook.Cookware{},
		Timers:      []cook.Timer{},
		Steps:       []cook.Step{},
	}
	got := EncodeToJson(&r)
	want := `{"name":"","metadata":[],"ingredients":[],"cookware":[],"timers":[],"steps":[]}`
	if got != want {
		t.Fatalf("Failed to encode: \"%v\"\ngot:\t%v\nwant:\t%v.", r, got, want)
	}
}

func TestEncodeToJsonRecipeMetadata(t *testing.T) {
	r := cook.Recipe{
		Name: "",
		Metadata: []cook.Metadata{{
			Tag:  "Author",
			Body: "Jayden",
		}},
		Ingredients: []cook.Ingredient{},
		Cookware:    []cook.Cookware{},
		Timers:      []cook.Timer{},
		Steps:       []cook.Step{},
	}
	got := EncodeToJson(&r)
	want := `{"name":"","metadata":[{"tag":"Author","body":"Jayden"}],"ingredients":[],"cookware":[],"timers":[],"steps":[]}`
	if got != want {
		t.Fatalf("Failed to encode: \"%v\"\ngot:\t%v\nwant:\t%v.", r, got, want)
	}
}

func TestEncodeToJsonRecipeIngredient(t *testing.T) {
	r := cook.Recipe{
		Name:     "",
		Metadata: []cook.Metadata{},
		Ingredients: []cook.Ingredient{{
			Name:   "tomato",
			Qty:    "1/2",
			QtyVal: 0.5,
			Unit:   "",
		}},
		Cookware: []cook.Cookware{},
		Timers:   []cook.Timer{},
		Steps:    []cook.Step{},
	}
	got := EncodeToJson(&r)
	want := `{"name":"","metadata":[],"ingredients":[{"name":"tomato","qty":"1/2","qtyVal":0.5,"unit":""}],"cookware":[],"timers":[],"steps":[]}`
	if got != want {
		t.Fatalf("Failed to encode: \"%v\"\ngot:\t%v\nwant:\t%v.", r, got, want)
	}
}

//TODO Cookware, Timer and Step tests + one full recipe
//TODO Decoder tests

// --------------------------------------------------------------
// Examples
// --------------------------------------------------------------
func ExampleEncodeToJson() {
	tomato := cook.Ingredient{
		Name:   "tomato",
		Qty:    "1",
		QtyVal: 1,
		Unit:   "",
	}
	r := cook.Recipe{
		Name:        "Test Recipe",
		Metadata:    []cook.Metadata{{Tag: "Author", Body: "Jayden"}},
		Ingredients: []cook.Ingredient{tomato},
		Cookware:    []cook.Cookware{},
		Timers:      []cook.Timer{},
		Steps: []cook.Step{
			{cook.Text("Slice whole "),
				tomato,
				cook.Text(" and eat fresh.")},
		},
	}

	json := EncodeToJson(&r)

	fmt.Println(json)
	// Output: {"name":"Test Recipe","metadata":[{"tag":"Author","body":"Jayden"}],"ingredients":[{"name":"tomato","qty":"1","qtyVal":1,"unit":""}],"cookware":[],"timers":[],"steps":[[{"tag":"text","data":"Slice whole "},{"tag":"ingredient","data":{"name":"tomato","qty":"1","qtyVal":1,"unit":""}},{"tag":"text","data":" and eat fresh."}]]}
}

func ExampleDecodeFromJson() {
	json := `{"name":"Test Recipe","metadata":[{"tag":"Author","body":"Jayden"}],"ingredients":[{"name":"tomato","qty":"1","qtyVal":1,"unit":""}],"cookware":[],"timers":[],"steps":[[{"tag":"text","data":"Slice whole "},{"tag":"ingredient","data":{"name":"tomato","qty":"1","qtyVal":1,"unit":""}},{"tag":"text","data":" and eat fresh."}]]}`

	r := DecodeFromJson(json)
	fmt.Println(r)
	// Output: {Test Recipe [{Author Jayden}] [{tomato 1 1 }] [] [] [[Slice whole  {tomato 1 1 }  and eat fresh.]]}
}
