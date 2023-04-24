package config

// TODO Platform independent configs

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Internal representations for configs, used to (de)serialize to toml
type (
	config struct {
		Recipe   recipeConfig   `toml:"recipe"`
		Shopping shoppingConfig `toml:"shopping"`
		Units    string         `toml:"units"`
	}

	recipeConfig struct {
		Dir string `toml:"dir"`
	}

	shoppingConfig struct {
		Dir string `toml:"dir"`
	}
)

// Loads a `cooklang-go` config file and returns the parsed `Config` struct.
//
// Leave path blank to use default location.
// NOTE: Defaults only work with Unix based OS atm
func LoadConfig(path string) bool {
	// Find default path if needed
	if path == "" {
		path = DefaultConfigPath()
	}

	// Ensure the target file actually exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	var cfg config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		panic(err)
	}

	configToVars(&cfg)

	return true
}

// Loads config from environment vars instead of a config file
func LoadConfigEnv() {
	for k := range Vars {
		v, _ := os.LookupEnv(k)
		Vars[k] = v
	}
}

// Creates a new config file, using the default template, using `path`.
//
// If `recipes` or `shopping` are not empty, the passed string will be passed
// on to the new config file as locations for each. Otherwise, XDG_DATA_HOME defaults
// will be used.
//
// If the recipe and shopping dirs don't already exist, they will be created.
func ConfigInit(path string, recipes string, shopping string) {
	if path == "" {
		path = DefaultConfigPath()
	}

	// Create new file if it doesn't already exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			panic("Failed to create config directory at: " + filepath.Dir(path))
		}
	} else {
		// Forbid accidental overwrites
		os.Stderr.WriteString("Config already exists at: " + path + ".. Exiting.\n")
		os.Exit(1)
	}

	// Try to create the file
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Populate recipes/shopping if needed
	if recipes == "" {
		recipes = DefaultRecipesPath()
	}
	if shopping == "" {
		shopping = DefaultShoppingPath()
	}

	// Push them to default config (built from default vars)
	cfg := varsToConfig()
	cfg.Recipe.Dir = recipes
	cfg.Shopping.Dir = shopping

	// Load the generated config back into Vars for use during execution
	configToVars(&cfg)

	// Encode into config file
	err = toml.NewEncoder(file).Encode(cfg)
	if err != nil {
		panic(err)
	}
}

// Returns the default config filepath
func DefaultConfigPath() string {
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
func DefaultRecipesPath() string {
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
func DefaultShoppingPath() string {
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
