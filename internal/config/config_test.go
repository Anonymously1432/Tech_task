package config_test

import (
	"os"
	"testing"

	"tech_task/internal/config"
)

func TestLoad_Success(t *testing.T) {
	oldUser := os.Getenv("POSTGRES_USER")
	oldPassword := os.Getenv("POSTGRES_PASSWORD")
	oldHost := os.Getenv("POSTGRES_HOST")
	oldPort := os.Getenv("POSTGRES_PORT")
	oldDB := os.Getenv("POSTGRES_DB")
	defer func() {
		os.Setenv("POSTGRES_USER", oldUser)
		os.Setenv("POSTGRES_PASSWORD", oldPassword)
		os.Setenv("POSTGRES_HOST", oldHost)
		os.Setenv("POSTGRES_PORT", oldPort)
		os.Setenv("POSTGRES_DB", oldDB)
	}()

	// Устанавливаем тестовые значения
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "testdb")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DBUser != "testuser" {
		t.Errorf("expected DBUser=testuser, got %s", cfg.DBUser)
	}
	if cfg.DBPassword != "testpass" {
		t.Errorf("expected DBPassword=testpass, got %s", cfg.DBPassword)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost=localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("expected DBPort=5432, got %s", cfg.DBPort)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("expected DBName=testdb, got %s", cfg.DBName)
	}
}

func TestLoad_MissingEnv(t *testing.T) {
	os.Setenv("POSTGRES_USER", "")
	os.Setenv("POSTGRES_PASSWORD", "")
	os.Setenv("POSTGRES_HOST", "")
	os.Setenv("POSTGRES_PORT", "")
	os.Setenv("POSTGRES_DB", "")

	cfg, err := config.Load()
	if err == nil {
		t.Fatal("expected error due to missing environment variables, got nil")
	}
	if cfg != nil {
		t.Errorf("expected cfg to be nil, got %+v", cfg)
	}
}

func TestDSN(t *testing.T) {
	cfg := &config.Config{
		DBUser:     "user",
		DBPassword: "pass",
		DBHost:     "host",
		DBPort:     "5432",
		DBName:     "dbname",
	}

	expected := "host=host user=user password=pass dbname=dbname port=5432 sslmode=disable TimeZone=UTC"
	if dsn := cfg.DSN(); dsn != expected {
		t.Errorf("expected DSN=%q, got %q", expected, dsn)
	}
}
