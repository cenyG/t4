package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

var Cfg Config

type (
	// Config .
	Config struct {
		Rest    `yaml:"rest"`
		Storage `yaml:"storage"`
	}

	// Rest .
	Rest struct {
		Name                        string        `env-required:"true" yaml:"name"    env:"REST_NAME"`
		Port                        string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		UpdateStorageServersTimeout time.Duration `env-required:"true" yaml:"update_storage_servers_timeout" env:"UPDATE_STORAGE_SERVERS_TIMEOUT"`
	}

	// Storage .
	Storage struct {
		Name string `env-required:"true" yaml:"name"    env:"STORAGE_NAME"`
		Port string `env-required:"true" yaml:"port" env:"STORAGE_PORT"`
	}
)

func init() {
	err := cleanenv.ReadConfig("./config/config.yml", &Cfg)
	if err != nil {
		panic("can't read config file: " + err.Error())
	}

	err = cleanenv.ReadEnv(&Cfg)
	if err != nil {
		panic("can't read config from env: " + err.Error())
	}
}
