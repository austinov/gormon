package main

import (
	"log"
	"os"
	"time"

	"github.com/austinov/gormon/ssh"
)

func main() {
	// TODO log into stdout/stderr
	log.SetOutput(os.Stdout)
	configs := []ssh.Config{
		ssh.Config{
			User:    "user1",
			Addr:    "192.168.1.13:22",
			Keypath: "id_rsa_1",
			Command: "redis-cli INFO",
		},
		ssh.Config{
			User:    "user2",
			Addr:    "192.168.1.14:22",
			Keypath: "id_rsa_2",
			Command: "redis-cli -s /var/lib/backend/sockets/redis/backend-redis.socket INFO",
		},
	}

	clients := make([]ssh.Client, len(configs))

	for i, cfg := range configs {
		clients[i] = ssh.New(cfg)
	}

	ticker := time.Tick(2 * time.Second)

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
					log.Println("***********************************")
					log.Printf("Command out:\n%s\n", out)
					log.Println("***********************************")
				}
			}
		}
	}
}
