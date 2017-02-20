package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/austinov/gormon/ssh"
	"github.com/austinov/gormon/types"
)

func main() {
	// TODO log into stdout/stderr
	log.SetOutput(os.Stdout)

	cfg := types.GetConfig()
	log.Printf("%#v\n", cfg)

	clients := make([]ssh.Client, 0, len(cfg.Hosts))
	for _, hostConfig := range cfg.Hosts {
		clients = append(clients, ssh.New(hostConfig))
	}

	ticker := time.Tick(2 * time.Second)

	m := NewMon(cfg)
	for {
		select {
		case <-ticker:
			for _, client := range clients {
				if err := client.Connect(); err != nil {
					log.Print(err)
					continue
				}
				if out, err := client.Run(); err != nil {
					log.Println("Command err: ", err)
				} else {
					m.Process(client.Host(), out)
				}
			}
		}
	}
}

type monitor struct {
	cfg types.Config
}

func NewMon(cfg types.Config) *monitor {
	m := &monitor{cfg}
	m.header()
	return m
}

var sep = "|"

func (m *monitor) header() {
	for i, f := range m.cfg.FieldsOut {
		if i == 0 {
			fmt.Printf("%s", f)
		} else {
			fmt.Printf("%s%s", sep, f)
		}
	}
	fmt.Println()
}

func (m *monitor) Process(host, output string) {
	stats := make(map[string]string)
	stats["host"] = host

	lines := strings.Split(output, "\n")
	for _, l := range lines {
		if l == "" || l[:1] == "#" {
			continue
		} else {
			stat := strings.Split(l, ":")
			if m.cfg.HasFieldOut(stat[0]) {
				stats[stat[0]] = stat[1]
			}
		}
	}
	m.printStat(stats)
}

func (m *monitor) printStat(stats map[string]string) {
	out := ""
	clear := func(s string) string {
		if strings.HasSuffix(s, "\r") {
			return s[:len(s)-1]
		}
		return s
	}
	for _, f := range m.cfg.FieldsOut {
		data := clear(stats[f])
		if out == "" {
			out = data
		} else {
			out += sep + data
		}
	}
	fmt.Println(out)
}
