package main

import (
	"RottenFishbone/cooklang/pkg/recipe"
	"fmt"
)

func main() {
    
    r := recipe.ParseRecipeString("Test Recipe", 
    `>> Author: Jayden
>> servings: 2

Preheat #deep fryer{} to 350°F. 
Cut @potatoes{3} into thin 1/4 slices. 
Optionally, blanch the cut potatoes in a #pot of boiling water for ~{4%minutes}.
Pat potatos dry and lower them into the deep fryer for ~frying{10%minutes} or until golden brown.
Season generously with #salt and serve with #ketchup.

-- Don't do this.
Consider adding @mayonnaise{equal parts} to the ketchup to make some fancy sauce.`)
    fmt.Printf("r: %v\n", r)
}
