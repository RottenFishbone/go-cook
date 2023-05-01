package api

// TODO cache recipes in memory

import (
)

// Returns a JSON list containing the `n`th page of recipes, using the
// provided page size.
func GetRecipesPage(n int, size int) string {
	return ""
}

// Returns a JSON list containing all recipes.
func GetAllRecipes() string {
	return ""
}

// Returns a JSON list containing all recipes with names that match using
// the provided regex pattern.
func GetRecipesMatching(regexPattern string) string {
	return ""
}

