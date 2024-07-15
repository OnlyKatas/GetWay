package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Http HTTPConfig `yaml:"http_server"`
	GRPC GRPC       `yaml:"grpc"`
}

type HTTPConfig struct {
	Port string `yaml:"port"`
}

type GRPC struct {
	AuthAddress string `yaml:"auth_address"`
	UserAddress string `yaml:"user_address"`
	PayAddress  string `yaml:"payment_address"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	configPath, _ := os.LookupEnv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Путь до конфига не найден в енв файле")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Файл конфига не найден")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Не можем прочитать конфиг %s", err)
	}
	log.Printf("Config: %+v", cfg)
	return &cfg
}
