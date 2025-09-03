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

type StorageConfig struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	BucketName string `yaml:"bucket_name"`
	UseSSL     bool   `yaml:"use_ssl"`
	Region     string `yaml:"region"`
	PublicURL  string `yaml:"public_url"`
}

type EmailConfig struct {
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     string `yaml:"smtp_port"`
	SMTPUser     string `yaml:"smtp_user"`
	SMTPPassword string `yaml:"smtp_password"`
	FromEmail    string `yaml:"from_email"`
	FromName     string `yaml:"from_name"`
	TemplatesDir string `yaml:"templates_dir"`
}

type Config struct {
	Port         string        `yaml:"port"`
	MongoURL     string        `yaml:"mongo_url"`
	LogLevel     string        `yaml:"log_level"`
	DatabaseName string        `yaml:"database_name"`
	JWTSecret    string        `yaml:"jwt_secret"`
	Redis        RedisConfig   `yaml:"redis"`
	Storage      StorageConfig `yaml:"storage"`
	Email        EmailConfig   `yaml:"email"`
	WorkerCount  int           `yaml:"worker_count"`
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
