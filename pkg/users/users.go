// -------------------------------------------
// This is a custom rolled user "database" into a TOML format. AS SUCH,
// it is not performant enough to handle large numbers of users, and is never
// intended to be.
//
// (initial)Reads and Writes require full memory loads of the users file and
// massive files WILL cause issues.
// -------------------------------------------

package users

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"

	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/bcrypt"
)

var users map[string]UserData
var usersLoaded = false

type UserData struct {
	PassHash string `toml:"pass"`
}

func ensureUsersLoaded() {
	if !usersLoaded {
		loadUsersFile()
	}
}

// Writes `users` to the toml file specified by config
func saveUsersFile() error {
	ensureUsersLoaded()
	var err error

	// Create temp file
	usersPath := config.GetConfig().Users
	var file *os.File
	if file, err = os.Create(usersPath + ".tmp"); err != nil {
		return err
	}

	// Write to temp
	if err = toml.NewEncoder(file).Encode(users); err != nil {
		return err
	}

	// Delete original users
	if err = os.Remove(usersPath); err != nil {
		return err
	}

	// Move temp into users
	if err = os.Rename(usersPath+".tmp", usersPath); err != nil {
		return err
	}

	return nil
}

// Loads the user file defined in config into `users`, setting `usersLoaded` to true
func loadUsersFile() {
	if !config.IsLoaded() {
		panic("Attempted to read users before loading config")
	}
	var err error

	// Parse the users file as `.toml` into `users`
	usersFile := config.GetConfig().Users
	if _, err = toml.DecodeFile(usersFile, &users); err != nil {
		errMsg := fmt.Sprintf(
			`Failed to open users file: %v\n\t%v`,
			usersFile, err)
		panic(fmt.Sprintf(errMsg))
	}

	usersLoaded = true
}

// Tests if a username contains illegal characters
func legalUsername(username string) bool {
	ok, _ := regexp.MatchString(`[a-zA-Z_\-0-9]+`, username)
	return ok
}

// Tests for existence of specified user in the user file
func UserExists(username string) bool {
	ensureUsersLoaded()
	_, ok := users[username]
	return ok
}

// Tests a raw password and username against the credentials in the `users.toml`
// file.
func ValidateUser(username string, password string) bool {
	ensureUsersLoaded()
	userData, ok := users[username]
	if ok {
		hashBytes, _ := hex.DecodeString(userData.PassHash)
		return bcrypt.CompareHashAndPassword(hashBytes, []byte(password)) == nil
	} else {
		return false
	}
}

// Appends a user to the `users.toml`, The username must not be taken.
// The password will be hashed and a salt will be generated and stored alongside
// it.
//
// Returns nil on success
func AddUser(username string, password string) error {
	ensureUsersLoaded()
	var err error

	// Sanity checks
	if UserExists(username) {
		return errors.New("Username already exists")
	}

	if !legalUsername(username) {
		return errors.New("Username contains illegal characters")
	}

	// Generate data to store
	var hashedPass string
	if hashedPass, err = HashPassword(password); err != nil {
		return err
	}

	// Push the record
	users[username] = UserData{
		PassHash: hashedPass,
	}
	if err = saveUsersFile(); err != nil {
		return err
	}

	return nil
}

func RemoveUser(username string) error {
	ensureUsersLoaded()

	_, ok := users[username]
	if !ok {
		return errors.New("User does not exist: " + username)
	}

	delete(users, username)
	return saveUsersFile()
}

// Edits the password record of a specific user. The user must exist.
//
// Returns nil on success
func ChangePassword(username string, password string) error {
	var err error

	// Grab the existing record
	user, ok := users[username]
	if !ok {
		return errors.New(
			"Attempted to change password of non-existent user: " + username)
	}

	// Hash the password
	var passHash string
	if passHash, err = HashPassword(password); err != nil {
		return err
	}

	// Assign the new password
	user.PassHash = passHash
	users[username] = user
	return saveUsersFile()
}

// Hashes a password to be stored in the users file, a salt is applied
// based on the salt specified in the config.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}
