package ssh

import (
	"time"

	c "github.com/austinov/gormon/config"
	m "github.com/austinov/gormon/monitor"
)

type Collector interface {
	Start()
	Stop()
}

type collector struct {
	cfg       c.Config
	client    Client
	mon       m.Monitor
	done      chan struct{}
	connected bool
}

func NewCollector(cfg c.Config, client Client, mon m.Monitor) Collector {
	return &collector{
		cfg:    cfg,
		client: client,
		mon:    mon,
		done:   make(chan struct{}, 1),
	}
}

func (c *collector) Start() {
	ticker := time.Tick(c.cfg.Interval)
	for {
		select {
		case <-ticker:
			if !c.connected {
				if err := c.client.Connect(); err != nil {
					c.mon.Process(c.client.Host(), "error:"+err.Error())
					continue
				}
				c.connected = true
			}
			if out, err := c.client.Run(); err != nil {
				c.mon.Process(c.client.Host(), "error:"+err.Error())
			} else {
				c.mon.Process(c.client.Host(), out)
			}
		case <-c.done:
			c.client.Disconnect()
			return
		}
	}
}

func (c *collector) Stop() {
	c.done <- struct{}{}
}
