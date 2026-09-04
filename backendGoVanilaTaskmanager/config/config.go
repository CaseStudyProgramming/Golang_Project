package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	CORS     CORSConfig     `yaml:"cors"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"` // "development" or "production"
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// LoadConfig loads configuration from yaml file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	// Fallback to environment variables for JWT secret
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	}

	// Default JWT secret for development (should be changed in production)
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "dev-secret-key-change-in-production"
	}

	// Default CORS origins for development (should be restricted in production)
	if len(cfg.CORS.AllowedOrigins) == 0 {
		// Check for environment variable
		if corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); corsOrigins != "" {
			cfg.CORS.AllowedOrigins = []string{corsOrigins}
		} else {
			cfg.CORS.AllowedOrigins = []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8080"}
		}
	}

	// Default environment to development
	if cfg.Server.Env == "" {
		cfg.Server.Env = os.Getenv("APP_ENV")
		if cfg.Server.Env == "" {
			cfg.Server.Env = "development"
		}
	}

	return &cfg, nil
}
