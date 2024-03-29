package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env"`
	Db         Db     `yaml:"db"`
	HTTPServer `yaml:"http_server"`
}

type Db struct {
	DbPath   string `yaml:"dbpath"`
	DbName   string `yaml:"dbname"`
	CollName string `yaml:"collectionname"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
}

func MustLoad() *Config {
	confPath := "../config/config.yaml"

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", confPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(confPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	return &cfg
}
