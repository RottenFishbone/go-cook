package config

// TODO Platform independent configs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
	"github.com/BurntSushi/toml"
)

// A flag to dermine if the config has been loaded yet
var loaded = false

var configPath string
var conf Config

// Internal representations for configs, used to (de)serialize to toml
type (
	Config struct {
		Units    string         `toml:"units"`
		Recipe   RecipeConfig   `toml:"recipe"`
		Shopping ShoppingConfig `toml:"shopping"`
		Users    string         `toml:"users"`
		HMACKey  string         `toml:"hmac-key"`
	}

	RecipeConfig struct {
		Dir string `toml:"dir"`
	}

	ShoppingConfig struct {
		Dir string `toml:"dir"`
	}
)

// Returns a copy of the config, this should be
// fetched before each use, as it may change during runtime
// (dat global mutability baby)
//
// Panics if used before a load
func GetConfig() Config {
	if !loaded {
		panic("Attempted to read an unloaded config")
	}
	return conf
}

func IsLoaded() bool {
	return loaded
}

// Loads a `go-cook` config file and returns the parsed `Config` struct.
//
// Leave path blank to use default location.
// NOTE: Defaults only work with Unix based OS atm
func LoadConfig(path string) bool {
	// Find default path if needed
	if path == "" {
		path = DefaultConfigPath()
	}
	configPath = path

	// Ensure the target file actually exists
	if !common.FileExists(path) {
		return false
	}

	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		panic(err)
	}

	loaded = true

	return true
}

// Loads config from environment vars instead of a config file
// WARNING: Unimplemented
func LoadConfigEnv() {
	panic("TODO")
}

// Creates a new config file, using the default template, using `path`.
//
// If `recipes` or `shopping` are not empty, the passed string will be passed
// on to the new config file as locations for each. Otherwise, XDG_DATA_HOME defaults
// will be used.
//
// If the recipe and shopping dirs don't already exist, they will be created.
//
// `users.toml` will be added to config but not created until needed
func ConfigInit(path string, recipes string, shopping string) bool {
	var err error

	if path == "" {
		path = DefaultConfigPath()
	}

	// Create new file if it doesn't already exist
	if !common.FileExists(path) {
		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			errstr := fmt.Sprintf("Error: %v\n", err.Error())
			panic(errstr)
		}
	} else {
		// Forbid accidental overwrites
		return false
	}

	// Try to create the file
	var file *os.File
	if file, err = os.Create(path); err != nil {
		panic(err)
	}
	defer file.Close()

	// Populate recipes/shopping if needed
	if recipes == "" {
		recipes = defaultRecipesPath()
	}
	if shopping == "" {
		shopping = defaultShoppingPath()
	}

	// Load defaults into config, if needed
	conf.Recipe.Dir = recipes
	conf.Shopping.Dir = shopping
	conf.Users = defaultUsersPath()
	conf.HMACKey = string(generateHMACKey(128))

	loaded = true

	// Encode into config file and write
	if err = toml.NewEncoder(file).Encode(conf); err != nil {
		panic(err)
	}

	return true
}

// Returns the default config filepath
func DefaultConfigPath() string {
	// Try env var first
	cookConfig, cookCfgExists := os.LookupEnv("COOK_CONFIG")
	if cookCfgExists {
		return cookConfig
	}

	// Fallback XDG standards
	var path string
	xdgConfig, xdgCfgExists := os.LookupEnv("XDG_CONFIG_HOME")
	if xdgCfgExists {
		path = filepath.Join(xdgConfig, "cook", "config.toml")
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".config", "cook", "config.toml")
	}

	return path
}

// Returns the default data path defined on a system
// TODO: windows support
func defaultDataPath(target string) string {
	var path string
	// Use XDG standards
	xdgConfig, xdgCfgExists := os.LookupEnv("XDG_DATA_HOME")
	if xdgCfgExists {
		path = filepath.Join(xdgConfig, "cook", target)
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".local", "share", "cook", target)
	}

	return path
}

// Returns the default recipes directory path
func defaultRecipesPath() string {
	return defaultDataPath("recipes")
}

// Returns the default shopping directory path
func defaultShoppingPath() string {
	return defaultDataPath("shopping")
}

// Returns the default users.toml path
func defaultUsersPath() string {
	return defaultDataPath("users.toml")
}

// Tests for the existence of the data directories, if they do not exist then
// they are created.
//
// Returns nil on success
func EnsureDataDirInit() error {
	if !loaded {
		return errors.New("Attempted to init data dirs before loading configs")
	}
	var err error
	rDir := conf.Recipe.Dir
	sDir := conf.Shopping.Dir

	// Spawn recipe dir
	if !common.FileExists(rDir) {
		if err = os.MkdirAll(rDir, os.ModePerm); err != nil {
			return err
		}

		fmt.Printf("Created recipes directory at: %v\n", rDir)
	}

	// Spawn shopping list dir
	if !common.FileExists(sDir) {
		if err = os.MkdirAll(sDir, os.ModePerm); err != nil {
			return err
		}

		fmt.Printf("Created shopping list directory at: %v\n", sDir)
	}

	return nil
}

// Tests for the existence of the users file, if it does not exist then an
// empty one will be created.
//
// Returns nil on success
func EnsureUsersInit() error {
	if !loaded {
		return errors.New("Attempted to init users file before loading configs")
	}
	var err error
	users := conf.Users
	ext := filepath.Ext(users)
	if ext != ".toml" {
		errMsg := fmt.Sprintf(
			"Users must be `.toml` or a directory, found: %v\n", users)
		return errors.New(errMsg)
	} else if ext == "" {
		users = filepath.Join(users, "users.toml")
	}

	// Create the users' parent directory if needed
	if !common.FileExists(filepath.Dir(users)) {
		if err = os.MkdirAll(filepath.Dir(users), os.ModePerm); err != nil {
			return err
		}
	}

	if !common.FileExists(users) {
		var file *os.File
		if file, err = os.Create(users); err != nil {
			return err
		}
		file.Close()

		fmt.Printf("Created new users.toml at: %v\n", users)
	}

	return nil
}
