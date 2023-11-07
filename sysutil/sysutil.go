package sysutil

import (
	"log"
	"os/user"
)

// IsRoot checks if the current user is root.
func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}
