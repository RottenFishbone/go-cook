package config

// TODO Platform independent configs

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type Config struct {
	Recipe   RecipeConfig   `toml:"recipe"`
	Shopping ShoppingConfig `toml:"shopping"`
}

type RecipeConfig struct {
	Dir   string `toml:"dir"`
	Units string `toml:"units"`
}

type ShoppingConfig struct {
	Dir string `toml:"dir"`
}

// Loads a `cooklang-go` config file and returns the parsed `Config` struct.
//
// Leave path blank to use default location.
// NOTE: This only works with Unix based OS atm
func LoadConfig(path string) (Config, bool) {
	// Find default path if needed
	if path == "" {
		path = findConfigPath()
	}

	// Ensure the target file actually exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, false
	}

	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg, true
}

// Creates a new config file, using the default template, using `path`.
//
// If `recipes` or `shopping` are not empty, the passed string will be passed
// on to the new config file as locations for each. Otherwise, XDG_DATA_HOME defaults
// will be used.
//
// If the recipe and shopping dirs don't already exist, they will be created.
func ConfigInit(path string, recipes string, shopping string) Config {
	if path == "" {
		path = findConfigPath()
	}

	// Create new file if it doesn't already exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			panic("Failed to create config directory at: " + filepath.Dir(path))
		}
	}

	// Try to create the file
	file, err := os.Create(path)
	if os.IsExist(err) {
		panic("Config already exists at: " + path)
	}
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Populate recipes/shopping if needed
	if recipes == "" {
		recipes = findRecipesPath()
	}
	if shopping == "" {
		shopping = findShoppingPath()
	}

	// Push them to default config
	cfg := default_config
	cfg.Recipe.Dir = recipes
	cfg.Shopping.Dir = shopping

	// Encode into config file
	err = toml.NewEncoder(file).Encode(cfg)
	if err != nil {
		panic(err)
	}

	//TODO Ensure the recipe/shopping dir exist and spawn otherwise

	return cfg
}

// Returns the default config filepath
func findConfigPath() string {
	var path string

	// Try env var first
	cookConfig, cookCfgExists := os.LookupEnv("COOK_CONFIG")
	if cookCfgExists {
		path = cookConfig
	}

	// Fallback XDG standards
	xdgConfig, xdgCfgExists := os.LookupEnv("XDG_CONFIG_HOME")
	if xdgCfgExists {
		path = filepath.Join(xdgConfig, "cook", "config.toml")
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".config", "cook", "config.toml")
	}

	return path
}

// Returns the default recipes directory path
func findRecipesPath() string {
	var path string

	// Use XDG standards
	xdgConfig, xdgCfgExists := os.LookupEnv("XDG_DATA_HOME")
	if xdgCfgExists {
		path = filepath.Join(xdgConfig, "cook", "recipes")
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".local", "share", "cook", "recipes")
	}

	return path
}

// Returns the default shopping directory path
func findShoppingPath() string {
	var path string

	// Use XDG standards
	xdgConfig, xdgCfgExists := os.LookupEnv("XDG_DATA_HOME")
	if xdgCfgExists {
		path = filepath.Join(xdgConfig, "cook", "shopping")
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".local", "share", "cook", "shopping")
	}

	return path
}
