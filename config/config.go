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
		Name                         string        `env-required:"true" yaml:"name"    env:"REST_NAME"`
		Port                         string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		UpdateStorageServersInterval time.Duration `yaml:"update_storage_servers_interval" env:"UPDATE_STORAGE_SERVERS_INTERVAL" env-default:"10s"`
		StorageServersCount          int32         `yaml:"storage_servers_count" env:"STORAGE_SERVERS_COUNT" env-default:"6"`
		UploadBytesWaitTime          time.Duration `yaml:"upload_bytes_wait_time" env:"UPLOAD_BYTES_WAIT_TIME" env-default:"15s"`
	}

	// Storage .
	Storage struct {
		Name string `yaml:"name"    env:"STORAGE_NAME"`
		Port string `yaml:"port" env:"STORAGE_PORT"`
	}
)

func Init() {
	err := cleanenv.ReadConfig("./config/config.yml", &Cfg)
	if err != nil {
		panic("can't read config file: " + err.Error())
	}

	err = cleanenv.ReadEnv(&Cfg)
	if err != nil {
		panic("can't read config from env: " + err.Error())
	}
}
