package cooklang

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

// --------------------------------------------------------------
// Unit Tests
// --------------------------------------------------------------

func TestTryParseFraction(t *testing.T) {
	// Function to easily test inputs
	testFrac := func(in string, want float64) {
		v := tryParseFraction(in)
		if (math.IsNaN(want) && !math.IsNaN(v)) || v != want {
			t.Fatalf("Failed to parse: \"%s\"\ngot: %v\nwant: %v.", in, v, want)
		}
	}

	// Should Succeed
	testFrac("1/2", 0.5)         // Larger b (decimal result)
	testFrac("2/1", 2.0)         // Larger a
	testFrac("10/10", 1.0)       // Two digit
	testFrac("500/1000", 0.5)    // Larger numbers
	testFrac("1.5", 1.5)         // Decimal values
	testFrac("100.084", 100.084) // Large Decimal values
	testFrac("0.084", 0.084)     // Values < 1 (leading 0)
	testFrac(".084", 0.084)      // Values < 1 (no leading #)
	testFrac("5", 5.0)           // Positive Integers
	testFrac("840", 840.0)       // Positive Large Integers

	// Should Fail
	inf := math.Inf(-1)
	testFrac("0/1", inf)   // a Only 0
	testFrac("1/0", inf)   // b Only 0
	testFrac("01/10", inf) // a Leading 0
	testFrac("10/01", inf) // b Leading 0
	testFrac("01.0", inf)  // Decimal leading with 0
	testFrac("1.0/1", inf) // Decimal with fraction
	testFrac("-1", inf)    // Negative int
	testFrac("-1.0", inf)  // Negative decimal
	testFrac("inf", inf)   // Float keyword
}

// --------------------------------------------------------------
// Examples
// --------------------------------------------------------------

func ExampleParseRecipeString() {
	recipeText :=
		`Preheat #deep fryer{} to 190°C.
Slice @potatoes{3} into 1/4" strips.
Optionally, blanch in boiling @water{2%cups}.
Drop into deep fryer for ~{7%mins}.
Remove from fryer and sprinkle @pink salt{} 
Enjoy with @ketchup, or mix in @mayonnaise{equal parts} for fancy sauce.`

	recipe := ParseRecipeString("Fries", recipeText)
	fmt.Println(recipe.Ingredients)
	// Output: [{potatoes 3 3 } {water 2 2 cups} {pink salt  -Inf } {ketchup  -Inf } {mayonnaise equal parts -Inf }]

}

func ExampleParseRecipe() {
	recipeText :=
		`Preheat #deep fryer{} to 190°C.
Slice @potatoes{3} into 1/4" strips.
Optionally, blanch in boiling @water{2%cups}.
Drop into deep fryer for ~{7%mins}.
Remove from fryer and sprinkle @pink salt{} 
Enjoy with @ketchup, or mix in @mayonnaise{equal parts} for fancy sauce.`

	data := []byte(recipeText)
	recipe := ParseRecipe("Fries", &data)
	fmt.Println(recipe.Ingredients)
	// Output: [{potatoes 3 3 } {water 2 2 cups} {pink salt  -Inf } {ketchup  -Inf } {mayonnaise equal parts -Inf }]

}

// --------------------------------------------------------------
// Canonical Unit Tests
//
// NOTE: I did ignore `Qty: "some"` for one word ingredients and `Qty: 1` for one word cookware
// It's fairly trivial to add, however, it is a little weird for some ingredients and some
// cookware i.e. "some egg" or "1 tongs".
// It's also presumptious concerning non-english languages.
//
// I'm going to leave it to the API interfacing this parser to work that out as seen fit.
//
// As defined here:
// https://github.com/cooklang/spec/tree/fa9bc51515b3317da434cb2b5a4a6ac12257e60b/tests
// --------------------------------------------------------------

// Deep compares two recipes, emitting t.Fatalf on inequality.
func assertRecipe(t *testing.T, got *Recipe, want *Recipe) {
	if !reflect.DeepEqual(*want, *got) {
		t.Fatalf("Assertion failed:\ngot:\t%+v\nwant:\t%+v", *got, *want)
	}
}

//-------------------------------------------------------------

