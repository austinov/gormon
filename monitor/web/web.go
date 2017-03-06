package web

import (
	"encoding/json"
	"fmt"
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
		http.HandleFunc("/", m.indexHandler)
		http.HandleFunc("/stat/", broker.statHandler)
		log.Fatal("HTTP server error: ", http.ListenAndServe(cfg.Server, nil))
	}()
	return m
}

var index = `
<!DOCTYPE html>
<html>

<head>
  <meta charset="utf-8">
  <script>
    var eventSource;

    function start() {

      if (!window.EventSource) {
        alert('This browser doesn\'t support EventSource.');
        return;
      }

      eventSource = new EventSource('/stat/');

      eventSource.onerror = function(e) {
        if (this.readyState == EventSource.CONNECTING) {
          log("Reconnecting...");
        } else {
          log("Connection error: " + this.readyState);
        }
      };

      eventSource.onmessage = function(e) {
        console.log(e);
        log(e.data);
      };
    }

    function stop() {
      eventSource.close();
    }

    function log(msg) {
      logElem.innerHTML += msg + "<br>";
    }
  </script>
</head>

<body onload="start();">

  <button onclick="start()">Start</button>
  <button onclick="stop()">Stop</button>

  <div id="logElem"></div>

</body>

</html>`

func (m *monitor) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, index)
}

func (m *monitor) Close() error {
	m.broker.Stop()
	return nil
}

func (m *monitor) Process(host, output string) {
	stats := make(map[string]string)
	stats["host"] = host
	stats["tstamp"] = fmt.Sprintf("%d", time.Now().Unix())

	lines := strings.Split(output, "\r\n")
	for _, l := range lines {
		if l == "" || l[:1] == "#" {
			continue
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
