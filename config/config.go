package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	fu "github.com/austinov/go-recipes/flagutils"
	"github.com/austinov/gormon/utils"
	"github.com/spf13/viper"
)

const (
	defPath = "."
	defName = "dev"
)

var (
	cfg     Config
	once    sync.Once
	cfgPath string
	cfgName string
	server  fu.StringFlag
	debug   bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Please, make sure you have the configuration file and use flags to setup path to it.\n")
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.StringVar(&cfgPath, "cfg-dir", defPath, "dir with app's config")
	flag.StringVar(&cfgName, "cfg-name", defName, "app's config base file name")
	flag.Var(&server, "server", "launch as server on address")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
}

type Config struct {
	OutFormat       string                `mapstructure:"out-format"`
	FieldsSeparator string                `mapstructure:"fields-separator"`
	HumanReadable   bool                  `mapstructure:"human-readable"`
	AdjustWidth     bool                  `mapstructure:"adjust-width"`
	Interval        time.Duration         `mapstructure:"interval"`
	Keydir          string                `mapstructure:"ssh-key-dir"`
	FieldsOut       []string              `mapstructure:"fields-out"`
	ChangeFactor    float32               `mapstructure:"change-factor"`
	Hosts           map[string]HostConfig `mapstructure:"hosts"`
	Server          string                `mapstructure:"server"`
	Debug           bool
	fieldsOutMap    map[string]struct{}
}

type HostConfig struct {
	User    string `mapstructure:"user"`
	Addr    string `mapstructure:"addr"`
	Keypath string `mapstructure:"keypath"`
	Command string `mapstructure:"cmd"`
}

func (c *Config) init() {
	formats := map[string]struct{}{
		"csv":   struct{}{},
		"csv+":  struct{}{},
		"csv++": struct{}{},
	}
	if _, ok := formats[c.OutFormat]; !ok {
		c.OutFormat = "csv+"
	}
	if c.FieldsSeparator == "" {
		c.FieldsSeparator = "|"
	}
	if c.Interval == time.Duration(0) {
		c.Interval = 2 * time.Second
	}
	if c.Keydir == "" {
		c.Keydir = "~/.ssh"
	}
	for name, host := range c.Hosts {
		host.ExpandKeypath(c.Keydir)
		c.Hosts[name] = host
	}
	if c.FieldsOut == nil {
		// default stat fields
		c.FieldsOut = []string{
			"host",
			"tstamp",
			"used_memory",
			"used_memory_rss",
			"connected_clients",
			"blocked_clients",
			"rejected_connections",
			"keyspace_hits",
			"keyspace_misses",
			"used_cpu_sys",
			"used_cpu_user",
			"aof_last_write_status",
			"error",
		}
	}
	c.fieldsOutMap = make(map[string]struct{})
	var host, tstamp bool
	for _, f := range c.FieldsOut {
		if f == "host" {
			host = true
		} else if f == "tstamp" {
			tstamp = true
		}
		c.fieldsOutMap[f] = struct{}{}
	}
	if !tstamp {
		c.FieldsOut = append([]string{"tstamp"}, c.FieldsOut...)
		c.fieldsOutMap["tstamp"] = struct{}{}
	}
	if !host {
		c.FieldsOut = append([]string{"host"}, c.FieldsOut...)
		c.fieldsOutMap["host"] = struct{}{}
	}
	if server.Exist {
		c.Server = server.Value
	}
	c.Debug = debug
}

func (c *Config) HasFieldOut(name string) bool {
	_, ok := c.fieldsOutMap[name]
	return ok
}

func (h *HostConfig) ExpandKeypath(keydir string) {
	if filepath.IsAbs(h.Keypath) {
		return
	}
	keydir = utils.ExpandPath(keydir)
	keypath := filepath.Join(keydir, h.Keypath)
	if _, err := os.Stat(keypath); err != nil {
		log.Fatalf("expand key path error: %s", err.(*os.PathError).Error())
	} else {
		h.Keypath = keypath
	}
}

func GetConfig() Config {
	once.Do(func() {
		if cfgPath != "" {
			viper.AddConfigPath(cfgPath)
		}
		viper.AddConfigPath(defPath)
		viper.AddConfigPath(".")
		if cfgName == "" {
			cfgName = defName
		}
		viper.SetConfigName(cfgName)

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("fatal error config file: %s \n", err)
		}
		err = viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatal(err)
		}
		cfg.init()
		if cfg.Debug {
			log.Printf("config: %#v", cfg)
		}
	})
	return cfg
}
