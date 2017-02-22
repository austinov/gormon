package utils

import (
	"fmt"
	"log"
	"math"
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

func HumanBytes(bytes uint64) string {
	var unit uint64 = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	exp := int(math.Log(float64(bytes)) / math.Log(float64(unit)))
	pre := "KMGTPE"[exp-1 : exp]
	return fmt.Sprintf("%.2f %sB", float64(bytes)/math.Pow(float64(unit), float64(exp)), pre)
}
