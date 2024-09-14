package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var Cfg Config

type (
	// Config .
	Config struct {
		Rest    `yaml:"rest"`
		Storage `yaml:"storage"`
		Common  `yaml:"common"`
	}

	// Common .
	Common struct {
		ConsulAddress   string `env-required:"true" yaml:"consul_address"    env:"CONSUL_ADDRESS"`
		DB              `yaml:"db"`
		HealthCheckPort string `yaml:"health_check_port" env:"HEALTH_CHECK_PORT"`
	}

	DB struct {
		Host     string `yaml:"postgres_host"    env:"POSTGRES_HOST"`
		Port     string `yaml:"postgres_port" env:"POSTGRES_PORT"`
		User     string `yaml:"postgres_user"    env:"POSTGRES_USER"`
		Password string `yaml:"postgres_password" env:"POSTGRES_PASSWORD"`
		DbName   string `yaml:"postgres_db" env:"POSTGRES_DB"`
	}

	// Rest .
	Rest struct {
		Name                        string        `env-required:"true" yaml:"name"    env:"REST_NAME"`
		Port                        string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		UpdateStorageServersTimeout time.Duration `env-required:"true" yaml:"update_storage_servers_timeout" env:"UPDATE_STORAGE_SERVERS_TIMEOUT"`
	}

	// Storage .
	Storage struct {
		Name string `yaml:"name"    env:"STORAGE_NAME"`
		Port string `yaml:"port" env:"STORAGE_PORT"`
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
