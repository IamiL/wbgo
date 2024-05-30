package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env            string              `yaml:"env" env-default:"local"`
	NatsStreaming  NatsStreamingConfig `yaml:"nats_streaming"`
	HTTPServer     `yaml:"http_server"`
	PostgresConfig `yaml:"postgresql"`
}

type NatsStreamingConfig struct {
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
}

type HTTPServer struct {
	Port    string        `yaml:"port" env-default:":8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type PostgresConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DBName string `yaml:"database_name"`
	User   string `yaml:"username"`
	Pass   string `yaml:"password"`
}

func MustLoad() *Config {
	configPath := "config/local.yaml"

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