func TestBasicDirection(t *testing.T) {
	got := ParseRecipeString("", "Add a bit of chilli")
	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Add a bit of chilli")}},
	}

	assertRecipe(t, &got, &want)
}

func TestComments(t *testing.T) {
	got := ParseRecipeString("", "-- testing comments")
	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{},
	}
	assertRecipe(t, &got, &want)
}

func TestCommentsWithIngredients(t *testing.T) {
	got := ParseRecipeString("",
		`-- testing comments
@thyme{2%sprigs}`)

	thyme := Ingredient{
		Name:   "thyme",
		Qty:    "2",
		QtyVal: 2.0,
		Unit:   "sprigs",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{thyme},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{thyme}},
	}

	assertRecipe(t, &got, &want)
}

func TestDirectionWithDegrees(t *testing.T) {
	got := ParseRecipeString("", "Heat oven up to 200°C")
	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Heat oven up to 200°C")}},
	}

	assertRecipe(t, &got, &want)
}

func TestDirectionWithNumbers(t *testing.T) {
	got := ParseRecipeString("", "Heat 5L of water")
	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Heat 5L of water")}},
	}

	assertRecipe(t, &got, &want)
}

func TestDirectionWithIngredient(t *testing.T) {
	got := ParseRecipeString("", "Add @chilli{3%items}, @ginger{10%g} and @milk{1%l}.")
	chili := Ingredient{
		Name:   "chilli",
		Qty:    "3",
		QtyVal: 3.0,
		Unit:   "items",
	}
	ginger := Ingredient{
		Name:   "ginger",
		Qty:    "10",
		QtyVal: 10.0,
		Unit:   "g",
	}
	milk := Ingredient{
		Name:   "milk",
		Qty:    "1",
		QtyVal: 1.0,
		Unit:   "l",
	}
	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chili, ginger, milk},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps: []Step{
			{Text("Add "), chili, Text(", "),
				ginger, Text(" and "), milk, Text(".")},
		},
	}

	assertRecipe(t, &got, &want)
}

func TestEquipmentMultipleWords(t *testing.T) {
	got := ParseRecipeString("", "Fry in #frying pan{}")

	pan := Cookware{
		Name:   "frying pan",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{pan},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Fry in "), pan}},
	}

	assertRecipe(t, &got, &want)
}

func TestEquipmentMultipleWordsWithLeadingNumber(t *testing.T) {
	got := ParseRecipeString("", "Fry in #7-inch nonstick frying pan{}")

	pan := Cookware{
		Name:   "7-inch nonstick frying pan",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{pan},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Fry in "), pan}},
	}

	assertRecipe(t, &got, &want)
}

func TestEquipmentOneWord(t *testing.T) {
	got := ParseRecipeString("", "Fry in #pan for some time")

	pan := Cookware{
		Name:   "pan",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{pan},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Fry in "), pan, Text(" for some time")}},
	}

	assertRecipe(t, &got, &want)
}

func TestEquipmentQuantity(t *testing.T) {
	got := ParseRecipeString("", "#frying pan{2}")

	pan := Cookware{
		Name:   "frying pan",
		Qty:    "2",
		QtyVal: 2.0,
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{pan},
		Timer:       []Timer{},
		Steps:       []Step{{pan}},
	}

	assertRecipe(t, &got, &want)
}

func TestEquipmentQuantityOneWord(t *testing.T) {
	got := ParseRecipeString("", "#frying pan{three}")

	pan := Cookware{
		Name:   "frying pan",
		Qty:    "three",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{pan},
		Timer:       []Timer{},
		Steps:       []Step{{pan}},
	}

	assertRecipe(t, &got, &want)
}

func TestEquipmentQuantityMultipleWords(t *testing.T) {
	got := ParseRecipeString("", "#frying pan{two small}")

	pan := Cookware{
		Name:   "frying pan",
		Qty:    "two small",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{pan},
		Timer:       []Timer{},
		Steps:       []Step{{pan}},
	}

	assertRecipe(t, &got, &want)
}

func TestFractions(t *testing.T) {
	got := ParseRecipeString("", "@milk{1/2%cup}")

	milk := Ingredient{
		Name:   "milk",
		Qty:    "1/2",
		QtyVal: 0.5,
		Unit:   "cup",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{milk},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{milk}},
	}

	assertRecipe(t, &got, &want)
}

