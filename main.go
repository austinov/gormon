package main

import (
	"sync"
	"syscall"

	c "github.com/austinov/gormon/config"
	"github.com/austinov/gormon/monitor/csv"
	"github.com/austinov/gormon/ssh"
	"github.com/austinov/gormon/utils"
)

func main() {
	cfg := c.GetConfig()

	monitor := csv.New(cfg)
	defer monitor.Close()

	var wg sync.WaitGroup

	collectors := make([]ssh.Collector, 0, len(cfg.Hosts))
	for _, hostConfig := range cfg.Hosts {
		client := ssh.NewClient(hostConfig)
		collector := ssh.NewCollector(cfg, client, monitor)
		collectors = append(collectors, collector)

		wg.Add(1)
		go func() {
			defer wg.Done()
			collector.Start()
		}()
	}

	// handle stop signals
	utils.SignalsHandle(func() {
		for _, collector := range collectors {
			collector.Stop()
		}
	}, syscall.SIGINT, syscall.SIGTERM)

	wg.Wait()
}
