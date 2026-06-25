package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	MigrationsPath string
	PostgresConfig PostgresConfig
}

type PostgresConfig struct {
	DSN               string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

func Load() (*Config, error) {
	cgf := &Config{
		Port: getenv("PORT", "8080"),
		PostgresConfig: PostgresConfig{
			DSN:               getDatabaseURL(),
			MaxConns:          int32(getenvInt("DB_MAX_CONNS", 20)),
			MinConns:          int32(getenvInt("DB_MIN_CONNS", 5)),
			MaxConnLifetime:   getenvDurationSeconds("DB_MAX_CONN_LIFETIME", 3600),
			MaxConnIdleTime:   getenvDurationSeconds("DB_MAX_CONN_IDLE_TIME", 900),
			HealthCheckPeriod: getenvDurationSeconds("DB_HEALTHCHECK_PERIOD", 60),
		},
		MigrationsPath: getenv("MIGRATIONS_PATH", "migrations"),
	}
	return cgf, nil
}

func getenv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getenvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}

	return i
}

func getenvDurationSeconds(key string, defSeconds int) time.Duration {
	return time.Duration(getenvInt(key, defSeconds)) * time.Second
}

func getDatabaseURL() string {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("POSTGRES_USER", "postgres")
	password := getenv("POSTGRES_PASSWORD", "password")
	dbname := getenv("POSTGRES_DB", "bloq")
	sslmode := getenv("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}
