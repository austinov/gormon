package csv

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	c "github.com/austinov/gormon/config"
	m "github.com/austinov/gormon/monitor"
	"github.com/austinov/gormon/utils"
)

type monitor struct {
	cfg  c.Config
	prev map[string]map[string]string
}

func New(cfg c.Config) m.Monitor {
	m := &monitor{
		cfg:  cfg,
		prev: make(map[string]map[string]string),
	}
	m.header()
	return m
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
			stat := strings.SplitN(l, ":", 2)
			if m.cfg.HasFieldOut(stat[0]) {
				stats[stat[0]] = stat[1]
			}
		}
	}
	m.printStat(stats)
}

func (m *monitor) header() {
	for i, f := range m.cfg.FieldsOut {
		if i == 0 {
			fmt.Printf("%s", f)
		} else {
			fmt.Printf("%s%s", m.cfg.FieldsSeparator, f)
		}
	}
	fmt.Println()
}

func (m *monitor) printStat(stats map[string]string) {
	out := ""
	needPrint := false
	prev := m.prev[stats["host"]]
	for _, f := range m.cfg.FieldsOut {
		data := stats[f]
		if f != "host" && f != "tstamp" {
			np := false
			if data, np = m.format(f, data, prev); np {
				needPrint = true
			}
		}
		if out == "" {
			out = data
		} else {
			out += m.cfg.FieldsSeparator + data
		}
	}
	if out != "" && needPrint {
		fmt.Println(out)
	}
	m.prev[stats["host"]] = stats
}

func (m *monitor) format(field, value string, prev map[string]string) (string, bool) {
	parseInts := func() (data, delta int64, factor float32) {
		var err error
		if data, err = strconv.ParseInt(value, 10, 64); err == nil {
			prevData := prev[field]
			if pv, err := strconv.ParseInt(prevData, 10, 64); err == nil && pv != 0 {
				delta = data - pv
				factor = float32(delta*100) / float32(pv)
				if factor < 0 {
					factor = -factor
				}
			}
		}
		return data, delta, factor
	}
	parseFloats := func() (data, delta float64, factor float32) {
		var err error
		if data, err = strconv.ParseFloat(value, 64); err == nil {
			prevData := prev[field]
			if pv, err := strconv.ParseFloat(prevData, 64); err == nil {
				delta := data - pv
				factor := float32(delta*100) / float32(pv)
				if factor < 0 {
					factor = -factor
				}
			}
		}
		return data, delta, factor
	}
	needPrint := false
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
		_, delta, factor := parseInts()
		if m.cfg.OutFormat == "csv++" {
			if factor > m.cfg.ChangeFactor {
				value += fmt.Sprintf(" (%d)", delta)
				needPrint = true
			} else if prev[field] == "" && value != "" {
				needPrint = true
			}
		} else {
			if m.cfg.OutFormat == "csv+" && factor > m.cfg.ChangeFactor {
				value += fmt.Sprintf(" (%d)", delta)
			}
			needPrint = true
		}
	case "used_memory",
		"used_memory_rss",
		"used_memory_peak",
		"used_memory_lua":
		data, delta, factor := parseInts()
		sign := ""
		if delta < 0 {
			sign = "-"
			delta = -delta
		}
		if m.cfg.OutFormat == "csv++" {
			if factor > m.cfg.ChangeFactor {
				if m.cfg.IsHumanReadable() {
					value = fmt.Sprintf("%s (%s%s)", utils.HumanBytes(uint64(data)), sign, utils.HumanBytes(uint64(delta)))
				} else {
					value = fmt.Sprintf("%d (%d)", data, delta)
				}
				needPrint = true
			} else if prev[field] == "" && value != "" {
				if m.cfg.IsHumanReadable() {
					value = fmt.Sprintf("%s", utils.HumanBytes(uint64(data)))
				} else {
					value = fmt.Sprintf("%d", data)
				}
				needPrint = true
			}
		} else {
			if m.cfg.OutFormat == "csv+" && factor > m.cfg.ChangeFactor {
				if m.cfg.IsHumanReadable() {
					value = fmt.Sprintf("%s (%s%s)", utils.HumanBytes(uint64(data)), sign, utils.HumanBytes(uint64(delta)))
				} else {
					value = fmt.Sprintf("%d (%d)", data, delta)
				}
			} else if m.cfg.IsHumanReadable() {
				value = fmt.Sprintf("%s", utils.HumanBytes(uint64(data)))
			}
			needPrint = true
		}
	case "mem_fragmentation_ratio",
		"instantaneous_input_kbps",
		"instantaneous_output_kbps",
		"used_cpu_sys",
		"used_cpu_user",
		"used_cpu_sys_children",
		"used_cpu_user_children":
		_, delta, factor := parseFloats()
		if m.cfg.OutFormat == "csv++" {
			if factor > m.cfg.ChangeFactor {
				value += fmt.Sprintf(" (%f)", delta)
				needPrint = true
			} else if prev[field] == "" && value != "" {
				needPrint = true
			}
		} else {
			if m.cfg.OutFormat == "csv+" && factor > m.cfg.ChangeFactor {
				value += fmt.Sprintf(" (%f)", delta)
			}
			needPrint = true
		}
	default:
		if m.cfg.OutFormat == "csv++" {
			needPrint = value != prev[field]
		} else {
			needPrint = true
		}
	}
	return value, needPrint
}
