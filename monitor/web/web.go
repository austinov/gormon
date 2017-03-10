package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	c "github.com/austinov/gormon/config"
	m "github.com/austinov/gormon/monitor"
)

type monitor struct {
	cfg    c.Config
	broker *Broker
}

func New(cfg c.Config) m.Monitor {
	broker := NewServer()
	m := &monitor{
		cfg:    cfg,
		broker: broker,
	}
	go func() {
		http.Handle("/", http.HandlerFunc(staticHandler("./monitor/web/static/", cfg)))
		http.HandleFunc("/stat/", broker.statHandler)
		log.Fatal("HTTP server error: ", http.ListenAndServe(cfg.Server, nil))
	}()
	return m
}

func (m *monitor) Close() error {
	m.broker.Stop()
	return nil
}

func staticHandler(filesPath string, cfg c.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if path == "" || path == "app.htm" {
			path = fmt.Sprintf("%s/%s", filesPath, "app.htm")
			t, err := template.ParseFiles(path)
			if err != nil {
				log.Fatal(err)
			}
			t.Execute(w, cfg.Hosts)
		} else {
			path = fmt.Sprintf("%s/%s", filesPath, path)
			http.ServeFile(w, r, path)
		}
	}
}

func (m *monitor) Process(host, output string) {
	stats := make(map[string]string)
	stats["host"] = host
	stats["tstamp"] = fmt.Sprintf("%d", time.Now().Unix())

	lines := strings.Split(output, "\r\n")
	for _, l := range lines {
		if l == "" {
			continue
		} else if l[:1] == "#" {
			if len(l) > 4 && l[:4] == "# PS" {
				cpu := strings.SplitN(l, "\n", 3)[1]
				stats["used_cpu_perc"] = strings.Replace(cpu, ",", ".", -1)
			} else {
				continue
			}
		} else {
			stat := strings.SplitN(l, ":", 2)
			if m.cfg.HasFieldOut(stat[0]) {
				stats[stat[0]] = stat[1]
			}
		}
	}
	bytes, err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err)
	}
	m.broker.Notifier <- bytes
}
