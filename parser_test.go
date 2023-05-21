package cook

import (
	"fmt"
	"math"
	"testing"
)

// --------------------------------------------------------------
// Unit Tests
// --------------------------------------------------------------

func TestTryParseQty(t *testing.T) {
	// Function to easily test inputs
	testFrac := func(in string, want float64) {
		v := TryParseQty(in)
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
	testFrac("0/1", NoQty)   // a Only 0
	testFrac("1/0", NoQty)   // b Only 0
	testFrac("01/10", NoQty) // a Leading 0
	testFrac("10/01", NoQty) // b Leading 0
	testFrac("01.0", NoQty)  // Decimal leading with 0
	testFrac("1.0/1", NoQty) // Decimal with fraction
	testFrac("-1", NoQty)    // Negative int
	testFrac("-1.0", NoQty)  // Negative decimal
	testFrac("NoQty", NoQty) // Float keyword
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
	// Output: [{potatoes 3 3 } {water 2 2 cups} {pink salt  0 } {ketchup  0 } {mayonnaise equal parts 0 }]

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
	// Output: [{potatoes 3 3 } {water 2 2 cups} {pink salt  0 } {ketchup  0 } {mayonnaise equal parts 0 }]

}
