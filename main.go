package main

import (
	"sync"
	"syscall"

	c "github.com/austinov/gormon/config"
	m "github.com/austinov/gormon/monitor"
	"github.com/austinov/gormon/monitor/csv"
	"github.com/austinov/gormon/monitor/web"
	"github.com/austinov/gormon/ssh"
	"github.com/austinov/gormon/utils"
)

func main() {
	cfg := c.GetConfig()

	csvMonitor := csv.New(cfg)
	defer csvMonitor.Close()

	monitors := make([]m.Monitor, 0)
	monitors = append(monitors, csvMonitor)

	if len(cfg.Server) > 0 {
		webMonitor := web.New(cfg)
		defer webMonitor.Close()
		monitors = append(monitors, webMonitor)
	}

	var wg sync.WaitGroup

	collectors := make([]ssh.Collector, 0, len(cfg.Hosts))
	for _, hostConfig := range cfg.Hosts {
		client := ssh.NewClient(hostConfig)
		collector := ssh.NewCollector(cfg, client, monitors...)
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
