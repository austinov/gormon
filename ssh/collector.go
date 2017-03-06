package ssh

import (
	"math/rand"
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
	monitors  []m.Monitor
	done      chan struct{}
	connected bool
}

func NewCollector(cfg c.Config, client Client, mon ...m.Monitor) Collector {
	return &collector{
		cfg:      cfg,
		client:   client,
		monitors: mon,
		done:     make(chan struct{}, 1),
	}
}

func (c *collector) Start() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	normal := time.Duration(rnd.NormFloat64() * float64(time.Millisecond))
	ticker := time.Tick(c.cfg.Interval + normal)
	for {
		select {
		case <-ticker:
			if !c.connected {
				if err := c.client.Connect(); err != nil {
					c.send(c.client.Host(), "error:"+err.Error())
					continue
				}
				c.connected = true
			}
			if out, err := c.client.Run(); err != nil {
				c.send(c.client.Host(), "error:"+err.Error())
			} else {
				c.send(c.client.Host(), out)
			}
		case _, ok := <-c.done:
			if !ok {
				c.client.Disconnect()
				return
			}
		}
	}
}

func (c *collector) send(host, output string) {
	for _, mon := range c.monitors {
		mon.Process(host, output)
	}
}

func (c *collector) Stop() {
	close(c.done)
}
