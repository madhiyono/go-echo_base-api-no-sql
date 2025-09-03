package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Port         string      `yaml:"port"`
	MongoURL     string      `yaml:"mongo_url"`
	LogLevel     string      `yaml:"log_level"`
	DatabaseName string      `yaml:"database_name"`
	JWTSecret    string      `yaml:"jwt_secret"`
	Redis        RedisConfig `yaml:"redis"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	if mongoURL := os.Getenv("MONGO_URL"); mongoURL != "" {
		cfg.MongoURL = mongoURL
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}

	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}

	return &cfg, nil
}
