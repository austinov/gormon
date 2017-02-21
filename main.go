package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/austinov/gormon/ssh"
	"github.com/austinov/gormon/types"
	"github.com/dustin/go-humanize"
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
	cfg  types.Config
	prev map[string]map[string]string
}

func NewMon(cfg types.Config) *monitor {
	m := &monitor{
		cfg:  cfg,
		prev: make(map[string]map[string]string),
	}
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
	stats["tstamp"] = time.Now().Format("2006-01-02 15:04:05.999")

	lines := strings.Split(output, "\r\n")
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
	prev := m.prev[stats["host"]]
	for _, f := range m.cfg.FieldsOut {
		data := stats[f]
		if f != "host" && f != "tstamp" && prev != nil {
			data = m.format(f, data, prev)
		}
		if out == "" {
			out = data
		} else {
			out += sep + data
		}
	}
	if out != "" {
		fmt.Println(out)
	}
	m.prev[stats["host"]] = stats
}

func (m *monitor) format(field, value string, prev map[string]string) string {
	parseInts := func() (data, delta int64, factor int) {
		var err error
		if data, err = strconv.ParseInt(value, 10, 64); err == nil {
			prevData := prev[field]
			if pv, err := strconv.ParseInt(prevData, 10, 64); err == nil && pv != 0 {
				delta = data - pv
				factor = int((delta * 100) / pv)
				if factor < 0 {
					factor = -factor
				}
			}
		}
		return data, delta, factor
	}
	parseFloats := func() (data, delta float64, factor int) {
		var err error
		if data, err = strconv.ParseFloat(value, 64); err == nil {
			prevData := prev[field]
			if pv, err := strconv.ParseFloat(prevData, 64); err == nil {
				delta := data - pv
				factor := int((delta * 100) / pv)
				if factor < 0 {
					factor = -factor
				}
			}
		}
		return data, delta, factor
	}
	switch field {
	case "connected_clients",
		"client_longest_output_list",
		"client_biggest_input_buf",
		"blocked_clients",
		"loading",
		"rdb_changes_since_last_save",
		"rdb_bgsave_in_progress",
		"rdb_last_save_time",
		"rdb_last_bgsave_status",
		"rdb_last_bgsave_time_sec",
		"rdb_current_bgsave_time_sec",
		"aof_enabled",
		"aof_rewrite_in_progress",
		"aof_rewrite_scheduled",
		"aof_last_rewrite_time_sec",
		"aof_current_rewrite_time_sec",
		"total_connections_received",
		"total_commands_processed",
		"instantaneous_ops_per_sec",
		"total_net_input_bytes",
		"total_net_output_bytes",
		"rejected_connections",
		"sync_full",
		"sync_partial_ok",
		"sync_partial_err",
		"expired_keys",
		"evicted_keys",
		"keyspace_hits",
		"keyspace_misses",
		"pubsub_channels",
		"pubsub_patterns",
		"latest_fork_usec",
		"migrate_cached_sockets",
		"connected_slaves",
		"master_repl_offset",
		"repl_backlog_active",
		"repl_backlog_size",
		"repl_backlog_first_byte_offset",
		"repl_backlog_histlen",
		"cluster_enabled":
		if _, delta, factor := parseInts(); factor > m.cfg.ChangeFactor {
			value += fmt.Sprintf(" (%d)", delta)
		}
	case "used_memory",
		"used_memory_rss",
		"used_memory_peak",
		"used_memory_lua":
		if data, delta, factor := parseInts(); factor > m.cfg.ChangeFactor {
			value = fmt.Sprintf("%s (%s)", humanize.Bytes(uint64(data)), humanize.Bytes(uint64(delta)))
		} else {
			value = fmt.Sprintf("%s", humanize.Bytes(uint64(data)))
		}
	case "mem_fragmentation_ratio",
		"instantaneous_input_kbps",
		"instantaneous_output_kbps",
		"used_cpu_sys",
		"used_cpu_user",
		"used_cpu_sys_children",
		"used_cpu_user_children":
		if _, delta, factor := parseFloats(); factor > m.cfg.ChangeFactor {
			value += fmt.Sprintf(" (%f)", delta)
		}
	}
	return value
}
