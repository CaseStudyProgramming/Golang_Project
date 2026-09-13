package config

import (
	"fmt"
	"os"
	"strings"

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

// LoadConfig loads configuration from yaml file or environment variables only
func LoadConfig(filePath string) (*Config, error) {
	var cfg Config

	// Load from file if path is provided and file exists
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&cfg); err != nil {
			return nil, err
		}
	}

	// Environment variables take precedence over config file values
	// This also provides defaults when no config file is loaded

	// Server port from environment
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	} else if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	// Fallback to environment variables for JWT secret
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	}

	// Default JWT secret for development (should be changed in production)
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "dev-secret-key-change-in-production"
	}

	// Database configuration from environment variables
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	} else if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		var portInt int
		if _, err := fmt.Sscanf(port, "%d", &portInt); err == nil {
			cfg.Database.Port = portInt
		}
	} else if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}

	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	} else if cfg.Database.User == "" {
		cfg.Database.User = "postgres"
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	} else if cfg.Database.Password == "" {
		cfg.Database.Password = "postgres"
	}

	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.Database.DBName = dbname
	} else if cfg.Database.DBName == "" {
		cfg.Database.DBName = "taskmanager"
	}

	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	} else if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	// Default CORS origins for development (should be restricted in production)
	if len(cfg.CORS.AllowedOrigins) == 0 {
		// Check for environment variable (comma-separated)
		if corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); corsOrigins != "" {
			// Parse comma-separated origins
			cfg.CORS.AllowedOrigins = strings.Split(corsOrigins, ",")
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
