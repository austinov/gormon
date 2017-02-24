package utils

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"os/user"
	"strings"
)

// ExpandPath expands file path if it starts with a tilde (~).
// E.g.: ~/.ssh will be expanded /home/current-user/.ssh
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

// HumanBytes returns human readable bytes count.
func HumanBytes(bytes uint64) string {
	var unit uint64 = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	exp := int(math.Log(float64(bytes)) / math.Log(float64(unit)))
	pre := "KMGTPE"[exp-1 : exp]
	return fmt.Sprintf("%.2f %sB", float64(bytes)/math.Pow(float64(unit), float64(exp)), pre)
}

// SignalsHandle executes the handler when any signal is received.
func SignalsHandle(handler func(), sig ...os.Signal) {
	interrupter := make(chan os.Signal, 1)
	signal.Notify(interrupter, sig...)
	go func() {
		defer close(interrupter)
		<-interrupter
		handler()
		signal.Stop(interrupter)
	}()
}