func TestFractionsInDirections(t *testing.T) {
	got := ParseRecipeString("", "knife cut about every 1/2 inches")

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("knife cut about every 1/2 inches")}},
	}

	assertRecipe(t, &got, &want)
}

func TestFractionsLike(t *testing.T) {
	got := ParseRecipeString("", "@milk{01/2%cup}")

	milk := Ingredient{
		Name:   "milk",
		Qty:    "01/2",
		QtyVal: math.Inf(-1),
		Unit:   "cup",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{milk},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{milk}},
	}

	assertRecipe(t, &got, &want)
}

func TestFractionsWithSpaces(t *testing.T) {
	got := ParseRecipeString("", "@milk{1 / 2%cup}")

	milk := Ingredient{
		Name:   "milk",
		Qty:    "1 / 2",
		QtyVal: 0.5,
		Unit:   "cup",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{milk},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{milk}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientMultipleWords(t *testing.T) {
	got := ParseRecipeString("", "@hot chilli{3}")

	chilli := Ingredient{
		Name:   "hot chilli",
		Qty:    "3",
		QtyVal: 3.0,
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chilli},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{chilli}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientMultipleWordsNoAmount(t *testing.T) {
	got := ParseRecipeString("", "@hot chilli{}")

	chilli := Ingredient{
		Name:   "hot chilli",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chilli},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{chilli}},
	}

	assertRecipe(t, &got, &want)
}

func TestMultipleIngredientWithoutStopper(t *testing.T) {
	got := ParseRecipeString("", "@chilli cut into pieces and @garlic")

	chilli := Ingredient{
		Name:   "chilli",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	garlic := Ingredient{
		Name:   "garlic",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chilli, garlic},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{chilli, Text(" cut into pieces and "), garlic}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientMultipleWordsWithLeadingNumber(t *testing.T) {
	got := ParseRecipeString("", "Top with @1000 island dressing{ }")

	dressing := Ingredient{
		Name:   "1000 island dressing",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{dressing},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Top with "), dressing}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientWithEmoji(t *testing.T) {
	got := ParseRecipeString("", "Add some @🧂")

	salt := Ingredient{
		Name:   "🧂",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{salt},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Add some "), salt}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientExplicitUnitsWithSpaces(t *testing.T) {
	got := ParseRecipeString("", "@chilli{ 3 % items }")

	chilli := Ingredient{
		Name:   "chilli",
		Qty:    "3",
		QtyVal: 3.0,
		Unit:   "items",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chilli},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{chilli}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientImplicitUnits(t *testing.T) {
	got := ParseRecipeString("", "@chilli{}")

	chilli := Ingredient{
		Name:   "chilli",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chilli},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{chilli}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientNoUnits(t *testing.T) {
	got := ParseRecipeString("", "@chilli")

	chilli := Ingredient{
		Name:   "chilli",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chilli},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{chilli}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientNoUnitsNotOnlyString(t *testing.T) {
	got := ParseRecipeString("", "@5peppers")

	peppers := Ingredient{
		Name:   "5peppers",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{peppers},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{peppers}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientWithNumbers(t *testing.T) {
	got := ParseRecipeString("", "@tipo 00 flour{250%g}")

	tipo := Ingredient{
		Name:   "tipo 00 flour",
		Qty:    "250",
		QtyVal: 250.0,
		Unit:   "g",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{tipo},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{tipo}},
	}

	assertRecipe(t, &got, &want)
}

func TestIngredientWithoutStopper(t *testing.T) {
	got := ParseRecipeString("", "@chilli cut into pieces")

	chilli := Ingredient{
		Name:   "chilli",
		Qty:    "",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{chilli},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{chilli, Text(" cut into pieces")}},
	}

	assertRecipe(t, &got, &want)
}

func TestMetadata(t *testing.T) {
	got := ParseRecipeString("", ">> sourced: babooshka")

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{{Tag: "sourced", Body: "babooshka"}},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{},
	}

	assertRecipe(t, &got, &want)
}

func TestMetadataBreak(t *testing.T) {
	got := ParseRecipeString("", "hello >> sourced: babooshka")

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("hello >> sourced: babooshka")}},
	}

	assertRecipe(t, &got, &want)
}

func TestMetadataMultiwordKey(t *testing.T) {
	got := ParseRecipeString("", ">> cooking time: 30 mins")

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{{Tag: "cooking time", Body: "30 mins"}},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{},
	}

	assertRecipe(t, &got, &want)
}

func TestMetadataMultiwordKeyWithSpaces(t *testing.T) {
	got := ParseRecipeString("", ">>cooking time    :30 mins")

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{{Tag: "cooking time", Body: "30 mins"}},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{},
	}

	assertRecipe(t, &got, &want)
}

func TestMetadataServings(t *testing.T) {
	got := ParseRecipeString("", ">> servings: 1|2|3")

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{{Tag: "servings", Body: "1|2|3"}},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{},
	}

	assertRecipe(t, &got, &want)
}

func TestMetadataMultipleLines(t *testing.T) {
	got := ParseRecipeString("",
		`>> Prep Time: 15 minutes
>> Cook Time: 30 minutes`)

	want := Recipe{
		Name: "",
		Metadata: []Metadata{
			{Tag: "Prep Time", Body: "15 minutes"},
			{Tag: "Cook Time", Body: "30 minutes"},
		},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{},
	}

	assertRecipe(t, &got, &want)
}

func TestMultilineDirections(t *testing.T) {
	got := ParseRecipeString("",
		`Add a bit of chilli

Add a bit of hummus`)

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps: []Step{
			{Text("Add a bit of chilli")},
			{Text("Add a bit of hummus")},
		},
	}

	assertRecipe(t, &got, &want)
}

