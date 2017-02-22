package main

import (
	"log"
	"os"
	"time"

	c "github.com/austinov/gormon/config"
	"github.com/austinov/gormon/monitor/csv"
	"github.com/austinov/gormon/ssh"
)

func main() {
	// TODO log into stdout/stderr
	log.SetOutput(os.Stdout)

	cfg := c.GetConfig()
	log.Printf("%#v\n", cfg)

	clients := make([]ssh.Client, 0, len(cfg.Hosts))
	for _, hostConfig := range cfg.Hosts {
		clients = append(clients, ssh.New(hostConfig))
	}

	ticker := time.Tick(2 * time.Second)

	m := csv.New(cfg)
	for {
		select {
		case <-ticker:
			for _, client := range clients {
				if err := client.Connect(); err != nil {
					m.Process(client.Host(), "error:"+err.Error())
					continue
				}
				if out, err := client.Run(); err != nil {
					m.Process(client.Host(), "error:"+err.Error())
				} else {
					m.Process(client.Host(), out)
				}
			}
		}
	}
}
