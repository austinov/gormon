package main

import (
	"log"
	"os"
	"strings"

	"github.com/austinov/gormon/ssh"
)

func main() {
	// TODO log into stdout/stderr
	log.SetOutput(os.Stdout)
	/*
		var err error
		currentUser, err = user.Current()
		if err != nil {
			log.Print(err)
			return
		}
		keypath := ""
		idrsap := filepath.Join(currentUser.HomeDir, ".ssh", keyfile)
		if _, err := os.Stat(idrsap); err == nil {
			keypath = idrsap
		}
	*/

	/*
		configs := []ssh.Config{
			ssh.Config{
				User:    "root",
				Addr:    "109.234.37.213:22",
				Keypath: "/home/dev/.ssh/id_rsa_vdsina",
				Command: "redis-cli -s /var/lib/letsrock/sockets/redis/letsrock-redis.socket INFO",
			},
			ssh.Config{
				User:    "vdev",
				Addr:    "192.168.1.104:22",
				Keypath: "/home/dev/.ssh/id_rsa_vbox_ub16",
				Command: "redis-cli INFO",
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
					out, err := client.Run()
					log.Println("*****************")
					//log.Println("Redis: ", cfg.Addr)
					//log.Println("*****************")
					log.Println("Command err: ", err)
					log.Printf("Command out:\n%s\n", out)
				}
			}
		}
	*/
	stats := make([]info, 0)
	stats = append(stats, info{"host", "192.168.1.104"})

	lines := strings.Split(output, "\n")
	for _, l := range lines {
		if l == "" || l[:1] == "#" {
			continue
		} else {
			stat := strings.Split(l, ":")
			stats = append(stats, info{stat[0], stat[1]})
		}
	}
	for _, s := range stats {
		log.Println(s.name, s.value)
	}
	log.Printf("%#v\n", ssh.GetConfig())
}

type info struct {
	name  string
	value string
}

var output = `# Server
redis_version:3.0.6
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:83b1283ee21c4d12
redis_mode:standalone
os:Linux 3.10.0-327.22.2.el7.x86_64 x86_64
arch_bits:64
multiplexing_api:epoll
gcc_version:4.8.5
process_id:6900
run_id:83438629d2113a5cd0c84865865ebfebfd043b9c
tcp_port:0
uptime_in_seconds:14487344
uptime_in_days:167
hz:10
lru_clock:11397207
config_file:/etc/redis.conf

# Clients
connected_clients:8
client_longest_output_list:0
client_biggest_input_buf:0
blocked_clients:0

# Memory
used_memory:11608352
used_memory_human:11.07M
used_memory_rss:17747968
used_memory_peak:206884888
used_memory_peak_human:197.30M
used_memory_lua:1146880
mem_fragmentation_ratio:1.53
mem_allocator:jemalloc-3.6.0

# Persistence
loading:0
rdb_changes_since_last_save:8
rdb_bgsave_in_progress:0
rdb_last_save_time:1487792134
rdb_last_bgsave_status:ok
rdb_last_bgsave_time_sec:0
rdb_current_bgsave_time_sec:-1
aof_enabled:0
aof_rewrite_in_progress:0
aof_rewrite_scheduled:0
aof_last_rewrite_time_sec:-1
aof_current_rewrite_time_sec:-1
aof_last_bgrewrite_status:ok
aof_last_write_status:ok

# Stats
total_connections_received:178729
total_commands_processed:336844261
instantaneous_ops_per_sec:0
total_net_input_bytes:4823003739
total_net_output_bytes:10535763829
instantaneous_input_kbps:0.00
instantaneous_output_kbps:0.00
rejected_connections:0
sync_full:0
sync_partial_ok:0
sync_partial_err:0
expired_keys:24
evicted_keys:0
keyspace_hits:85465985
keyspace_misses:378048
pubsub_channels:0
pubsub_patterns:0
latest_fork_usec:890
migrate_cached_sockets:0

# Replication
role:master
connected_slaves:0
master_repl_offset:0
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0

# CPU
used_cpu_sys:7863.76
used_cpu_user:8649.66
used_cpu_sys_children:378.82
used_cpu_user_children:3978.28

# Cluster
cluster_enabled:0

# Keyspace
db0:keys=28,expires=2,avg_ttl=547995183`
