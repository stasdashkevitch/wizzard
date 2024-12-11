package config

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"prod"`
	GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout"`
}

var instance *Config
var once sync.Once

func NewConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		path := fetchPath()
		if path == "" {
			panic("config is empty")
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic("config file does not exist: " + path)
		}

		err := cleanenv.ReadConfig(path, instance)
		if err != nil {
			panic("failed to read config: " + err.Error())
		}
	})

	return instance
}

func fetchPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
		fmt.Println(os.Getenv("CONFIG_PATH"))
	}

	return res
}
