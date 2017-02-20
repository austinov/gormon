package utils

import (
	"log"
	"os/user"
	"strings"
)

func ExpandPath(path string) string {
	if len(path) < 2 || path[:2] != "~/" {
		return path
	}
	currUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(path, "~", currUser.HomeDir, 1)
}
