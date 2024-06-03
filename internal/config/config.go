package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "configs/config.yaml"
)

type Config struct {
	DB   PosgreConfig `yaml:"db"`
	Http HttpConfig   `yaml:"http"`
}

type PosgreConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type HttpConfig struct {
	Port int `yaml:"port"`
}

var (
	config Config
)

func Read() *Config {
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic(err)
	}
	return &config
}
