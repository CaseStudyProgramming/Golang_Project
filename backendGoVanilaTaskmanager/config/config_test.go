package config

import (
	"os"
	"testing"
)

func TestLoadConfig_EnvironmentOnly(t *testing.T) {
	// Clean up environment after test
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		os.Unsetenv("APP_ENV")
	}()

	// Set some environment variables
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://example.com,http://test.com")
	os.Setenv("APP_ENV", "production")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Server.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", cfg.Server.Port)
	}

	if cfg.JWT.Secret != "test-secret" {
		t.Errorf("Expected JWT secret 'test-secret', got %s", cfg.JWT.Secret)
	}

	if cfg.Database.Host != "testhost" {
		t.Errorf("Expected DB host 'testhost', got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5433 {
		t.Errorf("Expected DB port 5433, got %d", cfg.Database.Port)
	}

	if cfg.Database.User != "testuser" {
		t.Errorf("Expected DB user 'testuser', got %s", cfg.Database.User)
	}

	if cfg.Database.Password != "testpass" {
		t.Errorf("Expected DB password 'testpass', got %s", cfg.Database.Password)
	}

	if cfg.Database.DBName != "testdb" {
		t.Errorf("Expected DB name 'testdb', got %s", cfg.Database.DBName)
	}

	if cfg.Database.SSLMode != "require" {
		t.Errorf("Expected SSL mode 'require', got %s", cfg.Database.SSLMode)
	}

	if cfg.Server.Env != "production" {
		t.Errorf("Expected env 'production', got %s", cfg.Server.Env)
	}

	expectedOrigins := []string{"http://example.com", "http://test.com"}
	if len(cfg.CORS.AllowedOrigins) != len(expectedOrigins) {
		t.Errorf("Expected %d CORS origins, got %d", len(expectedOrigins), len(cfg.CORS.AllowedOrigins))
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Clean up environment after test
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		os.Unsetenv("APP_ENV")
	}()

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check defaults
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Server.Port)
	}

	if cfg.JWT.Secret != "dev-secret-key-change-in-production" {
		t.Errorf("Expected default JWT secret, got %s", cfg.JWT.Secret)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected default DB host 'localhost', got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected default DB port 5432, got %d", cfg.Database.Port)
	}

	if cfg.Database.User != "postgres" {
		t.Errorf("Expected default DB user 'postgres', got %s", cfg.Database.User)
	}

	if cfg.Database.Password != "postgres" {
		t.Errorf("Expected default DB password 'postgres', got %s", cfg.Database.Password)
	}

	if cfg.Database.DBName != "taskmanager" {
		t.Errorf("Expected default DB name 'taskmanager', got %s", cfg.Database.DBName)
	}

	if cfg.Database.SSLMode != "disable" {
		t.Errorf("Expected default SSL mode 'disable', got %s", cfg.Database.SSLMode)
	}

	if cfg.Server.Env != "development" {
		t.Errorf("Expected default env 'development', got %s", cfg.Server.Env)
	}

	// Check default CORS origins
	if len(cfg.CORS.AllowedOrigins) == 0 {
		t.Error("Expected default CORS origins")
	}
}

func TestLoadConfig_InvalidDBPort(t *testing.T) {
	// Clean up environment after test
	defer func() {
		os.Unsetenv("DB_PORT")
	}()

	os.Setenv("DB_PORT", "invalid")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Should use default port when invalid port is provided
	if cfg.Database.Port != 5432 {
		t.Errorf("Expected default port 5432 when invalid port provided, got %d", cfg.Database.Port)
	}
}

func TestLoadConfig_EnvVarsOverrideDefaults(t *testing.T) {
	// Clean up environment after test
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		os.Unsetenv("APP_ENV")
	}()

	// Set only some environment variables
	os.Setenv("SERVER_PORT", "7000")
	os.Setenv("DB_HOST", "customhost")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check that env vars override defaults
	if cfg.Server.Port != "7000" {
		t.Errorf("Expected port 7000 from env, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "customhost" {
		t.Errorf("Expected DB host 'customhost' from env, got %s", cfg.Database.Host)
	}

	// Check that defaults are used for unset env vars
	if cfg.Database.Port != 5432 {
		t.Errorf("Expected default port 5432, got %d", cfg.Database.Port)
	}

	if cfg.Database.User != "postgres" {
		t.Errorf("Expected default user 'postgres', got %s", cfg.Database.User)
	}
}
