package config

var Vars = map[string]string{
	"RECIPE_DIR":   "",
	"SHOPPING_DIR": "",
	"UNITS":        "",
}

func varsToConfig() config {
	return config{
		Recipe: recipeConfig{
			Dir: Vars["RECIPE_DIR"],
		},
		Shopping: shoppingConfig{
			Dir: Vars["SHOPPING_DIR"],
		},
		Units: Vars["UNITS"],
	}
}

func configToVars(cfg *config) {
    Vars["RECIPE_DIR"] = cfg.Recipe.Dir
    Vars["SHOPPING_DIR"]= cfg.Shopping.Dir
    Vars["UNITS"] = cfg.Units
}
