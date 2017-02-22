package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	c "github.com/austinov/gormon/config"
	"github.com/austinov/gormon/monitor/csv"
	"github.com/austinov/gormon/ssh"
)

func main() {
	// TODO log into stdout/stderr
	log.SetOutput(os.Stdout)

	cfg := c.GetConfig()
	log.Printf("%#v\n", cfg)

	var wg sync.WaitGroup
	mon := csv.New(cfg)

	collectors := make([]ssh.Collector, 0, len(cfg.Hosts))
	for _, hostConfig := range cfg.Hosts {
		client := ssh.NewClient(hostConfig)
		collector := ssh.NewCollector(cfg, client, mon)
		collectors = append(collectors, collector)

		wg.Add(1)
		go func() {
			defer wg.Done()
			collector.Start()
		}()
	}

	// handle stop signals
	interrupt := make(chan os.Signal, 1)
	defer close(interrupt)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-interrupt
		for _, collector := range collectors {
			collector.Stop()
		}
		signal.Stop(interrupt)
	}()

	wg.Wait()
}
