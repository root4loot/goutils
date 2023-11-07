package sysutil

import (
	"os/user"
)

// IsRoot checks if the current user is root.
func IsRoot() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, err
	}
	return currentUser.Username == "root", nil
}
