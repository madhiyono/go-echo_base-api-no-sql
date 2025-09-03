package storage

// StorageConfig holds the configuration for storage service
type StorageConfig struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	BucketName string `yaml:"bucket_name"`
	UseSSL     bool   `yaml:"use_ssl"`
	Region     string `yaml:"region"`
	PublicURL  string `yaml:"public_url"`
}

// Config is an alias for StorageConfig for easier access
type Config = StorageConfig
