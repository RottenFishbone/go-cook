package config

// Constant key for reliably accessing `config.Get`
const (
	KeyRecipeDir   string = "RECIPE_DIR"
	KeyShoppingDir        = "SHOPPING_DIR"
	KeyUnits              = "UNITS"
)

// Globally accessible dictionary of config values.
var vars = map[string]string{
	KeyRecipeDir:   "",
	KeyShoppingDir: "",
	KeyUnits:       "",
}

// If the config has been loaded from file or environment vars
var loaded = false

// Convert the config dictionary into a `config` struct.
// This is used exclusively for (de)serialization.
func varsToConfig() config {
	return config{
		Recipe: recipeConfig{
			Dir: vars[KeyRecipeDir],
		},
		Shopping: shoppingConfig{
			Dir: vars[KeyShoppingDir],
		},
		Units: vars[KeyUnits],
	}
}

// Loads a `config` struct into the config dictionary.
// All values will be overwritten.
func configToVars(cfg *config) {
	vars[KeyRecipeDir] = cfg.Recipe.Dir
	vars[KeyShoppingDir] = cfg.Shopping.Dir
	vars[KeyUnits] = cfg.Units
}