func TestQuantityAsText(t *testing.T) {
	got := ParseRecipeString("", "@thyme{few%sprigs}")

	thyme := Ingredient{
		Name:   "thyme",
		Qty:    "few",
		QtyVal: math.Inf(-1),
		Unit:   "sprigs",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{thyme},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{thyme}},
	}

	assertRecipe(t, &got, &want)
}

func TestQuantityDigitalString(t *testing.T) {
	got := ParseRecipeString("", "@water{7 k }")

	water := Ingredient{
		Name:   "water",
		Qty:    "7 k",
		QtyVal: math.Inf(-1),
		Unit:   "",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{water},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{water}},
	}

	assertRecipe(t, &got, &want)
}

func TestSlashInText(t *testing.T) {
	got := ParseRecipeString("", "Preheat the oven to 200℃/Fan 180°C")

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{},
		Steps:       []Step{{Text("Preheat the oven to 200℃/Fan 180°C")}},
	}

	assertRecipe(t, &got, &want)
}

func TestTimerDecimal(t *testing.T) {
	got := ParseRecipeString("", "Fry for ~{1.5%minutes}")

	timer := Timer{
		Name:   "",
		Qty:    "1.5",
		QtyVal: 1.5,
		Unit:   "minutes",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{timer},
		Steps:       []Step{{Text("Fry for "), timer}},
	}

	assertRecipe(t, &got, &want)
}

func TestTimerFractional(t *testing.T) {
	got := ParseRecipeString("", "Fry for ~{1/2%hour}")

	timer := Timer{
		Name:   "",
		Qty:    "1/2",
		QtyVal: 0.5,
		Unit:   "hour",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{timer},
		Steps:       []Step{{Text("Fry for "), timer}},
	}

	assertRecipe(t, &got, &want)
}

func TestTimerInteger(t *testing.T) {
	got := ParseRecipeString("", "Fry for ~{10%minutes}")

	timer := Timer{
		Name:   "",
		Qty:    "10",
		QtyVal: 10.0,
		Unit:   "minutes",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{timer},
		Steps:       []Step{{Text("Fry for "), timer}},
	}

	assertRecipe(t, &got, &want)
}

func TestTimerWithName(t *testing.T) {
	got := ParseRecipeString("", "Fry for ~potato{42%minutes}")

	timer := Timer{
		Name:   "potato",
		Qty:    "42",
		QtyVal: 42.0,
		Unit:   "minutes",
	}

	want := Recipe{
		Name:        "",
		Metadata:    []Metadata{},
		Ingredients: []Ingredient{},
		Cookware:    []Cookware{},
		Timer:       []Timer{timer},
		Steps:       []Step{{Text("Fry for "), timer}},
	}

	assertRecipe(t, &got, &want)
}
