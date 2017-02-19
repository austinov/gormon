package ssh

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Interval  time.Duration     `mapstructure:"interval"`
	Keysdir   string            `mapstructure:"keys-dir"`
	FieldsOut string            `mapstructure:"fields-out"`
	FieldsWeb string            `mapstructure:"fields-web"`
	Hosts     map[string]Config `mapstructure:"hosts"`
}

type Config struct {
	User    string `mapstructure:"user"`
	Addr    string `mapstructure:"addr"`
	Keypath string `mapstructure:"keypath"`
	Command string `mapstructure:"cmd"`
}

var (
	cfg  AppConfig
	once sync.Once
)

func GetConfig() AppConfig {
	once.Do(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName("dev")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("fatal error config file: %s \n", err)
		}
		err = viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatal(err)
		}
	})
	return cfg
}
