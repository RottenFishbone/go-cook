package config

// Note: The `Dir` fields are overwritten using an intelligent lookup
// see `ConfigInit()` in `config.go`
var default_config = Config{
	Recipe: RecipeConfig{
		Dir:   "",
		Units: "metric",
	},
	Shopping: ShoppingConfig{
		Dir: "",
	},
}
