package config

// Constant key for reliably accessing `config.Vars`
const (
	RecipeDir   string = "RECIPE_DIR"
	ShoppingDir        = "SHOPPING_DIR"
	Units              = "UNITS"
)

// Globally accessible dictionary of config values.
var Vars = map[string]string{
	RecipeDir:   "",
	ShoppingDir: "",
	Units:       "",
}

// Convert the config dictionary into a `config` struct.
// This is used exclusively for (de)serialization.
func varsToConfig() config {
	return config{
		Recipe: recipeConfig{
			Dir: Vars[RecipeDir],
		},
		Shopping: shoppingConfig{
			Dir: Vars[ShoppingDir],
		},
		Units: Vars[Units],
	}
}

// Loads a `config` struct into the config dictionary.
// All values will be overwritten.
func configToVars(cfg *config) {
	Vars[RecipeDir] = cfg.Recipe.Dir
	Vars[ShoppingDir] = cfg.Shopping.Dir
	Vars[Units] = cfg.Units
}
