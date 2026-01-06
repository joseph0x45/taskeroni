package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"golang.org/x/crypto/bcrypt"
)

func EnsureDirExists(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func AppDataDir(user *user.User) string {
	return path.Join(
		user.HomeDir,
		".local/share/taskeroni",
	)
}

func AppConfigDir(user *user.User) string {
	return path.Join(
		user.HomeDir,
		".config/taskeroni",
	)
}

func AppDatabasePath(user *user.User) string {
	return path.Join(
		AppDataDir(user),
		"taskeroni.db",
	)
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("Error while hashing password: %w", err)
	}
	return string(hashed), nil
}

func HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}
